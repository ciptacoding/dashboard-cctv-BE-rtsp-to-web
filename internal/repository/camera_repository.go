package repository

import (
	"cctv-monitoring-backend/internal/models"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type CameraRepository interface {
	Create(camera *models.Camera, userID string) error
	GetByID(id string) (*models.Camera, error)
	GetAll(page, pageSize int) ([]*models.Camera, *models.PaginationMeta, error)
	Update(id string, camera *models.Camera) error
	Delete(id string) error
	GetByZone(zone string) ([]*models.Camera, error)
	GetNearby(lat, lng, radius float64) ([]*models.Camera, error)
}

type cameraRepository struct {
	db *sql.DB
}

func NewCameraRepository(db *sql.DB) CameraRepository {
	return &cameraRepository{db: db}
}

func (r *cameraRepository) Create(camera *models.Camera, userID string) error {
	query := `
		INSERT INTO cameras (
			id, name, description, rtsp_url, stream_id,
			latitude, longitude, building, zone,
			ip_address, port, manufacturer, model, resolution, fps,
			tags, status, is_active, created_by,
			created_at, updated_at
		) VALUES (
			uuid_generate_v4(), $1, $2, $3, $4,
			$5, $6, $7, $8,
			$9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18,
			NOW(), NOW()
		) RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		camera.Name,
		camera.Description,
		camera.RTSPUrl,
		camera.StreamID,
		camera.Latitude,
		camera.Longitude,
		camera.Building,
		camera.Zone,
		camera.IPAddress,
		camera.Port,
		camera.Manufacturer,
		camera.Model,
		camera.Resolution,
		camera.FPS,
		pq.Array(camera.Tags),
		camera.Status,
		camera.IsActive,
		userID,
	).Scan(&camera.ID, &camera.CreatedAt, &camera.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create camera: %w", err)
	}

	return nil
}

func (r *cameraRepository) GetByID(id string) (*models.Camera, error) {
	query := `
		SELECT 
			id, name, description, rtsp_url, stream_id,
			latitude, longitude, building, zone,
			ip_address, port, manufacturer, model, resolution, fps,
			tags, status, last_seen, is_active, created_by,
			created_at, updated_at
		FROM cameras
		WHERE id = $1 AND is_active = true
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

func (r *cameraRepository) GetAll(page, pageSize int) ([]*models.Camera, *models.PaginationMeta, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalItems int64
	countQuery := "SELECT COUNT(*) FROM cameras WHERE is_active = true"
	if err := r.db.QueryRow(countQuery).Scan(&totalItems); err != nil {
		return nil, nil, fmt.Errorf("failed to count cameras: %w", err)
	}

	// Calculate total pages
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		totalPages++
	}

	// Get cameras
	query := `
		SELECT 
			id, name, description, rtsp_url, stream_id,
			latitude, longitude, building, zone,
			ip_address, port, manufacturer, model, resolution, fps,
			tags, status, last_seen, is_active, created_by,
			created_at, updated_at
		FROM cameras
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get cameras: %w", err)
	}
	defer rows.Close()

	cameras := []*models.Camera{}
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
			return nil, nil, fmt.Errorf("failed to scan camera: %w", err)
		}
		cameras = append(cameras, camera)
	}

	meta := &models.PaginationMeta{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	return cameras, meta, nil
}

func (r *cameraRepository) Update(id string, camera *models.Camera) error {
	query := `
		UPDATE cameras SET
			name = $1,
			description = $2,
			rtsp_url = $3,
			stream_id = $4,
			latitude = $5,
			longitude = $6,
			building = $7,
			zone = $8,
			ip_address = $9,
			port = $10,
			manufacturer = $11,
			model = $12,
			resolution = $13,
			fps = $14,
			tags = $15,
			status = $16,
			last_seen = $17,
			is_active = $18,
			updated_at = NOW()
		WHERE id = $19
	`

	_, err := r.db.Exec(
		query,
		camera.Name,
		camera.Description,
		camera.RTSPUrl,
		camera.StreamID,
		camera.Latitude,
		camera.Longitude,
		camera.Building,
		camera.Zone,
		camera.IPAddress,
		camera.Port,
		camera.Manufacturer,
		camera.Model,
		camera.Resolution,
		camera.FPS,
		pq.Array(camera.Tags),
		camera.Status,
		camera.LastSeen,
		camera.IsActive,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update camera: %w", err)
	}

	return nil
}

func (r *cameraRepository) Delete(id string) error {
	query := "UPDATE cameras SET is_active = false, updated_at = NOW() WHERE id = $1"

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete camera: %w", err)
	}

	return nil
}

func (r *cameraRepository) GetByZone(zone string) ([]*models.Camera, error) {
	query := `
		SELECT 
			id, name, description, rtsp_url, stream_id,
			latitude, longitude, building, zone,
			ip_address, port, manufacturer, model, resolution, fps,
			tags, status, last_seen, is_active, created_by,
			created_at, updated_at
		FROM cameras
		WHERE zone = $1 AND is_active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, zone)
	if err != nil {
		return nil, fmt.Errorf("failed to get cameras by zone: %w", err)
	}
	defer rows.Close()

	cameras := []*models.Camera{}
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

func (r *cameraRepository) GetNearby(lat, lng, radius float64) ([]*models.Camera, error) {
	query := `
		SELECT 
			id, name, description, rtsp_url, stream_id,
			latitude, longitude, building, zone,
			ip_address, port, manufacturer, model, resolution, fps,
			tags, status, last_seen, is_active, created_by,
			created_at, updated_at,
			earth_distance(
				ll_to_earth(latitude, longitude),
				ll_to_earth($1, $2)
			) / 1000 as distance_km
		FROM cameras
		WHERE is_active = true
		AND earth_box(ll_to_earth($1, $2), $3 * 1000) @> ll_to_earth(latitude, longitude)
		ORDER BY distance_km ASC
	`

	rows, err := r.db.Query(query, lat, lng, radius)
	if err != nil {
		return nil, fmt.Errorf("failed to get nearby cameras: %w", err)
	}
	defer rows.Close()

	cameras := []*models.Camera{}
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
