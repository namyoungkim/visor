package widgets

import (
	"strings"
	"testing"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

func TestCostWidget_ZeroCost(t *testing.T) {
	w := &CostWidget{}
	session := &input.Session{}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "$0.00") {
		t.Errorf("Expected '$0.00', got '%s'", result)
	}
}

func TestCostWidget_SmallCost(t *testing.T) {
	w := &CostWidget{}
	session := &input.Session{
		Cost: input.Cost{TotalCostUSD: 0.005},
	}

	result := w.Render(session, &config.WidgetConfig{})

	// Small costs should show 3 decimal places
	if !strings.Contains(result, "$0.005") {
		t.Errorf("Expected '$0.005', got '%s'", result)
	}
}

func TestCostWidget_NormalCost(t *testing.T) {
	w := &CostWidget{}
	session := &input.Session{
		Cost: input.Cost{TotalCostUSD: 0.15},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "$0.15") {
		t.Errorf("Expected '$0.15', got '%s'", result)
	}
}

func TestCostWidget_LargeCost(t *testing.T) {
	w := &CostWidget{}
	session := &input.Session{
		Cost: input.Cost{TotalCostUSD: 2.50},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "$2.50") {
		t.Errorf("Expected '$2.50', got '%s'", result)
	}
}

func TestCostWidget_ColorThresholds(t *testing.T) {
	w := &CostWidget{}

	tests := []struct {
		cost          float64
		expectedColor string
	}{
		{0.0, render.FgGreen},
		{0.49, render.FgGreen},
		{0.50, render.FgYellow},
		{0.99, render.FgYellow},
		{1.00, render.FgRed},
		{5.00, render.FgRed},
	}

	for _, tt := range tests {
		session := &input.Session{
			Cost: input.Cost{TotalCostUSD: tt.cost},
		}
		result := w.Render(session, &config.WidgetConfig{})

		if !strings.Contains(result, tt.expectedColor) {
			t.Errorf("At $%.2f, expected color %s, got '%s'", tt.cost, tt.expectedColor, result)
		}
	}
}
