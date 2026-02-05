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
	cu := session.GetCurrentUsage()
	if cu == nil {
		t.Error("Expected current_usage to be present")
	}
	if cu.GetCacheReadTokens() != 400 {
		t.Errorf("Expected cache_read_tokens 400, got %d", cu.GetCacheReadTokens())
	}
}

func TestParse_ContextWindowCurrentUsage(t *testing.T) {
	jsonInput := `{
		"model": {"display_name": "Opus"},
		"context_window": {
			"used_percentage": 42.5,
			"total_input_tokens": 15000,
			"total_output_tokens": 3000,
			"current_usage": {"input_tokens": 200, "cache_read_input_tokens": 800}
		}
	}`

	session := Parse(strings.NewReader(jsonInput))

	if session.ContextWindow.TotalOutputTokens != 3000 {
		t.Errorf("Expected total_output_tokens 3000, got %d", session.ContextWindow.TotalOutputTokens)
	}
	cu := session.GetCurrentUsage()
	if cu == nil {
		t.Fatal("Expected current_usage from context_window")
	}
	if cu.GetCacheReadTokens() != 800 {
		t.Errorf("Expected cache_read_input_tokens 800, got %d", cu.GetCacheReadTokens())
	}
	if session.GetTotalOutputTokens() != 3000 {
		t.Errorf("Expected GetTotalOutputTokens() 3000, got %d", session.GetTotalOutputTokens())
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

	if session.GetCurrentUsage() != nil {
		t.Error("Expected current_usage to be nil")
	}
}

func TestGetTotalOutputTokens_FallbackToCost(t *testing.T) {
	session := &Session{
		Cost: Cost{TotalOutputTokens: 5000},
	}
	if session.GetTotalOutputTokens() != 5000 {
		t.Errorf("Expected fallback to Cost.TotalOutputTokens, got %d", session.GetTotalOutputTokens())
	}
}

func TestGetCacheReadTokens_PreferNewField(t *testing.T) {
	cu := &CurrentUsage{
		CacheReadTokens:      100,
		CacheReadInputTokens: 200,
	}
	if cu.GetCacheReadTokens() != 200 {
		t.Errorf("Expected cache_read_input_tokens preferred, got %d", cu.GetCacheReadTokens())
	}
}

func TestGetCacheReadTokens_FallbackToOldField(t *testing.T) {
	cu := &CurrentUsage{
		CacheReadTokens: 100,
	}
	if cu.GetCacheReadTokens() != 100 {
		t.Errorf("Expected fallback to cache_read_tokens, got %d", cu.GetCacheReadTokens())
	}
}
