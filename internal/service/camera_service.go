package service

import (
	"cctv-monitoring-backend/internal/models"
	"cctv-monitoring-backend/internal/repository"
	"database/sql"
	"fmt"
)

type CameraService interface {
	Create(req *models.CreateCameraRequest, userID string) (*models.Camera, error)
	GetByID(id string) (*models.Camera, error)
	GetAll(page, pageSize int) ([]*models.Camera, *models.PaginationMeta, error)
	Update(id string, req *models.UpdateCameraRequest) (*models.Camera, error)
	Delete(id string) error
	GetByZone(zone string) ([]*models.Camera, error)
	GetNearby(lat, lng, radius float64) ([]*models.Camera, error)
	StartStream(id string) (*models.Camera, error)
	StopStream(id string) error
}

type cameraService struct {
	cameraRepo  repository.CameraRepository
	rtspService RTSPService
}

func NewCameraService(cameraRepo repository.CameraRepository, rtspService RTSPService) CameraService {
	return &cameraService{
		cameraRepo:  cameraRepo,
		rtspService: rtspService,
	}
}

// enrichCameraWithStreamURLs menambahkan stream URLs ke camera
func (s *cameraService) enrichCameraWithStreamURLs(camera *models.Camera) {
	if camera.StreamID.Valid && camera.StreamID.String != "" {
		streamID := camera.StreamID.String
		camera.HLSUrl = s.rtspService.GetHLSURL(streamID)
		camera.SnapshotUrl = s.rtspService.GetSnapshotURL(streamID)
	}
}

// enrichCamerasWithStreamURLs menambahkan stream URLs ke array cameras
func (s *cameraService) enrichCamerasWithStreamURLs(cameras []*models.Camera) {
	for _, camera := range cameras {
		s.enrichCameraWithStreamURLs(camera)
	}
}

func (s *cameraService) Create(req *models.CreateCameraRequest, userID string) (*models.Camera, error) {
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
		Status:       req.Status,
		IsActive:     true,
		CreatedBy:    sql.NullString{String: userID, Valid: true},
	}

	// Create camera in database
	if err := s.cameraRepo.Create(camera, userID); err != nil {
		return nil, fmt.Errorf("failed to create camera: %w", err)
	}

	// Add stream to RTSPtoWeb
	streamID, hlsURL, snapshotURL, err := s.rtspService.AddStream(
		camera.ID,
		camera.Name,
		camera.RTSPUrl,
	)

	if err == nil {
		camera.StreamID = sql.NullString{String: streamID, Valid: true}
		camera.HLSUrl = hlsURL
		camera.SnapshotUrl = snapshotURL

		// Update camera dengan stream info
		s.cameraRepo.Update(camera.ID, camera)
	}

	// Enrich dengan stream URLs
	s.enrichCameraWithStreamURLs(camera)

	return camera, nil
}

func (s *cameraService) GetByID(id string) (*models.Camera, error) {
	camera, err := s.cameraRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("camera not found: %w", err)
	}

	// Enrich dengan stream URLs
	s.enrichCameraWithStreamURLs(camera)

	return camera, nil
}

func (s *cameraService) GetAll(page, pageSize int) ([]*models.Camera, *models.PaginationMeta, error) {
	cameras, meta, err := s.cameraRepo.GetAll(page, pageSize)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get cameras: %w", err)
	}

	// Enrich semua cameras dengan stream URLs
	s.enrichCamerasWithStreamURLs(cameras)

	return cameras, meta, nil
}

