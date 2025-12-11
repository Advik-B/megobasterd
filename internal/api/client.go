package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/go-resty/resty/v2"
)

const (
	// APIURL is the MEGA API endpoint
	APIURL = "https://g.api.mega.co.nz"
	
	// UserAgent for HTTP requests
	UserAgent = "MegaBasterd-Go/1.0"
)

var (
	// ErrInvalidResponse is returned when API response is invalid
	ErrInvalidResponse = errors.New("invalid API response")
	
	// ErrUnauthorized is returned when session is invalid
	ErrUnauthorized = errors.New("unauthorized")
)

// MegaClient represents a MEGA API client
type MegaClient struct {
	client    *resty.Client
	sessionID string
	seqNo     int64
}

// NewMegaClient creates a new MEGA API client
func NewMegaClient() *MegaClient {
	client := resty.New().
		SetHeader("User-Agent", UserAgent).
		SetHeader("Content-Type", "application/json")

	return &MegaClient{
		client: client,
		seqNo:  0,
	}
}

// Request sends a request to MEGA API
func (m *MegaClient) Request(ctx context.Context, commands []map[string]interface{}) ([]interface{}, error) {
	seqNo := atomic.AddInt64(&m.seqNo, 1)
	
	url := fmt.Sprintf("%s/cs?id=%d", APIURL, seqNo)
	if m.sessionID != "" {
		url += fmt.Sprintf("&sid=%s", m.sessionID)
	}

	resp, err := m.client.R().
		SetContext(ctx).
		SetBody(commands).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode(), resp.Status())
	}

	var result []interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for MEGA errors
	if len(result) > 0 {
		if errCode, ok := result[0].(float64); ok && errCode < 0 {
			return nil, m.parseMegaError(int(errCode))
		}
	}

	return result, nil
}

// Login authenticates with MEGA
func (m *MegaClient) Login(ctx context.Context, email, password string) error {
	// This is a simplified version - actual implementation would need:
	// 1. Prepare user hash
	// 2. Send login command
	// 3. Process RSA keys
	// 4. Set session ID
	
	// For now, return not implemented
	return errors.New("login not yet implemented - see Phase 3 for full implementation")
}

// GetAccountInfo retrieves account information
func (m *MegaClient) GetAccountInfo(ctx context.Context) (map[string]interface{}, error) {
	commands := []map[string]interface{}{
		{"a": "uq"},
	}

	result, err := m.Request(ctx, commands)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, ErrInvalidResponse
	}

	accountInfo, ok := result[0].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidResponse
	}

	return accountInfo, nil
}

// GetDownloadURL retrieves download URL for a file
func (m *MegaClient) GetDownloadURL(ctx context.Context, fileID, key string) (string, error) {
	commands := []map[string]interface{}{
		{
			"a": "g",
			"g": 1,
			"p": fileID,
		},
	}

	result, err := m.Request(ctx, commands)
	if err != nil {
		return "", err
	}

	if len(result) == 0 {
		return "", ErrInvalidResponse
	}

	response, ok := result[0].(map[string]interface{})
	if !ok {
		return "", ErrInvalidResponse
	}

	url, ok := response["g"].(string)
	if !ok {
		return "", ErrInvalidResponse
	}

	return url, nil
}

// parseMegaError converts MEGA error codes to Go errors
func (m *MegaClient) parseMegaError(code int) error {
	switch code {
	case -1:
		return errors.New("internal error")
	case -2:
		return errors.New("invalid argument")
	case -3:
		return errors.New("request failed, retrying")
	case -4:
		return errors.New("rate limit exceeded")
	case -6:
		return errors.New("file not found")
	case -9:
		return errors.New("file already exists")
	case -15:
		return ErrUnauthorized
	case -16:
		return errors.New("user blocked")
	case -17:
		return errors.New("quota exceeded")
	case -18:
		return errors.New("temporarily unavailable")
	default:
		return fmt.Errorf("MEGA error code: %d", code)
	}
}

// SetSessionID sets the session ID for authenticated requests
func (m *MegaClient) SetSessionID(sid string) {
	m.sessionID = sid
}

// GetSessionID returns the current session ID
func (m *MegaClient) GetSessionID() string {
	return m.sessionID
}
