package widgets

import (
	"strings"
	"testing"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
)

func TestCompactETAWidget_Name(t *testing.T) {
	w := &CompactETAWidget{}
	if w.Name() != "compact_eta" {
		t.Errorf("expected 'compact_eta', got %s", w.Name())
	}
}

func TestCompactETAWidget_Render(t *testing.T) {
	tests := []struct {
		name       string
		pct        float64
		durationMs int64
		wantValue  string
	}{
		{
			name:       "less than 1 minute ETA",
			pct:        42.0,
			durationMs: 60000, // 1 min, 42%/min burn rate, (80-42)/42 = ~0.9 min
			wantValue:  "<1m",
		},
		{
			name:       "7 minute estimate",
			pct:        10.0,
			durationMs: 60000, // 10%/min, (80-10)/10 = 7 min
			wantValue:  "~7m",
		},
		{
			name:       "already at compact threshold",
			pct:        80.0,
			durationMs: 60000,
			wantValue:  "compact soon",
		},
		{
			name:       "above compact threshold",
			pct:        85.0,
			durationMs: 60000,
			wantValue:  "compact soon",
		},
		{
			name:       "zero duration shows dash",
			pct:        50.0,
			durationMs: 0,
			wantValue:  "—",
		},
		{
			name:       "very slow burn - hours estimate",
			pct:        5.0,
			durationMs: 300000, // 5 min, 1%/min burn rate, (80-5)/1 = 75 min = 1h15m
			wantValue:  "~1h",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &CompactETAWidget{}
			session := &input.Session{
				ContextWindow: input.ContextWindow{
					UsedPercentage: tt.pct,
				},
				Cost: input.Cost{
					TotalDurationMs: tt.durationMs,
				},
			}
			cfg := &config.WidgetConfig{Name: "compact_eta"}

			result := w.Render(session, cfg)

			if !strings.Contains(result, tt.wantValue) {
				t.Errorf("expected output to contain %q, got %q", tt.wantValue, result)
			}
		})
	}
}

func TestCompactETAWidget_ShouldRender(t *testing.T) {
	w := &CompactETAWidget{}

	tests := []struct {
		name       string
		pct        float64
		durationMs int64
		threshold  string
		want       bool
	}{
		{"above default threshold", 50.0, 60000, "", true},
		{"at default threshold", 40.0, 60000, "", true},
		{"below default threshold", 30.0, 60000, "", false},
		{"custom threshold above", 60.0, 60000, "50", true},
		{"custom threshold below", 40.0, 60000, "50", false},
		{"no duration data", 50.0, 0, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &input.Session{
				ContextWindow: input.ContextWindow{
					UsedPercentage: tt.pct,
				},
				Cost: input.Cost{
					TotalDurationMs: tt.durationMs,
				},
			}
			cfg := &config.WidgetConfig{Name: "compact_eta"}
			if tt.threshold != "" {
				cfg.Extra = map[string]string{"show_when_above": tt.threshold}
			}

			got := w.ShouldRender(session, cfg)
			if got != tt.want {
				t.Errorf("ShouldRender() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompactETAWidget_ColorCoding(t *testing.T) {
	w := &CompactETAWidget{}

	tests := []struct {
		name       string
		pct        float64
		durationMs int64
		wantColor  string
	}{
		{
			name:       "green for long ETA",
			pct:        5.0,
			durationMs: 300000, // 5 min, 1%/min, (80-5)/1 = 75 min > 10 min warning
			wantColor:  "32", // green ANSI
		},
		{
			name:       "yellow for medium ETA",
			pct:        10.0,
			durationMs: 60000, // 10%/min, (80-10)/10 = 7 min (between 5-10)
			wantColor:  "33", // yellow ANSI
		},
		{
			name:       "red at compact threshold",
			pct:        80.0,
			durationMs: 60000,
			wantColor:  "31", // red
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &input.Session{
				ContextWindow: input.ContextWindow{
					UsedPercentage: tt.pct,
				},
				Cost: input.Cost{
					TotalDurationMs: tt.durationMs,
				},
			}
			cfg := &config.WidgetConfig{Name: "compact_eta"}

			result := w.Render(session, cfg)

			if !strings.Contains(result, tt.wantColor) {
				t.Errorf("expected ANSI color %q in output, got %q", tt.wantColor, result)
			}
		})
	}
}

func TestCompactETAWidget_WithFormat(t *testing.T) {
	w := &CompactETAWidget{}
	session := &input.Session{
		ContextWindow: input.ContextWindow{
			UsedPercentage: 40.0,
		},
		Cost: input.Cost{
			TotalDurationMs: 60000,
		},
	}
	cfg := &config.WidgetConfig{
		Name:   "compact_eta",
		Format: "Compact: {value}",
	}

	result := w.Render(session, cfg)
	if !strings.Contains(result, "Compact:") {
		t.Errorf("expected formatted output with 'Compact:', got %q", result)
	}
}
