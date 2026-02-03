package widgets

import (
	"strings"
	"testing"
	"time"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/usage"
)

func TestBlockLimitWidget_Name(t *testing.T) {
	w := &BlockLimitWidget{}
	if w.Name() != "block_limit" {
		t.Errorf("Name() = %q, want %q", w.Name(), "block_limit")
	}
}

func TestBlockLimitWidget_Render(t *testing.T) {
	tests := []struct {
		name   string
		limits *usage.Limits
		cfg    *config.WidgetConfig
		want   string
	}{
		{
			name: "renders utilization with remaining time",
			limits: &usage.Limits{
				FiveHour: usage.FiveHourLimit{
					Utilization: 50.0,
					ResetsAt:    time.Now().Add(2*time.Hour + 30*time.Minute),
				},
			},
			cfg:  &config.WidgetConfig{},
			want: "5h: 50%",
		},
		{
			name: "without remaining time",
			limits: &usage.Limits{
				FiveHour: usage.FiveHourLimit{
					Utilization: 75.0,
					ResetsAt:    time.Now().Add(1 * time.Hour),
				},
			},
			cfg: &config.WidgetConfig{
				Extra: map[string]string{"show_remaining": "false"},
			},
			want: "5h: 75%",
		},
		{
			name:   "no data returns dash",
			limits: nil,
			cfg:    &config.WidgetConfig{},
			want:   "—",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &BlockLimitWidget{}
			w.SetLimits(tt.limits)

			session := &input.Session{}
			got := w.Render(session, tt.cfg)

			// Strip ANSI codes for comparison
			got = stripANSI(got)
			if !strings.Contains(got, tt.want) {
				t.Errorf("Render() = %q, want to contain %q", got, tt.want)
			}
		})
	}
}

func TestBlockLimitWidget_ShouldRender(t *testing.T) {
	w := &BlockLimitWidget{}
	session := &input.Session{}
	cfg := &config.WidgetConfig{}

	// Without limits
	if w.ShouldRender(session, cfg) {
		t.Error("ShouldRender() = true without limits, want false")
	}

	// With zero utilization
	w.SetLimits(&usage.Limits{FiveHour: usage.FiveHourLimit{Utilization: 0}})
	if w.ShouldRender(session, cfg) {
		t.Error("ShouldRender() = true with zero utilization, want false")
	}

	// With positive utilization
	w.SetLimits(&usage.Limits{FiveHour: usage.FiveHourLimit{Utilization: 50.0}})
	if !w.ShouldRender(session, cfg) {
		t.Error("ShouldRender() = false with positive utilization, want true")
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		want     string
	}{
		{2*time.Hour + 30*time.Minute, "2h30m"},
		{45 * time.Minute, "45m"},
		{5 * time.Hour, "5h0m"},
		{0, "0m"},
	}

	for _, tt := range tests {
		got := formatDuration(tt.duration)
		if got != tt.want {
			t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, got, tt.want)
		}
	}
}
