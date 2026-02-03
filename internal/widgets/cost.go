package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// CostWidget displays the total API cost.
//
// Supported Extra options:
//   - show_label: "true"/"false" - whether to show "Cost:" prefix (default: false)
//   - warn_threshold: "0.5" - USD amount for warning color (default: 0.5)
//   - critical_threshold: "1.0" - USD amount for critical/red color (default: 1.0)
type CostWidget struct{}

func (w *CostWidget) Name() string {
	return "cost"
}

func (w *CostWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	cost := session.Cost.TotalCostUSD

	var value string
	switch {
	case cost >= 0.01:
		value = fmt.Sprintf("$%.2f", cost)
	case cost > 0:
		value = fmt.Sprintf("$%.3f", cost)
	default:
		value = "$0.00"
	}

	var text string
	if cfg.Format != "" {
		text = FormatOutput(cfg, "", value)
	} else if GetExtraBool(cfg, "show_label", false) {
		text = "Cost: " + value
	} else {
		text = value
	}

	warnThreshold := GetExtraFloat(cfg, "warn_threshold", CostWarningUSD)
	criticalThreshold := GetExtraFloat(cfg, "critical_threshold", CostDangerUSD)
	color := ColorByThreshold(cost, warnThreshold, criticalThreshold)
	return render.Colorize(text, color)
}

func (w *CostWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return true
}
