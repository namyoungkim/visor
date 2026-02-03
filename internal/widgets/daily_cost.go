package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/cost"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// Daily cost thresholds (USD).
const (
	DailyCostWarningUSD  = 5.0
	DailyCostCriticalUSD = 10.0
)

// DailyCostWidget displays today's aggregated cost.
//
// Supported Extra options:
//   - show_label: "true"/"false" - whether to show "Today:" prefix (default: false)
//   - warn_threshold: "5.0" - USD for warning color (default: 5.0)
//   - critical_threshold: "10.0" - USD for critical/red color (default: 10.0)
type DailyCostWidget struct {
	costData *cost.CostData
}

func (w *DailyCostWidget) Name() string {
	return "daily_cost"
}

// SetCostData sets the cost data for this widget.
func (w *DailyCostWidget) SetCostData(data *cost.CostData) {
	w.costData = data
}

func (w *DailyCostWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	if w.costData == nil {
		return render.Colorize("—", "gray")
	}

	value := formatCost(w.costData.Today)

	var text string
	if cfg.Format != "" {
		text = FormatOutput(cfg, "", value)
	} else if GetExtraBool(cfg, "show_label", false) {
		text = "Today: " + value
	} else {
		text = value
	}

	warnThreshold := GetExtraFloat(cfg, "warn_threshold", DailyCostWarningUSD)
	criticalThreshold := GetExtraFloat(cfg, "critical_threshold", DailyCostCriticalUSD)
	color := ColorByThreshold(w.costData.Today, warnThreshold, criticalThreshold)
	return render.Colorize(text, color)
}

func (w *DailyCostWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return w.costData != nil
}

// formatCost formats a cost value for display.
func formatCost(cost float64) string {
	switch {
	case cost >= 10.0:
		return fmt.Sprintf("$%.0f", cost)
	case cost >= 1.0:
		return fmt.Sprintf("$%.1f", cost)
	case cost >= 0.01:
		return fmt.Sprintf("$%.2f", cost)
	case cost > 0:
		return fmt.Sprintf("$%.3f", cost)
	default:
		return "$0.00"
	}
}
