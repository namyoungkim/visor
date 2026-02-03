//go:build darwin

package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"strings"
	"time"
)

// keychainProvider reads credentials from macOS Keychain.
type keychainProvider struct{}

// platformProvider returns the macOS-specific credential provider.
func platformProvider() CredentialProvider {
	return &keychainProvider{}
}

// Get retrieves credentials from macOS Keychain.
func (p *keychainProvider) Get() (*Credentials, error) {
	// Set timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Try to read from Keychain using security command
	// The service name used by Claude Code
	cmd := exec.CommandContext(ctx, "security", "find-generic-password",
		"-s", "claude.ai",
		"-a", "oauth",
		"-w", // Output password only
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// Keychain access denied or item not found
		if strings.Contains(stderr.String(), "could not be found") ||
			strings.Contains(stderr.String(), "SecKeychainSearchCopyNext") {
			return nil, ErrNoCredentials
		}
		if strings.Contains(stderr.String(), "User interaction is not allowed") ||
			strings.Contains(stderr.String(), "User canceled") {
			return nil, ErrKeychainAccess
		}
		return nil, ErrNoCredentials
	}

	// Parse the JSON credential
	password := strings.TrimSpace(stdout.String())
	if password == "" {
		return nil, ErrNoCredentials
	}

	var creds Credentials
	if err := json.Unmarshal([]byte(password), &creds); err != nil {
		// Maybe it's just a token, not JSON
		creds.AccessToken = password
	}

	if creds.AccessToken == "" {
		return nil, ErrNoCredentials
	}

	return &creds, nil
}
