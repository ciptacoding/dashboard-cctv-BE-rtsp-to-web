package service

import (
	"cctv-monitoring-backend/internal/models"
	"cctv-monitoring-backend/internal/repository"
	ws "cctv-monitoring-backend/internal/websocket"
	"database/sql"
	"fmt"
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

	// Rate limiting untuk health checks
	checkSemaphore chan struct{}

	// Track restart attempts untuk exponential backoff
	restartAttempts map[string]restartAttempt
	restartMu       sync.RWMutex
}

// restartAttempt tracks restart attempts for exponential backoff
type restartAttempt struct {
	count       int
	lastAttempt time.Time
	nextRetry   time.Time
}

// CameraHealthStatus tracks the health status of a camera
type CameraHealthStatus struct {
	Status           string
	LastChecked      time.Time
	ConsecutiveFails int
	LastSuccessTime  time.Time
	LastSnapshotHash string // Track snapshot hash for frozen detection
}

// NewCameraHealthMonitor creates a new camera health monitor
func NewCameraHealthMonitor(
	cameraRepo repository.CameraRepository,
	rtspService RTSPService,
	wsHub *ws.Hub,
	interval time.Duration,
) *CameraHealthMonitor {
	// Limit concurrent health checks to prevent overwhelming RTSP service
	// Allow max 5 concurrent checks at a time
	maxConcurrentChecks := 5
	if maxConcurrentChecks < 3 {
		maxConcurrentChecks = 3
	}

	return &CameraHealthMonitor{
		cameraRepo:      cameraRepo,
		rtspService:     rtspService,
		wsHub:           wsHub,
		interval:        interval,
		stopChan:        make(chan bool),
		lastCheckStatus: make(map[string]CameraHealthStatus),
		checkSemaphore:  make(chan struct{}, maxConcurrentChecks),
		restartAttempts: make(map[string]restartAttempt),
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

// checkAllCameras checks health of all active cameras with rate limiting
func (m *CameraHealthMonitor) checkAllCameras() {
	// Get all active cameras
	cameras, _, err := m.cameraRepo.GetAll(1, 1000) // Get all cameras
	if err != nil {
		log.Printf("Error fetching cameras for health check: %v", err)
		return
	}

	activeCameras := make([]*models.Camera, 0)
	for _, camera := range cameras {
		if camera.IsActive {
			activeCameras = append(activeCameras, camera)
		}
	}

	log.Printf("🔍 Checking health of %d active cameras...", len(activeCameras))

	// Use semaphore to limit concurrent checks
	var wg sync.WaitGroup
	for _, camera := range activeCameras {
		wg.Add(1)
		go func(cam *models.Camera) {
			defer wg.Done()

			// Acquire semaphore (rate limiting)
			m.checkSemaphore <- struct{}{}
			defer func() { <-m.checkSemaphore }()

			m.checkCameraHealth(cam)
		}(camera)
	}

	// Wait for all checks to complete (with timeout)
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All checks completed
	case <-time.After(2 * m.interval):
		log.Printf("⚠️  Health check timeout - some cameras may not have been checked")
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
	if err != nil {
		log.Printf("⚠️  Error checking status for camera %s (%s): %v", camera.Name, camera.ID, err)
	}

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
		newStatus.LastSnapshotHash = lastStatus.LastSnapshotHash

		// If this is a new failure or consistent failure
		if newStatus.ConsecutiveFails >= 2 {
			m.handleCameraOffline(camera, status)
		}
	} else if status == "READY" || status == "ONLINE" {
		newStatus.ConsecutiveFails = 0

		// Try to get snapshot hash to verify stream is actually updating
		if snapshotHash, err := m.rtspService.GetSnapshotHash(streamID); err == nil {
			newStatus.LastSnapshotHash = snapshotHash
			newStatus.LastSuccessTime = now

			// Log if snapshot hash changed (stream is updating)
			if lastStatus.LastSnapshotHash != "" && snapshotHash != lastStatus.LastSnapshotHash {
				log.Printf("✓ Camera %s (%s) stream is updating", camera.Name, camera.ID)
			}
		} else {
			// If we can't get snapshot, use last known values
			log.Printf("⚠️  Could not get snapshot for camera %s (%s): %v", camera.Name, camera.ID, err)
			newStatus.LastSnapshotHash = lastStatus.LastSnapshotHash
			if lastStatus.LastSuccessTime.IsZero() {
				newStatus.LastSuccessTime = now
			} else {
				newStatus.LastSuccessTime = lastStatus.LastSuccessTime
			}
		}

		// If camera was previously offline, notify it's back online
		if lastStatus.Status == "OFFLINE" && lastStatus.ConsecutiveFails >= 2 {
			m.handleCameraOnline(camera)
		}
	}

	// Detect frozen stream using snapshot comparison
	if status == "READY" {
		// Check if stream is actually frozen by comparing snapshots
		if m.isStreamFrozen(camera.ID, streamID) {
			m.handleCameraFrozen(camera)
		}
	}

	// Update status in memory
	m.mu.Lock()
	m.lastCheckStatus[camera.ID] = newStatus
	m.mu.Unlock()

	// Always update database with current status and last_seen
	// This ensures status is always synced with database
	camera.Status = status
	camera.LastSeen = sql.NullTime{Time: now, Valid: true}

	if err := m.cameraRepo.Update(camera.ID, camera); err != nil {
		log.Printf("Error updating camera status in DB: %v", err)
	} else {
		log.Printf("✓ Updated camera %s status to %s in database", camera.ID, status)
	}
}

// handleCameraOffline handles when a camera goes offline
func (m *CameraHealthMonitor) handleCameraOffline(camera *models.Camera, status string) {
	log.Printf("⚠️  Camera %s (%s) is OFFLINE", camera.Name, camera.ID)

	// Update camera status in database
	camera.Status = "OFFLINE"
	camera.LastSeen = sql.NullTime{Time: time.Now(), Valid: true}

	if err := m.cameraRepo.Update(camera.ID, camera); err != nil {
		log.Printf("Error updating camera status to OFFLINE in database: %v", err)
	} else {
		log.Printf("✓ Updated camera %s status to OFFLINE in database", camera.ID)
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
		log.Printf("Error updating camera status to ONLINE in database: %v", err)
	} else {
		log.Printf("✓ Updated camera %s status to ONLINE in database", camera.ID)
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

// isStreamFrozen checks if a stream is frozen by comparing snapshot hashes
func (m *CameraHealthMonitor) isStreamFrozen(cameraID, streamID string) bool {
	m.mu.RLock()
	lastStatus, exists := m.lastCheckStatus[cameraID]
	m.mu.RUnlock()

	if !exists || lastStatus.LastSnapshotHash == "" {
		// First check or no previous hash, can't determine if frozen
		return false
	}

	// Get current snapshot hash
	currentHash, err := m.rtspService.GetSnapshotHash(streamID)
	if err != nil {
		// If we can't get snapshot, check time-based fallback
		timeSinceLastSuccess := time.Since(lastStatus.LastSuccessTime)
		// If no update in 60 seconds, consider frozen
		return timeSinceLastSuccess > 60*time.Second
	}

	// Compare hashes - if same hash for too long, stream is frozen
	if currentHash == lastStatus.LastSnapshotHash {
		timeSinceLastSuccess := time.Since(lastStatus.LastSuccessTime)
		// Same snapshot for more than 45 seconds = frozen
		if timeSinceLastSuccess > 45*time.Second {
			log.Printf("🧊 Camera %s appears frozen - same snapshot hash for %v", cameraID, timeSinceLastSuccess)
			return true
		}
	}

	return false
}

// attemptStreamRestart attempts to restart a camera stream with exponential backoff
func (m *CameraHealthMonitor) attemptStreamRestart(camera *models.Camera) {
	cameraID := camera.ID

	// Check if we should retry (exponential backoff)
	m.restartMu.Lock()
	attempt, exists := m.restartAttempts[cameraID]
	now := time.Now()

	if exists && now.Before(attempt.nextRetry) {
		m.restartMu.Unlock()
		log.Printf("⏳ Skipping restart for camera %s - next retry in %v", camera.Name, attempt.nextRetry.Sub(now))
		return
	}

	// Update restart attempt tracking
	if !exists {
		attempt = restartAttempt{
			count:       0,
			lastAttempt: now,
		}
	}

	attempt.count++
	attempt.lastAttempt = now

	// Exponential backoff: 2^count seconds, max 5 minutes
	backoffSeconds := 1 << uint(attempt.count-1) // 1, 2, 4, 8, 16, 32, 64, 128...
	if backoffSeconds > 300 {                    // Max 5 minutes
		backoffSeconds = 300
	}
	attempt.nextRetry = now.Add(time.Duration(backoffSeconds) * time.Second)

	m.restartAttempts[cameraID] = attempt
	m.restartMu.Unlock()

	log.Printf("🔄 Attempting to restart stream for camera %s (attempt #%d, backoff: %ds)",
		camera.Name, attempt.count, backoffSeconds)

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
			fmt.Sprintf("Failed to restart stream (attempt #%d). Will retry in %ds...", attempt.count, backoffSeconds),
		)
		return
	}

	// Success! Reset restart attempts
	m.restartMu.Lock()
	delete(m.restartAttempts, cameraID)
	m.restartMu.Unlock()

	// Update camera with new stream info
	camera.StreamID = sql.NullString{String: streamID, Valid: true}
	camera.HLSUrl = hlsURL
	camera.SnapshotUrl = snapshotURL
	camera.Status = "ONLINE"
	camera.LastSeen = sql.NullTime{Time: time.Now(), Valid: true}

	if err := m.cameraRepo.Update(camera.ID, camera); err != nil {
		log.Printf("Error updating camera after restart: %v", err)
	} else {
		log.Printf("✓ Updated camera %s status to ONLINE in database after restart", camera.ID)
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
