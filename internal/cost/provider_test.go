package cost

import (
	"os"
	"testing"
)

// clearProviderEnv saves and clears all provider-related env vars.
// Returns a cleanup function to restore them.
func clearProviderEnv(t *testing.T) {
	t.Helper()
	envVars := []string{
		"ANTHROPIC_API_KEY",
		"AWS_ACCESS_KEY_ID", "AWS_PROFILE",
		"GOOGLE_APPLICATION_CREDENTIALS", "CLOUDSDK_CORE_PROJECT",
	}
	saved := make(map[string]string)
	for _, key := range envVars {
		saved[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
	t.Cleanup(func() {
		for key, val := range saved {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	})
}

func TestDetectProvider_EnvVars(t *testing.T) {
	clearProviderEnv(t)

	tests := []struct {
		name     string
		envKey   string
		envValue string
		want     Provider
	}{
		{"Anthropic API key", "ANTHROPIC_API_KEY", "sk-test", ProviderAnthropic},
		{"AWS access key", "AWS_ACCESS_KEY_ID", "AKIA...", ProviderAWS},
		{"AWS profile", "AWS_PROFILE", "bedrock", ProviderAWS},
		{"GCP credentials", "GOOGLE_APPLICATION_CREDENTIALS", "/path/to/creds.json", ProviderGCP},
		{"GCP project", "CLOUDSDK_CORE_PROJECT", "my-project", ProviderGCP},
	}

	envVars := []string{
		"ANTHROPIC_API_KEY",
		"AWS_ACCESS_KEY_ID", "AWS_PROFILE",
		"GOOGLE_APPLICATION_CREDENTIALS", "CLOUDSDK_CORE_PROJECT",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResetProviderCache()
			for _, key := range envVars {
				os.Unsetenv(key)
			}

			os.Setenv(tt.envKey, tt.envValue)
			defer os.Unsetenv(tt.envKey)

			got := DetectProvider()
			if got != tt.want {
				t.Errorf("DetectProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectProvider_EnvPriority(t *testing.T) {
	clearProviderEnv(t)
	ResetProviderCache()

	os.Setenv("ANTHROPIC_API_KEY", "sk-test")
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIA...")

	got := DetectProvider()
	if got != ProviderAnthropic {
		t.Errorf("DetectProvider() = %v, want %v (Anthropic should take priority)", got, ProviderAnthropic)
	}
}

func TestDetectProvider_Unknown(t *testing.T) {
	clearProviderEnv(t)
	ResetProviderCache()

	// Note: This test may return ProviderClaudePro if credentials exist
	// in the file or macOS Keychain. We only assert it doesn't return
	// an env-based provider.
	got := DetectProvider()
	if got == ProviderAnthropic || got == ProviderAWS || got == ProviderGCP {
		t.Errorf("DetectProvider() = %v, want non-env provider (claude_pro or unknown)", got)
	}
}

func TestDetectProvider_CacheWorks(t *testing.T) {
	clearProviderEnv(t)
	ResetProviderCache()

	os.Setenv("ANTHROPIC_API_KEY", "sk-test")
	first := DetectProvider()
	if first != ProviderAnthropic {
		t.Fatalf("first call = %v, want %v", first, ProviderAnthropic)
	}

	// Change env — cached result should persist
	os.Unsetenv("ANTHROPIC_API_KEY")
	second := DetectProvider()
	if second != ProviderAnthropic {
		t.Errorf("cached call = %v, want %v (should be cached)", second, ProviderAnthropic)
	}
}
