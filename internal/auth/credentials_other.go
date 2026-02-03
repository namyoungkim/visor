//go:build !darwin && !linux

package auth

// platformProvider returns a fallback credential provider for other platforms.
func platformProvider() CredentialProvider {
	return NewFileProvider()
}
