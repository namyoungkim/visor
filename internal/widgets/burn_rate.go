package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// Burn rate thresholds (cents per minute).
const (
	BurnRateWarningCents  = 10.0 // 10¢/min
	BurnRateDangerCents   = 25.0 // 25¢/min
)

// BurnRateWidget displays the cost burn rate ($/min or ¢/min).
//
// Calculation: total_cost_usd / (total_duration_ms / 60000)
//
// Supported Extra options:
//   - show_label: "true"/"false" - whether to show "Burn:" prefix (default: false)
//   - warn_threshold: "10" - cents/min for warning color (default: 10)
//   - critical_threshold: "25" - cents/min for critical/red color (default: 25)
type BurnRateWidget struct{}

func (w *BurnRateWidget) Name() string {
	return "burn_rate"
}

func (w *BurnRateWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	cost := session.Cost.TotalCostUSD
	durationMs := session.Cost.TotalDurationMs

	// Cannot calculate burn rate without duration
	if durationMs <= 0 {
		return render.Colorize("—", "dim")
	}

	// Calculate burn rate: $ per minute
	durationMinutes := float64(durationMs) / 60000.0
	burnRatePerMin := cost / durationMinutes
	burnRateCents := burnRatePerMin * 100

	// Format output
	var value string
	if burnRatePerMin >= 1.0 {
		// $1.00+/min: show as dollars
		value = fmt.Sprintf("$%.1f/min", burnRatePerMin)
	} else if burnRateCents >= 0.1 {
		// 0.1¢+ per min: show as cents
		value = fmt.Sprintf("%.1f¢/min", burnRateCents)
	} else {
		value = "0.0¢/min"
	}

	var text string
	if cfg.Format != "" {
		text = FormatOutput(cfg, "", value)
	} else if GetExtraBool(cfg, "show_label", false) {
		text = "Burn: " + value
	} else {
		text = value
	}

	warnThreshold := GetExtraFloat(cfg, "warn_threshold", BurnRateWarningCents)
	criticalThreshold := GetExtraFloat(cfg, "critical_threshold", BurnRateDangerCents)
	color := ColorByThreshold(burnRateCents, warnThreshold, criticalThreshold)
	return render.Colorize(text, color)
}

func (w *BurnRateWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	// Don't render if no duration data available
	return session.Cost.TotalDurationMs > 0
}
