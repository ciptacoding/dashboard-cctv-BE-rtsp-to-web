package models

import (
	"database/sql"
	"time"
)

// Camera merepresentasikan struktur data kamera CCTV
type Camera struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  sql.NullString `json:"description,omitempty"`
	RTSPUrl      string         `json:"rtsp_url"`
	StreamID     sql.NullString `json:"stream_id,omitempty"`
	Latitude     float64        `json:"latitude"`
	Longitude    float64        `json:"longitude"`
	Building     sql.NullString `json:"building,omitempty"`
	Zone         sql.NullString `json:"zone,omitempty"`
	IPAddress    sql.NullString `json:"ip_address,omitempty"`
	Port         sql.NullInt64  `json:"port,omitempty"`
	Manufacturer sql.NullString `json:"manufacturer,omitempty"`
	Model        sql.NullString `json:"model,omitempty"`
	Resolution   sql.NullString `json:"resolution,omitempty"`
	FPS          int            `json:"fps"`
	Tags         []string       `json:"tags"`
	Status       string         `json:"status"`
	LastSeen     sql.NullTime   `json:"last_seen,omitempty"`
	IsActive     bool           `json:"is_active"`
	CreatedBy    sql.NullString `json:"created_by,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// CreateCameraRequest adalah struktur untuk membuat kamera baru
type CreateCameraRequest struct {
	Name         string   `json:"name"`
	Description  string   `json:"description,omitempty"`
	RTSPUrl      string   `json:"rtsp_url"`
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	Building     string   `json:"building,omitempty"`
	Zone         string   `json:"zone,omitempty"`
	IPAddress    string   `json:"ip_address,omitempty"`
	Port         int      `json:"port,omitempty"`
	Manufacturer string   `json:"manufacturer,omitempty"`
	Model        string   `json:"model,omitempty"`
	Resolution   string   `json:"resolution,omitempty"`
	FPS          int      `json:"fps,omitempty"`
	Tags         []string `json:"tags,omitempty"`
}

// UpdateCameraRequest adalah struktur untuk update kamera
type UpdateCameraRequest struct {
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	RTSPUrl      string   `json:"rtsp_url,omitempty"`
	Latitude     float64  `json:"latitude,omitempty"`
	Longitude    float64  `json:"longitude,omitempty"`
	Building     string   `json:"building,omitempty"`
	Zone         string   `json:"zone,omitempty"`
	IPAddress    string   `json:"ip_address,omitempty"`
	Port         int      `json:"port,omitempty"`
	Manufacturer string   `json:"manufacturer,omitempty"`
	Model        string   `json:"model,omitempty"`
	Resolution   string   `json:"resolution,omitempty"`
	FPS          int      `json:"fps,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Status       string   `json:"status,omitempty"`
	IsActive     bool     `json:"is_active,omitempty"`
}

// CameraWithStream adalah camera dengan informasi stream URL
type CameraWithStream struct {
	Camera
	StreamURL   string `json:"stream_url,omitempty"`
	HLSUrl      string `json:"hls_url,omitempty"`
	WebRTCUrl   string `json:"webrtc_url,omitempty"`
	SnapshotURL string `json:"snapshot_url,omitempty"`
}
