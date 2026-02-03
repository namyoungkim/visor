package widgets

import (
	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/cost"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// Weekly cost thresholds (USD).
const (
	WeeklyCostWarningUSD  = 25.0
	WeeklyCostCriticalUSD = 50.0
)

// WeeklyCostWidget displays this week's aggregated cost.
//
// Supported Extra options:
//   - show_label: "true"/"false" - whether to show "Week:" prefix (default: false)
//   - warn_threshold: "25.0" - USD for warning color (default: 25.0)
//   - critical_threshold: "50.0" - USD for critical/red color (default: 50.0)
type WeeklyCostWidget struct {
	costData *cost.CostData
}

func (w *WeeklyCostWidget) Name() string {
	return "weekly_cost"
}

// SetCostData sets the cost data for this widget.
func (w *WeeklyCostWidget) SetCostData(data *cost.CostData) {
	w.costData = data
}

func (w *WeeklyCostWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	if w.costData == nil {
		return render.Colorize("—", "gray")
	}

	value := formatCost(w.costData.Week)

	var text string
	if cfg.Format != "" {
		text = FormatOutput(cfg, "", value)
	} else if GetExtraBool(cfg, "show_label", false) {
		text = "Week: " + value
	} else {
		text = value
	}

	warnThreshold := GetExtraFloat(cfg, "warn_threshold", WeeklyCostWarningUSD)
	criticalThreshold := GetExtraFloat(cfg, "critical_threshold", WeeklyCostCriticalUSD)
	color := ColorByThreshold(w.costData.Week, warnThreshold, criticalThreshold)
	return render.Colorize(text, color)
}

func (w *WeeklyCostWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return w.costData != nil
}
