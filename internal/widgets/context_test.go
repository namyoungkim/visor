package widgets

import (
	"strings"
	"testing"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

func TestContextWidget_LowUsage(t *testing.T) {
	w := &ContextWidget{}
	session := &input.Session{
		ContextWindow: input.ContextWindow{
			UsedPercentage: 30.0,
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "30%") {
		t.Errorf("Expected '30%%', got '%s'", result)
	}
	// Should be green
	if !strings.Contains(result, render.FgGreen) {
		t.Error("Expected green color for low usage")
	}
}

func TestContextWidget_MediumUsage(t *testing.T) {
	w := &ContextWidget{}
	session := &input.Session{
		ContextWindow: input.ContextWindow{
			UsedPercentage: 65.0,
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "65%") {
		t.Errorf("Expected '65%%', got '%s'", result)
	}
	// Should be yellow
	if !strings.Contains(result, render.FgYellow) {
		t.Error("Expected yellow color for medium usage")
	}
}

func TestContextWidget_HighUsage(t *testing.T) {
	w := &ContextWidget{}
	session := &input.Session{
		ContextWindow: input.ContextWindow{
			UsedPercentage: 85.0,
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "85%") {
		t.Errorf("Expected '85%%', got '%s'", result)
	}
	// Should be red
	if !strings.Contains(result, render.FgRed) {
		t.Error("Expected red color for high usage")
	}
}

func TestContextWidget_ZeroUsage(t *testing.T) {
	w := &ContextWidget{}
	session := &input.Session{}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "0%") {
		t.Errorf("Expected '0%%', got '%s'", result)
	}
}

func TestContextWidget_Thresholds(t *testing.T) {
	w := &ContextWidget{}

	tests := []struct {
		pct           float64
		expectedColor string
	}{
		{59.9, render.FgGreen},
		{60.0, render.FgYellow},
		{79.9, render.FgYellow},
		{80.0, render.FgRed},
		{100.0, render.FgRed},
	}

	for _, tt := range tests {
		session := &input.Session{
			ContextWindow: input.ContextWindow{UsedPercentage: tt.pct},
		}
		result := w.Render(session, &config.WidgetConfig{})

		if !strings.Contains(result, tt.expectedColor) {
			t.Errorf("At %.1f%%, expected color %s, got '%s'", tt.pct, tt.expectedColor, result)
		}
	}
}

func TestContextWidget_ProgressBar(t *testing.T) {
	w := &ContextWidget{}
	session := &input.Session{
		ContextWindow: input.ContextWindow{UsedPercentage: 42.0},
	}

	// Default: show_bar = true
	result := w.Render(session, &config.WidgetConfig{})
	if !strings.Contains(result, "████░░░░░░") {
		t.Errorf("Expected progress bar in output, got '%s'", result)
	}
	if !strings.Contains(result, "42%") {
		t.Errorf("Expected '42%%' in output, got '%s'", result)
	}
}

func TestContextWidget_ProgressBar_Disabled(t *testing.T) {
	w := &ContextWidget{}
	session := &input.Session{
		ContextWindow: input.ContextWindow{UsedPercentage: 42.0},
	}

	cfg := &config.WidgetConfig{
		Extra: map[string]string{"show_bar": "false"},
	}
	result := w.Render(session, cfg)

	// Should not contain progress bar characters
	if strings.Contains(result, "█") || strings.Contains(result, "░") {
		t.Errorf("Expected no progress bar when show_bar=false, got '%s'", result)
	}
	if !strings.Contains(result, "42%") {
		t.Errorf("Expected '42%%' in output, got '%s'", result)
	}
}

func TestContextWidget_ProgressBar_CustomWidth(t *testing.T) {
	w := &ContextWidget{}
	session := &input.Session{
		ContextWindow: input.ContextWindow{UsedPercentage: 50.0},
	}

	cfg := &config.WidgetConfig{
		Extra: map[string]string{"bar_width": "5"},
	}
	result := w.Render(session, cfg)

	// 50% of 5 chars = 2.5, rounds to 2
	if !strings.Contains(result, "██░░░") {
		t.Errorf("Expected 5-char progress bar, got '%s'", result)
	}
}
