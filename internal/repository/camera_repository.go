package repository

import (
	"database/sql"
	"fmt"

	"cctv-monitoring-backend/internal/models"

	"github.com/lib/pq"
)

// CameraRepository adalah interface untuk operasi database camera
type CameraRepository interface {
	Create(camera *models.Camera, createdBy string) error
	GetByID(id string) (*models.Camera, error)
	GetAll(limit, offset int) ([]*models.Camera, int64, error)
	Update(id string, camera *models.Camera) error
	Delete(id string) error
	UpdateStreamID(cameraID, streamID string) error
	UpdateStatus(cameraID, status string) error
	GetByZone(zone string) ([]*models.Camera, error)
	GetNearby(lat, lng, radiusKm float64) ([]*models.Camera, error)
}

type cameraRepository struct {
	db *sql.DB
}

// NewCameraRepository membuat instance baru dari CameraRepository
func NewCameraRepository(db *sql.DB) CameraRepository {
	return &cameraRepository{db: db}
}

// Create membuat camera baru di database
func (r *cameraRepository) Create(camera *models.Camera, createdBy string) error {
	query := `
		INSERT INTO cameras (
			name, description, rtsp_url, latitude, longitude,
			building, zone, ip_address, port, manufacturer, model,
			resolution, fps, tags, status, created_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		camera.Name,
		nullString(camera.Description),
		camera.RTSPUrl,
		camera.Latitude,
		camera.Longitude,
		nullString(camera.Building),
		nullString(camera.Zone),
		nullString(camera.IPAddress),
		nullInt64(camera.Port),
		nullString(camera.Manufacturer),
		nullString(camera.Model),
		nullString(camera.Resolution),
		camera.FPS,
		pq.Array(camera.Tags),
		camera.Status,
		createdBy,
	).Scan(&camera.ID, &camera.CreatedAt, &camera.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create camera: %w", err)
	}

	return nil
}

// GetByID mencari camera berdasarkan ID
func (r *cameraRepository) GetByID(id string) (*models.Camera, error) {
	query := `
		SELECT 
			id, name, description, rtsp_url, stream_id, latitude, longitude,
			building, zone, ip_address, port, manufacturer, model, resolution,
			fps, tags, status, last_seen, is_active, created_by, created_at, updated_at
		FROM cameras
		WHERE id = $1
	`

	camera := &models.Camera{}
	err := r.db.QueryRow(query, id).Scan(
		&camera.ID,
		&camera.Name,
		&camera.Description,
		&camera.RTSPUrl,
		&camera.StreamID,
		&camera.Latitude,
		&camera.Longitude,
		&camera.Building,
		&camera.Zone,
		&camera.IPAddress,
		&camera.Port,
		&camera.Manufacturer,
		&camera.Model,
		&camera.Resolution,
		&camera.FPS,
		pq.Array(&camera.Tags),
		&camera.Status,
		&camera.LastSeen,
		&camera.IsActive,
		&camera.CreatedBy,
		&camera.CreatedAt,
		&camera.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("camera not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get camera: %w", err)
	}

	return camera, nil
}

// GetAll mengambil semua camera dengan pagination
func (r *cameraRepository) GetAll(limit, offset int) ([]*models.Camera, int64, error) {
	// Hitung total records
	var total int64
	countQuery := "SELECT COUNT(*) FROM cameras WHERE is_active = true"
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cameras: %w", err)
	}

	// Ambil data dengan pagination
	query := `
		SELECT 
			id, name, description, rtsp_url, stream_id, latitude, longitude,
			building, zone, ip_address, port, manufacturer, model, resolution,
			fps, tags, status, last_seen, is_active, created_by, created_at, updated_at
		FROM cameras
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get cameras: %w", err)
	}
	defer rows.Close()

	cameras := make([]*models.Camera, 0)
	for rows.Next() {
		camera := &models.Camera{}
		err := rows.Scan(
			&camera.ID,
			&camera.Name,
			&camera.Description,
			&camera.RTSPUrl,
			&camera.StreamID,
			&camera.Latitude,
			&camera.Longitude,
			&camera.Building,
			&camera.Zone,
			&camera.IPAddress,
			&camera.Port,
			&camera.Manufacturer,
			&camera.Model,
			&camera.Resolution,
			&camera.FPS,
			pq.Array(&camera.Tags),
			&camera.Status,
			&camera.LastSeen,
			&camera.IsActive,
			&camera.CreatedBy,
			&camera.CreatedAt,
			&camera.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan camera: %w", err)
		}
		cameras = append(cameras, camera)
	}

	return cameras, total, nil
}

// Update mengupdate data camera
func (r *cameraRepository) Update(id string, camera *models.Camera) error {
	query := `
		UPDATE cameras
		SET name = $1, description = $2, rtsp_url = $3, latitude = $4, longitude = $5,
			building = $6, zone = $7, ip_address = $8, port = $9, manufacturer = $10,
			model = $11, resolution = $12, fps = $13, tags = $14, status = $15,
			is_active = $16, updated_at = NOW()
		WHERE id = $17
	`

	result, err := r.db.Exec(
		query,
		camera.Name,
		nullString(camera.Description),
		camera.RTSPUrl,
		camera.Latitude,
		camera.Longitude,
		nullString(camera.Building),
		nullString(camera.Zone),
		nullString(camera.IPAddress),
		nullInt64(camera.Port),
		nullString(camera.Manufacturer),
		nullString(camera.Model),
		nullString(camera.Resolution),
		camera.FPS,
		pq.Array(camera.Tags),
		camera.Status,
		camera.IsActive,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update camera: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("camera not found")
	}

	return nil
}

