package input

import (
	"strings"
	"testing"
)

func TestParse_ValidJSON(t *testing.T) {
	jsonInput := `{
		"model": {"display_name": "Opus", "id": "claude-opus-4"},
		"cost": {"total_cost_usd": 0.05, "total_api_duration_ms": 1234},
		"context_window": {"used_percentage": 42.5},
		"workspace": {"lines_added": 10, "lines_removed": 5},
		"current_usage": {"input_tokens": 100, "cache_read_tokens": 400}
	}`

	session := Parse(strings.NewReader(jsonInput))

	if session.Model.DisplayName != "Opus" {
		t.Errorf("Expected model name Opus, got %s", session.Model.DisplayName)
	}
	if session.Cost.TotalCostUSD != 0.05 {
		t.Errorf("Expected cost 0.05, got %f", session.Cost.TotalCostUSD)
	}
	if session.ContextWindow.UsedPercentage != 42.5 {
		t.Errorf("Expected context 42.5, got %f", session.ContextWindow.UsedPercentage)
	}
	if session.CurrentUsage == nil {
		t.Error("Expected current_usage to be present")
	}
	if session.CurrentUsage.CacheReadTokens != 400 {
		t.Errorf("Expected cache_read_tokens 400, got %d", session.CurrentUsage.CacheReadTokens)
	}
}

func TestParse_EmptyJSON(t *testing.T) {
	session := Parse(strings.NewReader("{}"))

	if session.Model.DisplayName != "" {
		t.Errorf("Expected empty model name, got %s", session.Model.DisplayName)
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	session := Parse(strings.NewReader("not json"))

	// Should return empty session without panic
	if session == nil {
		t.Error("Expected non-nil session on invalid JSON")
	}
}

func TestParse_NullCurrentUsage(t *testing.T) {
	jsonInput := `{"model": {"display_name": "Opus"}, "current_usage": null}`
	session := Parse(strings.NewReader(jsonInput))

	if session.CurrentUsage != nil {
		t.Error("Expected current_usage to be nil")
	}
}
