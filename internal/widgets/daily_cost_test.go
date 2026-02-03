package widgets

import (
	"strings"
	"testing"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/cost"
	"github.com/namyoungkim/visor/internal/input"
)

func TestDailyCostWidget_Name(t *testing.T) {
	w := &DailyCostWidget{}
	if w.Name() != "daily_cost" {
		t.Errorf("Name() = %q, want %q", w.Name(), "daily_cost")
	}
}

func TestDailyCostWidget_Render(t *testing.T) {
	tests := []struct {
		name     string
		costData *cost.CostData
		cfg      *config.WidgetConfig
		want     string
	}{
		{
			name: "renders today cost",
			costData: &cost.CostData{
				Today: 2.50,
			},
			cfg:  &config.WidgetConfig{},
			want: "$2.5",
		},
		{
			name: "with label",
			costData: &cost.CostData{
				Today: 5.00,
			},
			cfg: &config.WidgetConfig{
				Extra: map[string]string{"show_label": "true"},
			},
			want: "Today: $5.0",
		},
		{
			name:     "no data returns dash",
			costData: nil,
			cfg:      &config.WidgetConfig{},
			want:     "—",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &DailyCostWidget{}
			w.SetCostData(tt.costData)

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

func TestDailyCostWidget_ShouldRender(t *testing.T) {
	w := &DailyCostWidget{}
	session := &input.Session{}
	cfg := &config.WidgetConfig{}

	// Without cost data
	if w.ShouldRender(session, cfg) {
		t.Error("ShouldRender() = true without cost data, want false")
	}

	// With cost data
	w.SetCostData(&cost.CostData{Today: 1.0})
	if !w.ShouldRender(session, cfg) {
		t.Error("ShouldRender() = false with cost data, want true")
	}
}
