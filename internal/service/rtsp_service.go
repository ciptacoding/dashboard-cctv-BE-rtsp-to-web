package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RTSPService adalah interface untuk integrasi dengan RTSPtoWeb
type RTSPService interface {
	AddStream(cameraID, rtspURL string) (string, error)
	RemoveStream(streamID string) error
	GetHLSURL(streamID string) string
	GetWebRTCURL(streamID string) string
	GetSnapshotURL(streamID string) string
	GetPublicBaseURL() string
}

type rtspService struct {
	baseURL       string
	publicBaseURL string
	username      string
	password      string
}

// NewRTSPService membuat instance baru dari RTSPService
func NewRTSPService(baseURL, publicBaseURL, username, password string) RTSPService {
	return &rtspService{
		baseURL:       baseURL,
		publicBaseURL: publicBaseURL,
		username:      username,
		password:      password,
	}
}

// StreamConfig adalah struktur konfigurasi stream untuk RTSPtoWeb
type StreamConfig struct {
	Name string            `json:"name"`
	URL  string            `json:"url"`
	On   map[string]string `json:"on_demand,omitempty"`
}

// AddStream menambahkan stream baru ke RTSPtoWeb
func (s *rtspService) AddStream(cameraID, rtspURL string) (string, error) {
	// Konfigurasi stream sesuai format RTSPtoWeb
	streamConfig := map[string]interface{}{
		"uuid": cameraID,
		"name": cameraID,
		"channels": map[string]interface{}{
			"0": map[string]interface{}{
				"url":       rtspURL,
				"on_demand": true,
				"debug":     false,
			},
		},
	}

	jsonData, err := json.Marshal(streamConfig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}

	// Kirim request ke RTSPtoWeb API
	url := fmt.Sprintf("%s/stream/%s/add", s.baseURL, cameraID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// RTSPtoWeb menggunakan application/json (bukan form-urlencoded)
	req.Header.Set("Content-Type", "application/json")

	// Add basic auth jika ada (dari UI: Basic ZGVtbzpkZW1v = demo:demo)
	if s.username != "" && s.password != "" {
		req.SetBasicAuth(s.username, s.password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to add stream, status: %d, response: %s", resp.StatusCode, string(body))
	}

	// Log success
	fmt.Printf("✓ Stream registered successfully: %s, response: %s\n", cameraID, string(body))

	// Return stream ID
	return cameraID, nil
}

// RemoveStream menghapus stream dari RTSPtoWeb
func (s *rtspService) RemoveStream(streamID string) error {
	// RTSPtoWeb menggunakan endpoint edit untuk disable stream
	config := map[string]interface{}{
		"name": streamID,
		"channels": map[string]interface{}{
			"0": map[string]interface{}{
				"name":      streamID,
				"url":       "",
				"on_demand": false,
				"status":    0,
			},
		},
	}

	jsonData, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	url := fmt.Sprintf("%s/stream/%s/edit", s.baseURL, streamID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add basic auth jika ada
	if s.username != "" && s.password != "" {
		req.SetBasicAuth(s.username, s.password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 404 adalah ok karena stream mungkin sudah tidak ada
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("failed to remove stream, status: %d", resp.StatusCode)
	}

	return nil
}

// GetHLSURL mengembalikan URL HLS untuk streaming
func (s *rtspService) GetHLSURL(streamID string) string {
	return fmt.Sprintf("%s/stream/%s/channel/0/hls/live/index.m3u8", s.publicBaseURL, streamID)
}

// GetWebRTCURL mengembalikan URL WebRTC untuk streaming
func (s *rtspService) GetWebRTCURL(streamID string) string {
	return fmt.Sprintf("%s/stream/%s/channel/0/webrtc", s.publicBaseURL, streamID)
}

// GetSnapshotURL mengembalikan URL untuk snapshot JPEG
func (s *rtspService) GetSnapshotURL(streamID string) string {
	return fmt.Sprintf("%s/stream/%s/channel/0/jpeg", s.publicBaseURL, streamID)
}

// GetPublicBaseURL mengembalikan public base URL
func (s *rtspService) GetPublicBaseURL() string {
	return s.publicBaseURL
}
