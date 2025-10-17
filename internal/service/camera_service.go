package service

import (
	"database/sql"
	"fmt"

	"cctv-monitoring-backend/internal/models"
	"cctv-monitoring-backend/internal/repository"
)

// CameraService adalah interface untuk business logic camera
type CameraService interface {
	Create(req *models.CreateCameraRequest, userID string) (*models.CameraWithStream, error)
	GetByID(id string) (*models.CameraWithStream, error)
	GetAll(page, pageSize int) ([]*models.CameraWithStream, *models.PaginationMeta, error)
	Update(id string, req *models.UpdateCameraRequest) (*models.CameraWithStream, error)
	Delete(id string) error
	GetByZone(zone string) ([]*models.CameraWithStream, error)
	GetNearby(lat, lng, radiusKm float64) ([]*models.CameraWithStream, error)
	StartStream(cameraID string) (*models.CameraWithStream, error)
	StopStream(cameraID string) error
}

type cameraService struct {
	cameraRepo  repository.CameraRepository
	rtspService RTSPService
}

// NewCameraService membuat instance baru dari CameraService
func NewCameraService(cameraRepo repository.CameraRepository, rtspService RTSPService) CameraService {
	return &cameraService{
		cameraRepo:  cameraRepo,
		rtspService: rtspService,
	}
}

// Create membuat camera baru
func (s *cameraService) Create(req *models.CreateCameraRequest, userID string) (*models.CameraWithStream, error) {
	// Validasi input
	if req.Name == "" {
		return nil, fmt.Errorf("camera name is required")
	}
	if req.RTSPUrl == "" {
		return nil, fmt.Errorf("RTSP URL is required")
	}

	// Set default FPS jika tidak diisi
	if req.FPS == 0 {
		req.FPS = 25
	}

	// Buat object camera
	camera := &models.Camera{
		Name:         req.Name,
		Description:  sql.NullString{String: req.Description, Valid: req.Description != ""},
		RTSPUrl:      req.RTSPUrl,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Building:     sql.NullString{String: req.Building, Valid: req.Building != ""},
		Zone:         sql.NullString{String: req.Zone, Valid: req.Zone != ""},
		IPAddress:    sql.NullString{String: req.IPAddress, Valid: req.IPAddress != ""},
		Port:         sql.NullInt64{Int64: int64(req.Port), Valid: req.Port > 0},
		Manufacturer: sql.NullString{String: req.Manufacturer, Valid: req.Manufacturer != ""},
		Model:        sql.NullString{String: req.Model, Valid: req.Model != ""},
		Resolution:   sql.NullString{String: req.Resolution, Valid: req.Resolution != ""},
		FPS:          req.FPS,
		Tags:         req.Tags,
		Status:       "OFFLINE",
		IsActive:     true,
	}

	// Simpan ke database
	if err := s.cameraRepo.Create(camera, userID); err != nil {
		return nil, fmt.Errorf("failed to create camera: %w", err)
	}

	// Auto-register ke RTSPtoWeb untuk on-demand streaming
	streamID, err := s.rtspService.AddStream(camera.ID, camera.RTSPUrl)
	if err != nil {
		// Log error tapi tidak fail, karena bisa di-register nanti
		fmt.Printf("Warning: failed to register stream: %v\n", err)
	} else {
		// Update stream ID di database
		if err := s.cameraRepo.UpdateStreamID(camera.ID, streamID); err != nil {
			fmt.Printf("Warning: failed to update stream ID: %v\n", err)
		}
		// Update status ke READY (sudah registered, siap on-demand)
		s.cameraRepo.UpdateStatus(camera.ID, "READY")
	}

	// Ambil data camera terbaru
	updatedCamera, err := s.cameraRepo.GetByID(camera.ID)
	if err != nil {
		return nil, err
	}

	// Return camera dengan stream URLs
	return s.buildCameraWithStream(updatedCamera), nil
}

// GetByID mengambil camera berdasarkan ID
func (s *cameraService) GetByID(id string) (*models.CameraWithStream, error) {
	camera, err := s.cameraRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.buildCameraWithStream(camera), nil
}

// GetAll mengambil semua camera dengan pagination
func (s *cameraService) GetAll(page, pageSize int) ([]*models.CameraWithStream, *models.PaginationMeta, error) {
	// Set default pagination jika tidak diisi
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Ambil data dari repository
	cameras, total, err := s.cameraRepo.GetAll(pageSize, offset)
	if err != nil {
		return nil, nil, err
	}

	// Convert ke CameraWithStream
	result := make([]*models.CameraWithStream, 0, len(cameras))
	for _, camera := range cameras {
		result = append(result, s.buildCameraWithStream(camera))
	}

	// Hitung total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	// Buat pagination meta
	meta := &models.PaginationMeta{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}

	return result, meta, nil
}

