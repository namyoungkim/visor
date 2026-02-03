package widgets

import (
	"fmt"
	"time"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
	"github.com/namyoungkim/visor/internal/usage"
)

// Block limit thresholds (percentage utilization).
const (
	BlockLimitWarningPct  = 70.0
	BlockLimitCriticalPct = 90.0
)

// BlockLimitWidget displays the 5-hour rate limit utilization.
// This is for Claude Pro users to see their usage against the rate limit.
//
// Supported Extra options:
//   - show_label: "true"/"false" - whether to show "5h:" prefix (default: true)
//   - show_remaining: "true"/"false" - show remaining time (default: true)
//   - warn_threshold: "70" - % utilization for warning color (default: 70)
//   - critical_threshold: "90" - % utilization for critical color (default: 90)
type BlockLimitWidget struct {
	limits *usage.Limits
}

func (w *BlockLimitWidget) Name() string {
	return "block_limit"
}

// SetLimits sets the usage limits for this widget.
func (w *BlockLimitWidget) SetLimits(limits *usage.Limits) {
	w.limits = limits
}

func (w *BlockLimitWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	if w.limits == nil {
		return render.Colorize("—", "gray")
	}

	pct := w.limits.FiveHour.Utilization

	var value string
	if GetExtraBool(cfg, "show_remaining", true) {
		remaining := w.limits.FiveHourRemaining()
		if remaining > 0 {
			value = fmt.Sprintf("%.0f%% (%s)", pct, formatDuration(remaining))
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
		text = "5h: " + value
	} else {
		text = value
	}

	warnThreshold := GetExtraFloat(cfg, "warn_threshold", BlockLimitWarningPct)
	criticalThreshold := GetExtraFloat(cfg, "critical_threshold", BlockLimitCriticalPct)
	color := ColorByThreshold(pct, warnThreshold, criticalThreshold)
	return render.Colorize(text, color)
}

func (w *BlockLimitWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return w.limits != nil && w.limits.FiveHour.Utilization > 0
}

// formatDuration formats a duration for display.
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
