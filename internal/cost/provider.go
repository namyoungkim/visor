package cost

import (
	"os"
	"path/filepath"
	"strings"
)

// Provider represents a Claude API provider.
type Provider string

const (
	ProviderUnknown   Provider = "unknown"
	ProviderAnthropic Provider = "anthropic" // Direct API (pay-per-token)
	ProviderClaudePro Provider = "claude_pro" // Claude Pro subscription
	ProviderAWS       Provider = "aws"        // AWS Bedrock
	ProviderGCP       Provider = "gcp"        // Google Cloud Vertex AI
)

// DetectProvider determines the provider from environment and config.
func DetectProvider() Provider {
	// Check environment variables
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		return ProviderAnthropic
	}

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" || os.Getenv("AWS_PROFILE") != "" {
		return ProviderAWS
	}

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" || os.Getenv("CLOUDSDK_CORE_PROJECT") != "" {
		return ProviderGCP
	}

	// Check for Claude Pro by looking for OAuth credentials
	home, err := os.UserHomeDir()
	if err == nil {
		credPath := filepath.Join(home, ".claude", "credentials.json")
		if _, err := os.Stat(credPath); err == nil {
			return ProviderClaudePro
		}
	}

	return ProviderUnknown
}

// IsPaidPerToken returns true if the provider charges per token.
func (p Provider) IsPaidPerToken() bool {
	switch p {
	case ProviderAnthropic, ProviderAWS, ProviderGCP:
		return true
	default:
		return false
	}
}

// IsSubscription returns true if the provider uses subscription-based billing.
func (p Provider) IsSubscription() bool {
	return p == ProviderClaudePro
}

// String returns the provider name.
func (p Provider) String() string {
	return string(p)
}

// GetProjectsDir returns the Claude projects directory.
func GetProjectsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".claude", "projects")
}

// GetSessionDir returns the directory for a specific session.
func GetSessionDir(projectPath string) string {
	if projectPath == "" {
		return ""
	}

	projectsDir := GetProjectsDir()
	if projectsDir == "" {
		return ""
	}

	// Convert project path to directory name
	// e.g., /Users/foo/bar → -Users-foo-bar
	safePath := strings.ReplaceAll(projectPath, "/", "-")
	safePath = strings.ReplaceAll(safePath, "\\", "-")
	if safePath[0] == '-' {
		safePath = safePath[1:]
	}

	return filepath.Join(projectsDir, safePath)
}
