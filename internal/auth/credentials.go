package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Credentials represents OAuth credentials for Claude subscriptions.
type Credentials struct {
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	ExpiresAt        int64  `json:"expiresAt"`
	SubscriptionType string `json:"subscriptionType"` // "pro", "max", "team", etc.
	RateLimitTier    string `json:"rateLimitTier"`    // e.g. "default_claude_max_5x"
}

// keychainEnvelope wraps credentials stored in macOS Keychain.
// Claude Code stores credentials under {"claudeAiOauth": {...}}.
type keychainEnvelope struct {
	ClaudeAiOauth *Credentials `json:"claudeAiOauth"`
}

// parseCredentialJSON tries to unmarshal JSON as either direct Credentials
// or as a keychainEnvelope with nested claudeAiOauth key.
func parseCredentialJSON(data []byte) (*Credentials, error) {
	// Try nested format first: {"claudeAiOauth": {...}}
	var envelope keychainEnvelope
	if err := json.Unmarshal(data, &envelope); err == nil && envelope.ClaudeAiOauth != nil && envelope.ClaudeAiOauth.AccessToken != "" {
		return envelope.ClaudeAiOauth, nil
	}

	// Try direct format: {"accessToken": "..."}
	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}
	if creds.AccessToken == "" {
		return nil, ErrNoCredentials
	}
	return &creds, nil
}

// CredentialProvider defines the interface for credential access.
type CredentialProvider interface {
	// Get returns OAuth credentials.
	Get() (*Credentials, error)
}

// FileCredentialProvider reads credentials from a JSON file.
type FileCredentialProvider struct {
	Path string
}

// NewFileProvider creates a credential provider that reads from the default location.
func NewFileProvider() *FileCredentialProvider {
	home, err := os.UserHomeDir()
	if err != nil {
		return &FileCredentialProvider{}
	}
	return &FileCredentialProvider{
		Path: filepath.Join(home, ".claude", "credentials.json"),
	}
}

// Get reads credentials from the file.
func (p *FileCredentialProvider) Get() (*Credentials, error) {
	if p.Path == "" {
		return nil, ErrNoCredentials
	}

	data, err := os.ReadFile(p.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoCredentials
		}
		return nil, err
	}

	return parseCredentialJSON(data)
}

// DefaultProvider returns the default credential provider for the current platform.
func DefaultProvider() CredentialProvider {
	// First try file-based credentials
	fileProvider := NewFileProvider()
	if _, err := fileProvider.Get(); err == nil {
		return fileProvider
	}

	// Fall back to platform-specific provider (keychain, etc.)
	return platformProvider()
}
