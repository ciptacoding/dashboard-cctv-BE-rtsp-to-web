package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RTSPService interface {
	AddStream(cameraID, name, rtspURL string) (streamID, hlsURL, snapshotURL string, err error)
	RemoveStream(streamID string) error
	GetStreamStatus(streamID string) (string, error)
	GetHLSURL(streamID string) string
	GetSnapshotURL(streamID string) string
	GetSnapshotHash(streamID string) (string, error) // New: Get snapshot hash for frozen detection
}

type rtspService struct {
	apiURL        string
	publicBaseURL string
	username      string
	password      string
	httpClient    *http.Client
}

func NewRTSPService(apiURL, publicBaseURL, username, password string) RTSPService {
	// Configure HTTP client with proper timeouts and connection pooling
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
		DisableCompression:  false,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second, // 10 second timeout for all requests
	}

	return &rtspService{
		apiURL:        apiURL,
		publicBaseURL: publicBaseURL,
		username:      username,
		password:      password,
		httpClient:    httpClient,
	}
}

// GetHLSURL generates HLS URL for a stream
func (s *rtspService) GetHLSURL(streamID string) string {
	if streamID == "" {
		return ""
	}
	return fmt.Sprintf("%s/stream/%s/channel/0/hls/live/index.m3u8", s.publicBaseURL, streamID)
}

// GetSnapshotURL generates snapshot URL for a stream
func (s *rtspService) GetSnapshotURL(streamID string) string {
	if streamID == "" {
		return ""
	}
	return fmt.Sprintf("%s/stream/%s/channel/0/jpeg", s.publicBaseURL, streamID)
}

func (s *rtspService) AddStream(cameraID, name, rtspURL string) (streamID, hlsURL, snapshotURL string, err error) {
	payload := map[string]interface{}{
		"name": name,
		"channels": map[string]interface{}{
			"0": map[string]interface{}{
				"url":        rtspURL,
				"on_demand":  false, // Always running (instant playback)
				"persistent": true,  // Auto-reconnect on failure
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/stream/%s/add", s.apiURL, cameraID), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(s.username, s.password)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", "", "", fmt.Errorf("RTSPtoWeb API returned status %d: %s", resp.StatusCode, string(body))
	}

	streamID = cameraID
	hlsURL = s.GetHLSURL(streamID)
	snapshotURL = s.GetSnapshotURL(streamID)

	return streamID, hlsURL, snapshotURL, nil
}

func (s *rtspService) RemoveStream(streamID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/stream/%s/delete", s.apiURL, streamID), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(s.username, s.password)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("RTSPtoWeb API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *rtspService) GetStreamStatus(streamID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/stream/%s/info", s.apiURL, streamID), nil)
	if err != nil {
		return "ERROR", fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(s.username, s.password)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "ERROR", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Try to parse response to get actual status
		var info map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&info); err == nil {
			if channels, ok := info["channels"].(map[string]interface{}); ok {
				if ch0, ok := channels["0"].(map[string]interface{}); ok {
					if status, ok := ch0["status"].(string); ok {
						return status, nil
					}
				}
			}
		}
		return "READY", nil
	}

	return "OFFLINE", nil
}

// GetSnapshotHash gets the MD5 hash of the current snapshot for frozen detection
func (s *rtspService) GetSnapshotHash(streamID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	snapshotURL := s.GetSnapshotURL(streamID)
	if snapshotURL == "" {
		return "", fmt.Errorf("invalid stream ID")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", snapshotURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(s.username, s.password)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch snapshot: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("snapshot returned status: %d", resp.StatusCode)
	}

	// Read snapshot data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read snapshot: %w", err)
	}

	// Calculate MD5 hash
	hash := fmt.Sprintf("%x", md5.Sum(data))
	return hash, nil
}
