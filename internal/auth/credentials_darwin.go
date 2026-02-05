//go:build darwin

package auth

import (
	"bytes"
	"context"
	"os"
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

// keychainEntries defines the keychain service/account pairs to try, in order.
// Claude Code has changed its keychain entry format over time.
var keychainEntries = []struct {
	service string
	account string
}{
	{"Claude Code-credentials", ""}, // Current: account = OS username (auto-detected)
	{"claude.ai", "oauth"},          // Legacy
}

// Get retrieves credentials from macOS Keychain.
func (p *keychainProvider) Get() (*Credentials, error) {
	for _, entry := range keychainEntries {
		creds, err := readKeychain(entry.service, entry.account)
		if err == nil {
			return creds, nil
		}
	}
	return nil, ErrNoCredentials
}

// readKeychain reads a credential from macOS Keychain.
// If account is empty, the current OS username is used.
func readKeychain(service, account string) (*Credentials, error) {
	if account == "" {
		account = currentUsername()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "security", "find-generic-password",
		"-s", service,
		"-a", account,
		"-w",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
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

	password := strings.TrimSpace(stdout.String())
	if password == "" {
		return nil, ErrNoCredentials
	}

	return parseCredentialJSON([]byte(password))
}

// currentUsername returns the current OS username.
func currentUsername() string {
	if u := os.Getenv("USER"); u != "" {
		return u
	}
	return "unknown"
}
