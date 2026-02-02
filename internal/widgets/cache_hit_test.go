package widgets

import (
	"strings"
	"testing"

	"github.com/leo/visor/internal/config"
	"github.com/leo/visor/internal/input"
)

func TestCacheHitWidget_NilCurrentUsage(t *testing.T) {
	w := &CacheHitWidget{}
	session := &input.Session{}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "—") {
		t.Errorf("Expected dash for nil current_usage, got '%s'", result)
	}
}

func TestCacheHitWidget_ZeroTokens(t *testing.T) {
	w := &CacheHitWidget{}
	session := &input.Session{
		CurrentUsage: &input.CurrentUsage{
			InputTokens:     0,
			CacheReadTokens: 0,
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "—") {
		t.Errorf("Expected dash for zero tokens, got '%s'", result)
	}
}

func TestCacheHitWidget_HighCacheRate(t *testing.T) {
	w := &CacheHitWidget{}
	session := &input.Session{
		CurrentUsage: &input.CurrentUsage{
			InputTokens:     100,
			CacheReadTokens: 400, // 80% cache hit
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "80%") {
		t.Errorf("Expected 80%%, got '%s'", result)
	}
	// Should be green (high cache rate)
	if !strings.Contains(result, "\033[32m") {
		t.Errorf("Expected green color for high cache rate")
	}
}

func TestCacheHitWidget_LowCacheRate(t *testing.T) {
	w := &CacheHitWidget{}
	session := &input.Session{
		CurrentUsage: &input.CurrentUsage{
			InputTokens:     400,
			CacheReadTokens: 100, // 20% cache hit
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "20%") {
		t.Errorf("Expected 20%%, got '%s'", result)
	}
	// Should be red (low cache rate)
	if !strings.Contains(result, "\033[31m") {
		t.Errorf("Expected red color for low cache rate")
	}
}
