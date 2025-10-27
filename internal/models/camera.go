package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Camera merepresentasikan struktur data kamera CCTV
// Camera merepresentasikan struktur data kamera CCTV
type Camera struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  sql.NullString `json:"-"`
	RTSPUrl      string         `json:"rtsp_url"`
	StreamID     sql.NullString `json:"-"`
	Latitude     float64        `json:"latitude"`
	Longitude    float64        `json:"longitude"`
	Building     sql.NullString `json:"-"`
	Zone         sql.NullString `json:"-"`
	IPAddress    sql.NullString `json:"-"`
	Port         sql.NullInt64  `json:"-"`
	Manufacturer sql.NullString `json:"-"`
	Model        sql.NullString `json:"-"`
	Resolution   sql.NullString `json:"-"`
	FPS          int            `json:"fps"`
	Tags         []string       `json:"tags"`
	Status       string         `json:"status"`
	LastSeen     sql.NullTime   `json:"-"`
	IsActive     bool           `json:"is_active"`
	CreatedBy    sql.NullString `json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`

	// Stream URLs (not stored in DB, generated dynamically)
	HLSUrl      string `json:"hls_url,omitempty"`      // NEW
	SnapshotUrl string `json:"snapshot_url,omitempty"` // NEW
}

// MarshalJSON custom JSON marshaling untuk Camera
func (c Camera) MarshalJSON() ([]byte, error) {
	type Alias Camera
	return json.Marshal(&struct {
		*Alias
		Description  string `json:"description,omitempty"`
		StreamID     string `json:"stream_id,omitempty"`
		Building     string `json:"building,omitempty"`
		Zone         string `json:"zone,omitempty"`
		IPAddress    string `json:"ip_address,omitempty"`
		Port         *int64 `json:"port,omitempty"`
		Manufacturer string `json:"manufacturer,omitempty"`
		Model        string `json:"model,omitempty"`
		Resolution   string `json:"resolution,omitempty"`
		LastSeen     string `json:"last_seen,omitempty"`
		CreatedBy    string `json:"created_by,omitempty"`
	}{
		Alias:        (*Alias)(&c),
		Description:  c.Description.String,
		StreamID:     c.StreamID.String,
		Building:     c.Building.String,
		Zone:         c.Zone.String,
		IPAddress:    c.IPAddress.String,
		Port:         nullInt64ToPtr(c.Port),
		Manufacturer: c.Manufacturer.String,
		Model:        c.Model.String,
		Resolution:   c.Resolution.String,
		LastSeen:     formatNullTime(c.LastSeen),
		CreatedBy:    c.CreatedBy.String,
	})
}

// Helper functions
func nullInt64ToPtr(ni sql.NullInt64) *int64 {
	if ni.Valid {
		return &ni.Int64
	}
	return nil
}

func formatNullTime(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format(time.RFC3339)
	}
	return ""
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
	Status       string   `json:"status,omitempty"` // NEW: tambahkan ini
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
