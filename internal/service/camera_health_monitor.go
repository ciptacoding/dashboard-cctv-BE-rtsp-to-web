package service

import (
	"cctv-monitoring-backend/internal/models"
	"cctv-monitoring-backend/internal/repository"
	ws "cctv-monitoring-backend/internal/websocket"
	"database/sql"
	"log"
	"sync"
	"time"
)

// CameraHealthMonitor monitors camera health and broadcasts updates
type CameraHealthMonitor struct {
	cameraRepo  repository.CameraRepository
	rtspService RTSPService
	wsHub       *ws.Hub
	interval    time.Duration
	stopChan    chan bool
	mu          sync.RWMutex

	// Track camera check history untuk detect frozen streams
	lastCheckStatus map[string]CameraHealthStatus
}

// CameraHealthStatus tracks the health status of a camera
type CameraHealthStatus struct {
	Status           string
	LastChecked      time.Time
	ConsecutiveFails int
	LastSuccessTime  time.Time
}

// NewCameraHealthMonitor creates a new camera health monitor
func NewCameraHealthMonitor(
	cameraRepo repository.CameraRepository,
	rtspService RTSPService,
	wsHub *ws.Hub,
	interval time.Duration,
) *CameraHealthMonitor {
	return &CameraHealthMonitor{
		cameraRepo:      cameraRepo,
		rtspService:     rtspService,
		wsHub:           wsHub,
		interval:        interval,
		stopChan:        make(chan bool),
		lastCheckStatus: make(map[string]CameraHealthStatus),
	}
}

// Start begins monitoring camera health
func (m *CameraHealthMonitor) Start() {
	log.Printf("✓ Camera health monitor started (interval: %v)", m.interval)

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	// Run initial check immediately
	go m.checkAllCameras()

	for {
		select {
		case <-ticker.C:
			go m.checkAllCameras()
		case <-m.stopChan:
			log.Println("✓ Camera health monitor stopped")
			return
		}
	}
}

// Stop stops the health monitor
func (m *CameraHealthMonitor) Stop() {
	m.stopChan <- true
}

// checkAllCameras checks health of all active cameras
func (m *CameraHealthMonitor) checkAllCameras() {
	// Get all active cameras
	cameras, _, err := m.cameraRepo.GetAll(1, 1000) // Get all cameras
	if err != nil {
		log.Printf("Error fetching cameras for health check: %v", err)
		return
	}

	log.Printf("🔍 Checking health of %d cameras...", len(cameras))

	for _, camera := range cameras {
		if !camera.IsActive {
			continue
		}

		go m.checkCameraHealth(camera)
	}
}

// checkCameraHealth checks health of a single camera
func (m *CameraHealthMonitor) checkCameraHealth(camera *models.Camera) {
	// Skip if camera doesn't have a stream
	if !camera.StreamID.Valid || camera.StreamID.String == "" {
		return
	}

	streamID := camera.StreamID.String

	// Check stream status from RTSPtoWeb
	status, err := m.rtspService.GetStreamStatus(streamID)

	m.mu.Lock()
	lastStatus := m.lastCheckStatus[camera.ID]
	m.mu.Unlock()

	now := time.Now()
	newStatus := CameraHealthStatus{
		Status:      status,
		LastChecked: now,
	}

	// Determine if camera is having issues
	if err != nil || status == "OFFLINE" || status == "ERROR" {
		newStatus.ConsecutiveFails = lastStatus.ConsecutiveFails + 1
		newStatus.LastSuccessTime = lastStatus.LastSuccessTime

		// If this is a new failure or consistent failure
		if newStatus.ConsecutiveFails >= 2 {
			m.handleCameraOffline(camera, status)
		}
	} else if status == "READY" || status == "ONLINE" {
		newStatus.ConsecutiveFails = 0
		newStatus.LastSuccessTime = now

		// If camera was previously offline, notify it's back online
		if lastStatus.Status == "OFFLINE" && lastStatus.ConsecutiveFails >= 2 {
			m.handleCameraOnline(camera)
		}
	}

	// Detect frozen stream (stream is "READY" but no frame updates)
	if status == "READY" {
		// Check if camera hasn't updated in a while
		if camera.LastSeen.Valid {
			lastSeenTime := camera.LastSeen.Time
			timeSinceLastSeen := now.Sub(lastSeenTime)

			// If no update in 60 seconds, consider frozen
			if timeSinceLastSeen > 60*time.Second {
				m.handleCameraFrozen(camera)
			}
		}
	}

	// Update status in memory
	m.mu.Lock()
	m.lastCheckStatus[camera.ID] = newStatus
	m.mu.Unlock()

	// Update database if status changed
	if camera.Status != status {
		camera.Status = status
		camera.LastSeen = sql.NullTime{Time: now, Valid: true}

		if err := m.cameraRepo.Update(camera.ID, camera); err != nil {
			log.Printf("Error updating camera status in DB: %v", err)
		}
	}
}

