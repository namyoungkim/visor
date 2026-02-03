//go:build linux

package auth

// linuxProvider is a placeholder for Linux credential access.
// On Linux, credentials are typically stored in the file system
// rather than a system keychain.
type linuxProvider struct{}

// platformProvider returns the Linux-specific credential provider.
func platformProvider() CredentialProvider {
	// On Linux, we rely on file-based credentials
	return NewFileProvider()
}
