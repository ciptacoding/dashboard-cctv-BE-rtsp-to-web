package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type RTSPService interface {
	AddStream(cameraID, name, rtspURL string) (streamID, hlsURL, snapshotURL string, err error)
	RemoveStream(streamID string) error
	GetStreamStatus(streamID string) (string, error)
	GetHLSURL(streamID string) string      // NEW
	GetSnapshotURL(streamID string) string // NEW
}

type rtspService struct {
	apiURL        string
	publicBaseURL string
	username      string
	password      string
	httpClient    *http.Client
}

func NewRTSPService(apiURL, publicBaseURL, username, password string) RTSPService {
	return &rtspService{
		apiURL:        apiURL,
		publicBaseURL: publicBaseURL,
		username:      username,
		password:      password,
		httpClient:    &http.Client{},
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
				"url": rtspURL,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/stream/%s/add", s.apiURL, cameraID), bytes.NewBuffer(jsonData))
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
		return "", "", "", fmt.Errorf("RTSPtoWeb API returned status: %d", resp.StatusCode)
	}

	streamID = cameraID
	hlsURL = s.GetHLSURL(streamID)
	snapshotURL = s.GetSnapshotURL(streamID)

	return streamID, hlsURL, snapshotURL, nil
}

func (s *rtspService) RemoveStream(streamID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/stream/%s/delete", s.apiURL, streamID), nil)
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
		return fmt.Errorf("RTSPtoWeb API returned status: %d", resp.StatusCode)
	}

	return nil
}

func (s *rtspService) GetStreamStatus(streamID string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/stream/%s/info", s.apiURL, streamID), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(s.username, s.password)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "READY", nil
	}

	return "OFFLINE", nil
}
