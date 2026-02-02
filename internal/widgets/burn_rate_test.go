package widgets

import (
	"strings"
	"testing"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
)

func TestBurnRateWidget_Name(t *testing.T) {
	w := &BurnRateWidget{}
	if w.Name() != "burn_rate" {
		t.Errorf("expected 'burn_rate', got %s", w.Name())
	}
}

func TestBurnRateWidget_Render(t *testing.T) {
	tests := []struct {
		name       string
		cost       float64
		durationMs int64
		wantValue  string
		wantColor  string
	}{
		{
			name:       "low burn rate",
			cost:       0.10,
			durationMs: 60000, // 1 minute
			wantValue:  "10.0¢/min",
			wantColor:  "yellow", // 10¢ = warning threshold
		},
		{
			name:       "high burn rate",
			cost:       0.50,
			durationMs: 60000, // 1 minute = 50¢/min
			wantValue:  "50.0¢/min",
			wantColor:  "red",
		},
		{
			name:       "very low burn rate",
			cost:       0.01,
			durationMs: 120000, // 2 minutes = 0.5¢/min
			wantValue:  "0.5¢/min",
			wantColor:  "green",
		},
		{
			name:       "dollar per minute",
			cost:       2.00,
			durationMs: 60000, // 1 minute = $2/min
			wantValue:  "$2.0/min",
			wantColor:  "red",
		},
		{
			name:       "zero duration shows dash",
			cost:       0.50,
			durationMs: 0,
			wantValue:  "—",
			wantColor:  "dim",
		},
		{
			name:       "realistic scenario",
			cost:       0.48,
			durationMs: 45000, // 0.75 minutes = 64¢/min
			wantValue:  "64.0¢/min",
			wantColor:  "red",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &BurnRateWidget{}
			session := &input.Session{
				Cost: input.Cost{
					TotalCostUSD:    tt.cost,
					TotalDurationMs: tt.durationMs,
				},
			}
			cfg := &config.WidgetConfig{Name: "burn_rate"}

			result := w.Render(session, cfg)

			if !strings.Contains(result, tt.wantValue) {
				t.Errorf("expected output to contain %q, got %q", tt.wantValue, result)
			}
		})
	}
}

func TestBurnRateWidget_ShouldRender(t *testing.T) {
	w := &BurnRateWidget{}

	tests := []struct {
		name       string
		durationMs int64
		want       bool
	}{
		{"with duration", 60000, true},
		{"zero duration", 0, false},
		{"negative duration", -1000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &input.Session{
				Cost: input.Cost{TotalDurationMs: tt.durationMs},
			}
			cfg := &config.WidgetConfig{Name: "burn_rate"}

			got := w.ShouldRender(session, cfg)
			if got != tt.want {
				t.Errorf("ShouldRender() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBurnRateWidget_WithFormat(t *testing.T) {
	w := &BurnRateWidget{}
	session := &input.Session{
		Cost: input.Cost{
			TotalCostUSD:    0.10,
			TotalDurationMs: 60000,
		},
	}
	cfg := &config.WidgetConfig{
		Name:   "burn_rate",
		Format: "Rate: {value}",
	}

	result := w.Render(session, cfg)
	if !strings.Contains(result, "Rate: 10.0¢/min") {
		t.Errorf("expected formatted output with 'Rate:', got %q", result)
	}
}

func TestBurnRateWidget_WithShowLabel(t *testing.T) {
	w := &BurnRateWidget{}
	session := &input.Session{
		Cost: input.Cost{
			TotalCostUSD:    0.10,
			TotalDurationMs: 60000,
		},
	}
	cfg := &config.WidgetConfig{
		Name:  "burn_rate",
		Extra: map[string]string{"show_label": "true"},
	}

	result := w.Render(session, cfg)
	if !strings.Contains(result, "Burn:") {
		t.Errorf("expected output with 'Burn:' label, got %q", result)
	}
}