// handleCameraOffline handles when a camera goes offline
func (m *CameraHealthMonitor) handleCameraOffline(camera *models.Camera, status string) {
	log.Printf("⚠️  Camera %s (%s) is OFFLINE", camera.Name, camera.ID)

	// Update camera status in database
	camera.Status = "OFFLINE"
	camera.LastSeen = sql.NullTime{Time: time.Now(), Valid: true}

	if err := m.cameraRepo.Update(camera.ID, camera); err != nil {
		log.Printf("Error updating camera status: %v", err)
	}

	// Broadcast to WebSocket clients
	m.wsHub.BroadcastCameraStatus(
		camera.ID,
		"OFFLINE",
		camera.LastSeen.Time.Format(time.RFC3339),
	)

	// Send stream update notification
	m.wsHub.BroadcastStreamUpdate(
		camera.ID,
		camera.Name,
		"offline",
		"Camera stream is offline. Attempting to reconnect...",
	)

	// Attempt to restart stream
	go m.attemptStreamRestart(camera)
}

// handleCameraOnline handles when a camera comes back online
func (m *CameraHealthMonitor) handleCameraOnline(camera *models.Camera) {
	log.Printf("✓ Camera %s (%s) is back ONLINE", camera.Name, camera.ID)

	// Update camera status in database
	camera.Status = "ONLINE"
	camera.LastSeen = sql.NullTime{Time: time.Now(), Valid: true}

	if err := m.cameraRepo.Update(camera.ID, camera); err != nil {
		log.Printf("Error updating camera status: %v", err)
	}

	// Broadcast to WebSocket clients
	m.wsHub.BroadcastCameraStatus(
		camera.ID,
		"ONLINE",
		camera.LastSeen.Time.Format(time.RFC3339),
	)

	// Send stream update notification
	m.wsHub.BroadcastStreamUpdate(
		camera.ID,
		camera.Name,
		"online",
		"Camera stream is back online",
	)
}

// handleCameraFrozen handles when a camera stream is frozen
func (m *CameraHealthMonitor) handleCameraFrozen(camera *models.Camera) {
	log.Printf("🧊 Camera %s (%s) stream appears FROZEN", camera.Name, camera.ID)

	// Broadcast to WebSocket clients
	m.wsHub.BroadcastStreamUpdate(
		camera.ID,
		camera.Name,
		"frozen",
		"Camera stream appears frozen. Refreshing...",
	)

	// Attempt to restart stream
	go m.attemptStreamRestart(camera)
}

// attemptStreamRestart attempts to restart a camera stream
func (m *CameraHealthMonitor) attemptStreamRestart(camera *models.Camera) {
	log.Printf("🔄 Attempting to restart stream for camera %s", camera.Name)

	// Stop existing stream
	if camera.StreamID.Valid && camera.StreamID.String != "" {
		if err := m.rtspService.RemoveStream(camera.StreamID.String); err != nil {
			log.Printf("Error stopping stream: %v", err)
		}

		// Wait a bit before restarting
		time.Sleep(2 * time.Second)
	}

	// Start new stream
	streamID, hlsURL, snapshotURL, err := m.rtspService.AddStream(
		camera.ID,
		camera.Name,
		camera.RTSPUrl,
	)

	if err != nil {
		log.Printf("Error restarting stream: %v", err)

		// Notify failure
		m.wsHub.BroadcastStreamUpdate(
			camera.ID,
			camera.Name,
			"restart_failed",
			"Failed to restart stream. Will retry...",
		)
		return
	}

	// Update camera with new stream info
	camera.StreamID = sql.NullString{String: streamID, Valid: true}
	camera.HLSUrl = hlsURL
	camera.SnapshotUrl = snapshotURL
	camera.Status = "ONLINE"
	camera.LastSeen = sql.NullTime{Time: time.Now(), Valid: true}

	if err := m.cameraRepo.Update(camera.ID, camera); err != nil {
		log.Printf("Error updating camera after restart: %v", err)
	}

	// Notify success
	log.Printf("✓ Stream restarted successfully for camera %s", camera.Name)

	m.wsHub.BroadcastStreamUpdate(
		camera.ID,
		camera.Name,
		"restarted",
		"Stream restarted successfully",
	)

	m.wsHub.BroadcastCameraStatus(
		camera.ID,
		"ONLINE",
		camera.LastSeen.Time.Format(time.RFC3339),
	)
}

// GetHealthStatus returns current health status of all cameras
func (m *CameraHealthMonitor) GetHealthStatus() map[string]CameraHealthStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy to avoid race conditions
	status := make(map[string]CameraHealthStatus)
	for k, v := range m.lastCheckStatus {
		status[k] = v
	}

	return status
}
