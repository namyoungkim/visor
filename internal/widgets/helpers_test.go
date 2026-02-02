package widgets

import (
	"testing"

	"github.com/namyoungkim/visor/internal/config"
)

func TestColorByThreshold(t *testing.T) {
	tests := []struct {
		value    float64
		warning  float64
		danger   float64
		expected string
	}{
		{0.0, 50.0, 80.0, "green"},
		{49.9, 50.0, 80.0, "green"},
		{50.0, 50.0, 80.0, "yellow"},
		{79.9, 50.0, 80.0, "yellow"},
		{80.0, 50.0, 80.0, "red"},
		{100.0, 50.0, 80.0, "red"},
	}

	for _, tt := range tests {
		result := ColorByThreshold(tt.value, tt.warning, tt.danger)
		if result != tt.expected {
			t.Errorf("ColorByThreshold(%.1f, %.1f, %.1f) = %s, expected %s",
				tt.value, tt.warning, tt.danger, result, tt.expected)
		}
	}
}

func TestColorByThresholdInverse(t *testing.T) {
	tests := []struct {
		value    float64
		good     float64
		warning  float64
		expected string
	}{
		{100.0, 80.0, 50.0, "green"},
		{80.0, 80.0, 50.0, "green"},
		{79.9, 80.0, 50.0, "yellow"},
		{50.0, 80.0, 50.0, "yellow"},
		{49.9, 80.0, 50.0, "red"},
		{0.0, 80.0, 50.0, "red"},
	}

	for _, tt := range tests {
		result := ColorByThresholdInverse(tt.value, tt.good, tt.warning)
		if result != tt.expected {
			t.Errorf("ColorByThresholdInverse(%.1f, %.1f, %.1f) = %s, expected %s",
				tt.value, tt.good, tt.warning, result, tt.expected)
		}
	}
}

func TestRegistry(t *testing.T) {
	// Test that all widgets are registered
	expectedWidgets := []string{
		"model", "context", "git", "cost",
		"cache_hit", "api_latency", "code_changes",
	}

	for _, name := range expectedWidgets {
		w, ok := Get(name)
		if !ok {
			t.Errorf("Widget '%s' not found in registry", name)
			continue
		}
		if w.Name() != name {
			t.Errorf("Widget name mismatch: expected '%s', got '%s'", name, w.Name())
		}
	}
}

func TestGet_NonexistentWidget(t *testing.T) {
	_, ok := Get("nonexistent")
	if ok {
		t.Error("Expected false for nonexistent widget")
	}
}

func TestFormatOutput(t *testing.T) {
	tests := []struct {
		format   string
		defFmt   string
		value    string
		expected string
	}{
		{"", "", "42%", "42%"},
		{"Context: {value}", "", "42%", "Context: 42%"},
		{"{value} used", "", "42%", "42% used"},
		{"Value is {value}!", "", "100", "Value is 100!"},
	}

	for _, tt := range tests {
		cfg := &config.WidgetConfig{Format: tt.format}
		result := FormatOutput(cfg, tt.defFmt, tt.value)
		if result != tt.expected {
			t.Errorf("FormatOutput(%q, %q, %q) = %q, expected %q",
				tt.format, tt.defFmt, tt.value, result, tt.expected)
		}
	}
}

func TestGetExtra(t *testing.T) {
	cfg := &config.WidgetConfig{
		Extra: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	if v := GetExtra(cfg, "key1", "default"); v != "value1" {
		t.Errorf("Expected 'value1', got '%s'", v)
	}

	if v := GetExtra(cfg, "missing", "default"); v != "default" {
		t.Errorf("Expected 'default', got '%s'", v)
	}

	// Nil Extra map
	cfgNil := &config.WidgetConfig{}
	if v := GetExtra(cfgNil, "key", "default"); v != "default" {
		t.Errorf("Expected 'default' for nil Extra, got '%s'", v)
	}
}

func TestGetExtraBool(t *testing.T) {
	cfg := &config.WidgetConfig{
		Extra: map[string]string{
			"enabled":  "true",
			"disabled": "false",
			"yes":      "yes",
			"one":      "1",
		},
	}

	tests := []struct {
		key      string
		defVal   bool
		expected bool
	}{
		{"enabled", false, true},
		{"disabled", true, false},
		{"yes", false, true},
		{"one", false, true},
		{"missing", true, true},
		{"missing", false, false},
	}

	for _, tt := range tests {
		result := GetExtraBool(cfg, tt.key, tt.defVal)
		if result != tt.expected {
			t.Errorf("GetExtraBool(%q, %v) = %v, expected %v",
				tt.key, tt.defVal, result, tt.expected)
		}
	}
}