// Update mengupdate data camera
func (s *cameraService) Update(id string, req *models.UpdateCameraRequest) (*models.CameraWithStream, error) {
	// Ambil data camera yang ada
	existingCamera, err := s.cameraRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields jika ada di request
	if req.Name != "" {
		existingCamera.Name = req.Name
	}
	if req.Description != "" {
		existingCamera.Description = sql.NullString{String: req.Description, Valid: true}
	}
	if req.RTSPUrl != "" {
		existingCamera.RTSPUrl = req.RTSPUrl
	}
	if req.Latitude != 0 {
		existingCamera.Latitude = req.Latitude
	}
	if req.Longitude != 0 {
		existingCamera.Longitude = req.Longitude
	}
	if req.Building != "" {
		existingCamera.Building = sql.NullString{String: req.Building, Valid: true}
	}
	if req.Zone != "" {
		existingCamera.Zone = sql.NullString{String: req.Zone, Valid: true}
	}
	if req.IPAddress != "" {
		existingCamera.IPAddress = sql.NullString{String: req.IPAddress, Valid: true}
	}
	if req.Port > 0 {
		existingCamera.Port = sql.NullInt64{Int64: int64(req.Port), Valid: true}
	}
	if req.Manufacturer != "" {
		existingCamera.Manufacturer = sql.NullString{String: req.Manufacturer, Valid: true}
	}
	if req.Model != "" {
		existingCamera.Model = sql.NullString{String: req.Model, Valid: true}
	}
	if req.Resolution != "" {
		existingCamera.Resolution = sql.NullString{String: req.Resolution, Valid: true}
	}
	if req.FPS > 0 {
		existingCamera.FPS = req.FPS
	}
	if len(req.Tags) > 0 {
		existingCamera.Tags = req.Tags
	}
	if req.Status != "" {
		existingCamera.Status = req.Status
	}

	// Update di database
	if err := s.cameraRepo.Update(id, existingCamera); err != nil {
		return nil, fmt.Errorf("failed to update camera: %w", err)
	}

	// Return updated camera
	return s.buildCameraWithStream(existingCamera), nil
}

// Delete menghapus camera (soft delete)
func (s *cameraService) Delete(id string) error {
	// Cek apakah camera ada
	camera, err := s.cameraRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Jika ada stream ID, hapus dari RTSPtoWeb
	if camera.StreamID.Valid && camera.StreamID.String != "" {
		if err := s.rtspService.RemoveStream(camera.StreamID.String); err != nil {
			// Log error tapi tetap lanjut hapus dari database
			fmt.Printf("Warning: failed to remove stream from RTSPtoWeb: %v\n", err)
		}
	}

	// Hapus dari database (soft delete)
	return s.cameraRepo.Delete(id)
}

// GetByZone mengambil camera berdasarkan zone
func (s *cameraService) GetByZone(zone string) ([]*models.CameraWithStream, error) {
	cameras, err := s.cameraRepo.GetByZone(zone)
	if err != nil {
		return nil, err
	}

	result := make([]*models.CameraWithStream, 0, len(cameras))
	for _, camera := range cameras {
		result = append(result, s.buildCameraWithStream(camera))
	}

	return result, nil
}

// GetNearby mengambil camera dalam radius tertentu
func (s *cameraService) GetNearby(lat, lng, radiusKm float64) ([]*models.CameraWithStream, error) {
	cameras, err := s.cameraRepo.GetNearby(lat, lng, radiusKm)
	if err != nil {
		return nil, err
	}

	result := make([]*models.CameraWithStream, 0, len(cameras))
	for _, camera := range cameras {
		result = append(result, s.buildCameraWithStream(camera))
	}

	return result, nil
}

// StartStream memulai streaming camera (on-demand mode)
func (s *cameraService) StartStream(cameraID string) (*models.CameraWithStream, error) {
	// Ambil data camera
	camera, err := s.cameraRepo.GetByID(cameraID)
	if err != nil {
		return nil, err
	}

	// Check apakah sudah ada stream ID
	if !camera.StreamID.Valid || camera.StreamID.String == "" {
		// Register stream ke RTSPtoWeb (akan auto-start on-demand)
		streamID, err := s.rtspService.AddStream(cameraID, camera.RTSPUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to register stream: %w", err)
		}

		// Update stream ID di database
		if err := s.cameraRepo.UpdateStreamID(cameraID, streamID); err != nil {
			return nil, fmt.Errorf("failed to update stream ID: %w", err)
		}
	}

	// Update status ke READY (bukan ONLINE, karena on-demand)
	if err := s.cameraRepo.UpdateStatus(cameraID, "READY"); err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	// Ambil data camera terbaru
	updatedCamera, err := s.cameraRepo.GetByID(cameraID)
	if err != nil {
		return nil, err
	}

	return s.buildCameraWithStream(updatedCamera), nil
}

// StopStream menghentikan streaming camera (on-demand mode)
func (s *cameraService) StopStream(cameraID string) error {
	// Pada on-demand mode, kita tidak perlu hapus dari RTSPtoWeb
	// Cukup update status ke OFFLINE
	// Stream akan auto-stop saat tidak ada viewer

	if err := s.cameraRepo.UpdateStatus(cameraID, "OFFLINE"); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}

// buildCameraWithStream menambahkan stream URLs ke camera object
func (s *cameraService) buildCameraWithStream(camera *models.Camera) *models.CameraWithStream {
	result := &models.CameraWithStream{
		Camera: *camera,
	}

	// Jika ada stream ID, tambahkan URLs
	if camera.StreamID.Valid && camera.StreamID.String != "" {
		result.HLSUrl = s.rtspService.GetHLSURL(camera.StreamID.String)
		result.WebRTCUrl = s.rtspService.GetWebRTCURL(camera.StreamID.String)
		result.SnapshotURL = s.rtspService.GetSnapshotURL(camera.StreamID.String)
		result.StreamURL = result.HLSUrl // Backward compatibility
	}

	return result
}
