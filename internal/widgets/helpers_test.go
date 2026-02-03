package widgets

import (
	"testing"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
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

func TestGetExtraInt(t *testing.T) {
	cfg := &config.WidgetConfig{
		Extra: map[string]string{
			"width":   "15",
			"count":   "0",
			"invalid": "abc",
		},
	}

	tests := []struct {
		key      string
		defVal   int
		expected int
	}{
		{"width", 10, 15},
		{"count", 5, 0},
		{"invalid", 10, 10}, // Invalid number returns default
		{"missing", 10, 10},
	}

	for _, tt := range tests {
		result := GetExtraInt(cfg, tt.key, tt.defVal)
		if result != tt.expected {
			t.Errorf("GetExtraInt(%q, %d) = %d, expected %d",
				tt.key, tt.defVal, result, tt.expected)
		}
	}
}

func TestGetExtraFloat(t *testing.T) {
	cfg := &config.WidgetConfig{
		Extra: map[string]string{
			"threshold": "0.5",
			"rate":      "10.25",
			"zero":      "0",
			"negative":  "-5.5",
			"invalid":   "abc",
		},
	}

	tests := []struct {
		key      string
		defVal   float64
		expected float64
	}{
		{"threshold", 1.0, 0.5},
		{"rate", 0.0, 10.25},
		{"zero", 1.0, 0.0},
		{"negative", 0.0, -5.5},
		{"invalid", 1.0, 1.0},   // Invalid number returns default
		{"missing", 99.9, 99.9}, // Missing key returns default
	}

	for _, tt := range tests {
		result := GetExtraFloat(cfg, tt.key, tt.defVal)
		if result != tt.expected {
			t.Errorf("GetExtraFloat(%q, %.1f) = %.2f, expected %.2f",
				tt.key, tt.defVal, result, tt.expected)
		}
	}

	// Test nil Extra map
	cfgNil := &config.WidgetConfig{}
	if v := GetExtraFloat(cfgNil, "key", 42.0); v != 42.0 {
		t.Errorf("Expected 42.0 for nil Extra, got %.1f", v)
	}
}

func TestProgressBar(t *testing.T) {
	tests := []struct {
		pct      float64
		width    int
		expected string
	}{
		{0, 10, "░░░░░░░░░░"},
		{50, 10, "█████░░░░░"},
		{100, 10, "██████████"},
		{42, 10, "████░░░░░░"},
		{25, 4, "█░░░"},
		{75, 4, "███░"},
		{33, 6, "█░░░░░"},  // 33% of 6 = 1.98 → 1
		// Edge cases
		{-10, 10, "░░░░░░░░░░"},  // Negative clamped to 0
		{150, 10, "██████████"},  // Over 100 clamped to 100
		{50, 0, "█████░░░░░"},   // Zero width uses default (10)
		{50, -5, "█████░░░░░"},  // Negative width uses default (10)
	}

	for _, tt := range tests {
		result := ProgressBar(tt.pct, tt.width)
		if result != tt.expected {
			t.Errorf("ProgressBar(%.0f, %d) = %q, expected %q",
				tt.pct, tt.width, result, tt.expected)
		}
	}
}

func TestRenderAll(t *testing.T) {
	session := &input.Session{
		Model: input.Model{
			DisplayName: "Opus",
		},
		Cost: input.Cost{
			TotalCostUSD: 0.25,
		},
	}

	widgets := []config.WidgetConfig{
		{Name: "model"},
		{Name: "cost"},
	}

	result := RenderAll(session, widgets)

	if len(result) != 2 {
		t.Errorf("Expected 2 rendered widgets, got %d", len(result))
	}
}

func TestRenderAll_UnknownWidget(t *testing.T) {
	session := &input.Session{
		Model: input.Model{DisplayName: "Opus"},
	}

	widgets := []config.WidgetConfig{
		{Name: "unknown_widget"},
		{Name: "model"},
	}

	result := RenderAll(session, widgets)

	// Unknown widget should be skipped, model should render
	if len(result) != 1 {
		t.Errorf("Expected 1 rendered widget (unknown skipped), got %d", len(result))
	}
}

func TestRenderAll_EmptyWidgets(t *testing.T) {
	session := &input.Session{}

	result := RenderAll(session, []config.WidgetConfig{})

	if len(result) != 0 {
		t.Errorf("Expected 0 rendered widgets, got %d", len(result))
	}
}

func TestRenderAll_ShouldRenderFalse(t *testing.T) {
	session := &input.Session{
		// No code changes, so code_changes widget should not render
		Workspace: input.Workspace{
			LinesAdded:   0,
			LinesRemoved: 0,
		},
	}

	widgets := []config.WidgetConfig{
		{Name: "code_changes"},
	}

	result := RenderAll(session, widgets)

	// code_changes widget should not render when there are no changes
	if len(result) != 0 {
		t.Errorf("Expected 0 rendered widgets (code_changes should not render), got %d", len(result))
	}
}

func TestShouldRender_AlwaysTrue(t *testing.T) {
	// These widgets always render regardless of data
	alwaysTrueWidgets := []string{"context", "cost", "cache_hit", "api_latency"}
	session := &input.Session{}
	cfg := &config.WidgetConfig{}

	for _, name := range alwaysTrueWidgets {
		w, ok := Get(name)
		if !ok {
			t.Errorf("Widget %q not found", name)
			continue
		}
		if !w.ShouldRender(session, cfg) {
			t.Errorf("Widget %q.ShouldRender should return true", name)
		}
	}
}

func TestShouldRender_Model(t *testing.T) {
	w, _ := Get("model")
	cfg := &config.WidgetConfig{}

	// Model widget requires model data
	emptySession := &input.Session{}
	if w.ShouldRender(emptySession, cfg) {
		t.Error("Model widget should not render with empty session")
	}

	sessionWithModel := &input.Session{
		Model: input.Model{DisplayName: "Opus"},
	}
	if !w.ShouldRender(sessionWithModel, cfg) {
		t.Error("Model widget should render with model data")
	}
}

func TestShouldRender_Git(t *testing.T) {
	w, _ := Get("git")
	cfg := &config.WidgetConfig{}

	// Git widget always returns true (checks repo status internally)
	session := &input.Session{}
	// Note: actual git status depends on current directory
	_ = w.ShouldRender(session, cfg)
}
