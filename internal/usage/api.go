package usage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/namyoungkim/visor/internal/auth"
)

const (
	// API endpoints
	baseURL = "https://api.claude.ai"

	// Default timeout for API requests
	defaultTimeout = 5 * time.Second
)

// Client is an OAuth API client for Claude Pro usage data.
type Client struct {
	httpClient *http.Client
	authProvider auth.CredentialProvider
}

// NewClient creates a new usage API client.
func NewClient(provider auth.CredentialProvider) *Client {
	if provider == nil {
		provider = auth.DefaultProvider()
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		authProvider: provider,
	}
}

// Limits represents usage limits and current utilization.
type Limits struct {
	FiveHour FiveHourLimit `json:"fiveHour"`
	SevenDay SevenDayLimit `json:"sevenDay"`
}

// FiveHourLimit represents the 5-hour rate limit.
type FiveHourLimit struct {
	Utilization float64   `json:"utilization"` // 0-100%
	ResetsAt    time.Time `json:"resetsAt"`
	Remaining   int       `json:"remaining,omitempty"` // Messages remaining
	Total       int       `json:"total,omitempty"`     // Total messages allowed
}

// SevenDayLimit represents the 7-day rate limit.
type SevenDayLimit struct {
	Utilization float64   `json:"utilization"` // 0-100%
	ResetsAt    time.Time `json:"resetsAt"`
	Remaining   int       `json:"remaining,omitempty"`
	Total       int       `json:"total,omitempty"`
}

// apiResponse represents the raw API response.
type apiResponse struct {
	FiveHourBlock struct {
		UtilizationPct float64 `json:"utilization_pct"`
		ResetsAt       string  `json:"resets_at"`
		MessagesLeft   int     `json:"messages_left,omitempty"`
		TotalMessages  int     `json:"total_messages,omitempty"`
	} `json:"five_hour_block"`
	SevenDayBlock struct {
		UtilizationPct float64 `json:"utilization_pct"`
		ResetsAt       string  `json:"resets_at"`
		MessagesLeft   int     `json:"messages_left,omitempty"`
		TotalMessages  int     `json:"total_messages,omitempty"`
	} `json:"seven_day_block"`
}

// GetLimits retrieves current usage limits from the API.
func (c *Client) GetLimits() (*Limits, error) {
	creds, err := c.authProvider.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	req, err := http.NewRequest("GET", baseURL+"/api/organizations/default/usage", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "visor")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, auth.ErrExpiredCredentials
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return c.parseResponse(&apiResp)
}

// parseResponse converts the API response to our Limits type.
func (c *Client) parseResponse(resp *apiResponse) (*Limits, error) {
	limits := &Limits{}

	// Parse 5-hour block
	limits.FiveHour.Utilization = resp.FiveHourBlock.UtilizationPct
	limits.FiveHour.Remaining = resp.FiveHourBlock.MessagesLeft
	limits.FiveHour.Total = resp.FiveHourBlock.TotalMessages
	if resp.FiveHourBlock.ResetsAt != "" {
		if t, err := time.Parse(time.RFC3339, resp.FiveHourBlock.ResetsAt); err == nil {
			limits.FiveHour.ResetsAt = t
		}
	}

	// Parse 7-day block
	limits.SevenDay.Utilization = resp.SevenDayBlock.UtilizationPct
	limits.SevenDay.Remaining = resp.SevenDayBlock.MessagesLeft
	limits.SevenDay.Total = resp.SevenDayBlock.TotalMessages
	if resp.SevenDayBlock.ResetsAt != "" {
		if t, err := time.Parse(time.RFC3339, resp.SevenDayBlock.ResetsAt); err == nil {
			limits.SevenDay.ResetsAt = t
		}
	}

	return limits, nil
}

// FiveHourRemaining returns the remaining time in the 5-hour block.
func (l *Limits) FiveHourRemaining() time.Duration {
	if l.FiveHour.ResetsAt.IsZero() {
		return 0
	}
	remaining := time.Until(l.FiveHour.ResetsAt)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// SevenDayRemaining returns the remaining time in the 7-day block.
func (l *Limits) SevenDayRemaining() time.Duration {
	if l.SevenDay.ResetsAt.IsZero() {
		return 0
	}
	remaining := time.Until(l.SevenDay.ResetsAt)
	if remaining < 0 {
		return 0
	}
	return remaining
}
