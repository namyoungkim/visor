package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
	"github.com/namyoungkim/visor/internal/usage"
)

// Week limit thresholds (percentage utilization).
const (
	WeekLimitWarningPct  = 70.0
	WeekLimitCriticalPct = 90.0
)

// WeekLimitWidget displays the 7-day rate limit utilization.
// This is for Claude Pro users to see their weekly usage against the rate limit.
//
// Supported Extra options:
//   - show_label: "true"/"false" - whether to show "7d:" prefix (default: true)
//   - show_remaining: "true"/"false" - show remaining time (default: false)
//   - warn_threshold: "70" - % utilization for warning color (default: 70)
//   - critical_threshold: "90" - % utilization for critical color (default: 90)
type WeekLimitWidget struct {
	limits *usage.Limits
}

func (w *WeekLimitWidget) Name() string {
	return "week_limit"
}

// SetLimits sets the usage limits for this widget.
func (w *WeekLimitWidget) SetLimits(limits *usage.Limits) {
	w.limits = limits
}

func (w *WeekLimitWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	if w.limits == nil {
		return render.Colorize("—", "gray")
	}

	pct := w.limits.SevenDay.Utilization

	var value string
	if GetExtraBool(cfg, "show_remaining", false) {
		remaining := w.limits.SevenDayRemaining()
		if remaining > 0 {
			days := int(remaining.Hours()) / 24
			hours := int(remaining.Hours()) % 24
			if days > 0 {
				value = fmt.Sprintf("%.0f%% (%dd%dh)", pct, days, hours)
			} else {
				value = fmt.Sprintf("%.0f%% (%dh)", pct, hours)
			}
		} else {
			value = fmt.Sprintf("%.0f%%", pct)
		}
	} else {
		value = fmt.Sprintf("%.0f%%", pct)
	}

	var text string
	if cfg.Format != "" {
		text = FormatOutput(cfg, "", value)
	} else if GetExtraBool(cfg, "show_label", true) {
		text = "7d: " + value
	} else {
		text = value
	}

	warnThreshold := GetExtraFloat(cfg, "warn_threshold", WeekLimitWarningPct)
	criticalThreshold := GetExtraFloat(cfg, "critical_threshold", WeekLimitCriticalPct)
	color := ColorByThreshold(pct, warnThreshold, criticalThreshold)
	return render.Colorize(text, color)
}

func (w *WeekLimitWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return w.limits != nil && w.limits.SevenDay.Utilization > 0
}
