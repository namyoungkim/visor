package widgets

import (
	"strings"
	"testing"
	"time"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/history"
	"github.com/namyoungkim/visor/internal/input"
)

func TestBlockTimerWidget_Name(t *testing.T) {
	w := &BlockTimerWidget{}
	if w.Name() != "block_timer" {
		t.Errorf("expected 'block_timer', got %s", w.Name())
	}
}

func TestBlockTimerWidget_NoHistory(t *testing.T) {
	w := &BlockTimerWidget{}
	session := &input.Session{}
	cfg := &config.WidgetConfig{Name: "block_timer"}

	result := w.Render(session, cfg)
	if result != "" {
		t.Errorf("expected empty string for no history, got %q", result)
	}
}

func TestBlockTimerWidget_Render(t *testing.T) {
	tests := []struct {
		name           string
		hoursRemaining float64
		wantContains   string
		wantColor      string
	}{
		{
			name:           "4 hours remaining (green)",
			hoursRemaining: 4,
			wantContains:   "4h",
			wantColor:      "32m", // green ANSI
		},
		{
			name:           "1 hour remaining (yellow, 80% elapsed)",
			hoursRemaining: 1,
			wantContains:   "1h",
			wantColor:      "33m", // yellow ANSI
		},
		{
			name:           "10 minutes remaining (red, >95% elapsed)",
			hoursRemaining: 0.167, // 10 min
			wantContains:   "m",
			wantColor:      "31m", // red ANSI
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &BlockTimerWidget{}
			h := &history.History{SessionID: "test"}

			// Calculate block start time based on desired remaining time
			elapsedMs := int64(history.BlockDurationMs) - int64(tt.hoursRemaining*60*60*1000)
			h.BlockStartTime = time.Now().UnixMilli() - elapsedMs
			w.SetHistory(h)

			session := &input.Session{}
			cfg := &config.WidgetConfig{Name: "block_timer"}

			result := w.Render(session, cfg)

			if !strings.Contains(result, tt.wantContains) {
				t.Errorf("expected result to contain %q, got %q", tt.wantContains, result)
			}

			if !strings.Contains(result, tt.wantColor) {
				t.Errorf("expected color code %q, got %q", tt.wantColor, result)
			}
		})
	}
}

func TestBlockTimerWidget_ShouldRender(t *testing.T) {
	tests := []struct {
		name     string
		history  *history.History
		expected bool
	}{
		{
			name:     "nil history",
			history:  nil,
			expected: false,
		},
		{
			name:     "no block start time",
			history:  &history.History{SessionID: "test", BlockStartTime: 0},
			expected: false,
		},
		{
			name: "with block start time",
			history: &history.History{
				SessionID:      "test",
				BlockStartTime: time.Now().UnixMilli(),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &BlockTimerWidget{}
			if tt.history != nil {
				w.SetHistory(tt.history)
			}

			session := &input.Session{}
			cfg := &config.WidgetConfig{Name: "block_timer"}

			if got := w.ShouldRender(session, cfg); got != tt.expected {
				t.Errorf("ShouldRender() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBlockTimerWidget_WithShowLabel(t *testing.T) {
	w := &BlockTimerWidget{}
	h := &history.History{SessionID: "test"}
	h.BlockStartTime = time.Now().UnixMilli() // Just started
	w.SetHistory(h)

	session := &input.Session{}

	// With label (default)
	cfg := &config.WidgetConfig{Name: "block_timer"}
	result := w.Render(session, cfg)
	if !strings.Contains(result, "Block:") {
		t.Errorf("expected 'Block:' label, got %q", result)
	}

	// Without label
	cfg = &config.WidgetConfig{
		Name:  "block_timer",
		Extra: map[string]string{"show_label": "false"},
	}
	result = w.Render(session, cfg)
	if strings.Contains(result, "Block:") {
		t.Errorf("expected no 'Block:' label, got %q", result)
	}
}

func TestBlockTimerWidget_CustomThresholds(t *testing.T) {
	w := &BlockTimerWidget{}
	h := &history.History{SessionID: "test"}

	// 30% elapsed = 3.5 hours remaining
	elapsedMs := int64(float64(history.BlockDurationMs) * 0.30)
	h.BlockStartTime = time.Now().UnixMilli() - elapsedMs
	w.SetHistory(h)

	session := &input.Session{}

	// Default thresholds: 30% elapsed should be green
	cfg := &config.WidgetConfig{Name: "block_timer"}
	result := w.Render(session, cfg)
	if !strings.Contains(result, "32m") { // green
		t.Errorf("expected green at 30%% elapsed with default thresholds, got %q", result)
	}

	// Custom thresholds: warn at 20%, critical at 40%
	// 30% elapsed should now be yellow
	cfg = &config.WidgetConfig{
		Name: "block_timer",
		Extra: map[string]string{
			"warn_threshold":     "20",
			"critical_threshold": "40",
		},
	}
	result = w.Render(session, cfg)
	if !strings.Contains(result, "33m") { // yellow
		t.Errorf("expected yellow at 30%% elapsed with custom thresholds, got %q", result)
	}
}