// Delete menghapus camera (soft delete)
func (r *cameraRepository) Delete(id string) error {
	query := "UPDATE cameras SET is_active = false, updated_at = NOW() WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete camera: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("camera not found")
	}

	return nil
}

// UpdateStreamID mengupdate stream ID dari RTSPtoWeb
func (r *cameraRepository) UpdateStreamID(cameraID, streamID string) error {
	query := "UPDATE cameras SET stream_id = $1, updated_at = NOW() WHERE id = $2"
	_, err := r.db.Exec(query, streamID, cameraID)
	return err
}

// UpdateStatus mengupdate status camera
func (r *cameraRepository) UpdateStatus(cameraID, status string) error {
	query := "UPDATE cameras SET status = $1, last_seen = NOW(), updated_at = NOW() WHERE id = $2"
	_, err := r.db.Exec(query, status, cameraID)
	return err
}

// GetByZone mengambil camera berdasarkan zone
func (r *cameraRepository) GetByZone(zone string) ([]*models.Camera, error) {
	query := `
		SELECT 
			id, name, description, rtsp_url, stream_id, latitude, longitude,
			building, zone, ip_address, port, manufacturer, model, resolution,
			fps, tags, status, last_seen, is_active, created_by, created_at, updated_at
		FROM cameras
		WHERE zone = $1 AND is_active = true
		ORDER BY name
	`

	rows, err := r.db.Query(query, zone)
	if err != nil {
		return nil, fmt.Errorf("failed to get cameras by zone: %w", err)
	}
	defer rows.Close()

	cameras := make([]*models.Camera, 0)
	for rows.Next() {
		camera := &models.Camera{}
		err := rows.Scan(
			&camera.ID,
			&camera.Name,
			&camera.Description,
			&camera.RTSPUrl,
			&camera.StreamID,
			&camera.Latitude,
			&camera.Longitude,
			&camera.Building,
			&camera.Zone,
			&camera.IPAddress,
			&camera.Port,
			&camera.Manufacturer,
			&camera.Model,
			&camera.Resolution,
			&camera.FPS,
			pq.Array(&camera.Tags),
			&camera.Status,
			&camera.LastSeen,
			&camera.IsActive,
			&camera.CreatedBy,
			&camera.CreatedAt,
			&camera.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan camera: %w", err)
		}
		cameras = append(cameras, camera)
	}

	return cameras, nil
}

// GetNearby mengambil camera dalam radius tertentu (dalam kilometer)
func (r *cameraRepository) GetNearby(lat, lng, radiusKm float64) ([]*models.Camera, error) {
	// Menggunakan earthdistance extension untuk query geospatial
	query := `
		SELECT 
			id, name, description, rtsp_url, stream_id, latitude, longitude,
			building, zone, ip_address, port, manufacturer, model, resolution,
			fps, tags, status, last_seen, is_active, created_by, created_at, updated_at,
			earth_distance(ll_to_earth($1, $2), ll_to_earth(latitude, longitude)) as distance
		FROM cameras
		WHERE is_active = true
			AND earth_distance(ll_to_earth($1, $2), ll_to_earth(latitude, longitude)) <= $3
		ORDER BY distance
	`

	radiusMeters := radiusKm * 1000 // Convert km to meters

	rows, err := r.db.Query(query, lat, lng, radiusMeters)
	if err != nil {
		return nil, fmt.Errorf("failed to get nearby cameras: %w", err)
	}
	defer rows.Close()

	cameras := make([]*models.Camera, 0)
	for rows.Next() {
		camera := &models.Camera{}
		var distance float64
		err := rows.Scan(
			&camera.ID,
			&camera.Name,
			&camera.Description,
			&camera.RTSPUrl,
			&camera.StreamID,
			&camera.Latitude,
			&camera.Longitude,
			&camera.Building,
			&camera.Zone,
			&camera.IPAddress,
			&camera.Port,
			&camera.Manufacturer,
			&camera.Model,
			&camera.Resolution,
			&camera.FPS,
			pq.Array(&camera.Tags),
			&camera.Status,
			&camera.LastSeen,
			&camera.IsActive,
			&camera.CreatedBy,
			&camera.CreatedAt,
			&camera.UpdatedAt,
			&distance,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan camera: %w", err)
		}
		cameras = append(cameras, camera)
	}

	return cameras, nil
}

// Helper functions untuk handle NULL values
func nullString(s sql.NullString) interface{} {
	if s.Valid {
		return s.String
	}
	return nil
}

func nullInt64(i sql.NullInt64) interface{} {
	if i.Valid {
		return i.Int64
	}
	return nil
}
