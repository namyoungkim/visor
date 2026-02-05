package auth

import (
	"testing"
)

func TestParseCredentialJSON_NestedEnvelope(t *testing.T) {
	data := []byte(`{"claudeAiOauth":{"accessToken":"tok123","refreshToken":"ref456","expiresAt":999,"subscriptionType":"pro","rateLimitTier":"default"}}`)
	creds, err := parseCredentialJSON(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds.AccessToken != "tok123" {
		t.Errorf("AccessToken = %q, want %q", creds.AccessToken, "tok123")
	}
	if creds.SubscriptionType != "pro" {
		t.Errorf("SubscriptionType = %q, want %q", creds.SubscriptionType, "pro")
	}
	if creds.RateLimitTier != "default" {
		t.Errorf("RateLimitTier = %q, want %q", creds.RateLimitTier, "default")
	}
}

func TestParseCredentialJSON_DirectFormat(t *testing.T) {
	data := []byte(`{"accessToken":"direct-tok","refreshToken":"ref","expiresAt":100}`)
	creds, err := parseCredentialJSON(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds.AccessToken != "direct-tok" {
		t.Errorf("AccessToken = %q, want %q", creds.AccessToken, "direct-tok")
	}
}

func TestParseCredentialJSON_NestedEmptyAccessToken(t *testing.T) {
	// Nested format with empty accessToken should fall through to direct
	data := []byte(`{"claudeAiOauth":{"accessToken":"","refreshToken":"ref"},"accessToken":"fallback-tok"}`)
	creds, err := parseCredentialJSON(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds.AccessToken != "fallback-tok" {
		t.Errorf("AccessToken = %q, want %q", creds.AccessToken, "fallback-tok")
	}
}

func TestParseCredentialJSON_MalformedJSON(t *testing.T) {
	data := []byte(`{invalid json`)
	_, err := parseCredentialJSON(data)
	if err == nil {
		t.Error("expected error for malformed JSON")
	}
}

func TestParseCredentialJSON_EmptyObject(t *testing.T) {
	data := []byte(`{}`)
	_, err := parseCredentialJSON(data)
	if err != ErrNoCredentials {
		t.Errorf("expected ErrNoCredentials, got %v", err)
	}
}

func TestParseCredentialJSON_NilClaudeAiOauth(t *testing.T) {
	data := []byte(`{"claudeAiOauth":null,"accessToken":"tok"}`)
	creds, err := parseCredentialJSON(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds.AccessToken != "tok" {
		t.Errorf("AccessToken = %q, want %q", creds.AccessToken, "tok")
	}
}
