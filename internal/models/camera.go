package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	
	// Status information (not stored in DB, generated dynamically)
	StatusMessage string `json:"status_message,omitempty"` // Human readable status message
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

// GetStatusMessage returns a human-readable message based on camera status
func GetStatusMessage(status string, hasStream bool, lastSeen sql.NullTime) string {
	switch status {
	case "ONLINE", "READY":
		if hasStream {
			return "Camera is online and streaming"
		}
		return "Camera is online but stream not started"
	case "OFFLINE":
		if lastSeen.Valid {
			timeSince := time.Since(lastSeen.Time)
			if timeSince < 1*time.Minute {
				return "Camera is offline (just disconnected)"
			} else if timeSince < 5*time.Minute {
				return "Camera is offline (disconnected " + formatDuration(timeSince) + " ago)"
			} else {
				return "Camera is offline (disconnected " + formatDuration(timeSince) + " ago). Attempting to reconnect..."
			}
		}
		return "Camera is offline. Attempting to reconnect..."
	case "ERROR":
		return "Camera encountered an error. Please check the connection"
	case "FROZEN":
		return "Camera stream appears frozen. Refreshing..."
	default:
		if !hasStream {
			return "Camera stream not started"
		}
		return "Camera status unknown"
	}
}

// formatDuration formats duration to human readable string
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f seconds", d.Seconds())
	} else if d < time.Hour {
		minutes := int(d.Minutes())
		return fmt.Sprintf("%d minute(s)", minutes)
	} else {
		hours := int(d.Hours())
		return fmt.Sprintf("%d hour(s)", hours)
	}
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

// CameraPreview adalah struktur untuk preview video kamera
type CameraPreview struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Status        string       `json:"status"`
	StatusMessage string       `json:"status_message"`
	HLSUrl        string       `json:"hls_url,omitempty"`
	SnapshotUrl   string       `json:"snapshot_url,omitempty"`
	HasStream     bool         `json:"has_stream"`
	LastSeen      sql.NullTime `json:"-"`
}

// StreamErrorReport adalah request untuk report stream error
type StreamErrorReport struct {
	ErrorType string `json:"error_type"` // "timeout", "hls_error", "network_error", etc.
	Message   string `json:"message,omitempty"`
}