func (s *cameraService) Update(id string, req *models.UpdateCameraRequest) (*models.Camera, error) {
	camera, err := s.cameraRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("camera not found: %w", err)
	}

	// Update fields
	if req.Name != "" {
		camera.Name = req.Name
	}
	if req.Description != "" {
		camera.Description = sql.NullString{String: req.Description, Valid: true}
	}
	if req.RTSPUrl != "" {
		camera.RTSPUrl = req.RTSPUrl
	}
	if req.Latitude != 0 {
		camera.Latitude = req.Latitude
	}
	if req.Longitude != 0 {
		camera.Longitude = req.Longitude
	}
	if req.Building != "" {
		camera.Building = sql.NullString{String: req.Building, Valid: true}
	}
	if req.Zone != "" {
		camera.Zone = sql.NullString{String: req.Zone, Valid: true}
	}
	if req.IPAddress != "" {
		camera.IPAddress = sql.NullString{String: req.IPAddress, Valid: true}
	}
	if req.Port != 0 {
		camera.Port = sql.NullInt64{Int64: int64(req.Port), Valid: true}
	}
	if req.Manufacturer != "" {
		camera.Manufacturer = sql.NullString{String: req.Manufacturer, Valid: true}
	}
	if req.Model != "" {
		camera.Model = sql.NullString{String: req.Model, Valid: true}
	}
	if req.Resolution != "" {
		camera.Resolution = sql.NullString{String: req.Resolution, Valid: true}
	}
	if req.FPS != 0 {
		camera.FPS = req.FPS
	}
	if len(req.Tags) > 0 {
		camera.Tags = req.Tags
	}
	if req.Status != "" {
		camera.Status = req.Status
	}

	if err := s.cameraRepo.Update(id, camera); err != nil {
		return nil, fmt.Errorf("failed to update camera: %w", err)
	}

	// Enrich dengan stream URLs
	s.enrichCameraWithStreamURLs(camera)

	return camera, nil
}

func (s *cameraService) Delete(id string) error {
	camera, err := s.cameraRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("camera not found: %w", err)
	}

	// Stop stream jika ada
	if camera.StreamID.Valid {
		s.rtspService.RemoveStream(camera.StreamID.String)
	}

	if err := s.cameraRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete camera: %w", err)
	}

	return nil
}

func (s *cameraService) GetByZone(zone string) ([]*models.Camera, error) {
	cameras, err := s.cameraRepo.GetByZone(zone)
	if err != nil {
		return nil, fmt.Errorf("failed to get cameras by zone: %w", err)
	}

	// Enrich semua cameras dengan stream URLs
	s.enrichCamerasWithStreamURLs(cameras)

	return cameras, nil
}

func (s *cameraService) GetNearby(lat, lng, radius float64) ([]*models.Camera, error) {
	cameras, err := s.cameraRepo.GetNearby(lat, lng, radius)
	if err != nil {
		return nil, fmt.Errorf("failed to get nearby cameras: %w", err)
	}

	// Enrich semua cameras dengan stream URLs
	s.enrichCamerasWithStreamURLs(cameras)

	return cameras, nil
}

func (s *cameraService) StartStream(id string) (*models.Camera, error) {
	camera, err := s.cameraRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("camera not found: %w", err)
	}

	// Add stream ke RTSPtoWeb jika belum ada
	if !camera.StreamID.Valid || camera.StreamID.String == "" {
		streamID, hlsURL, snapshotURL, err := s.rtspService.AddStream(
			camera.ID,
			camera.Name,
			camera.RTSPUrl,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to start stream: %w", err)
		}

		camera.StreamID = sql.NullString{String: streamID, Valid: true}
		camera.HLSUrl = hlsURL
		camera.SnapshotUrl = snapshotURL
		camera.Status = "READY"

		if err := s.cameraRepo.Update(id, camera); err != nil {
			return nil, fmt.Errorf("failed to update camera: %w", err)
		}
	}

	// Enrich dengan stream URLs
	s.enrichCameraWithStreamURLs(camera)

	return camera, nil
}

func (s *cameraService) StopStream(id string) error {
	camera, err := s.cameraRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("camera not found: %w", err)
	}

	if camera.StreamID.Valid {
		if err := s.rtspService.RemoveStream(camera.StreamID.String); err != nil {
			return fmt.Errorf("failed to stop stream: %w", err)
		}

		camera.StreamID = sql.NullString{Valid: false}
		camera.Status = "OFFLINE"

		if err := s.cameraRepo.Update(id, camera); err != nil {
			return fmt.Errorf("failed to update camera: %w", err)
		}
	}

	return nil
}
