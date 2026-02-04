package widgets

import (
	"fmt"
	"strings"
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
//   - show_bar: "true"/"false" - show progress bar (default: false)
//   - bar_width: "10" - progress bar width in characters (default: 10)
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
	pctStr := fmt.Sprintf("%.0f%%", pct)

	// Build value with optional progress bar
	var valueParts []string
	valueParts = append(valueParts, pctStr)

	if GetExtraBool(cfg, "show_bar", false) {
		barWidth := GetExtraInt(cfg, "bar_width", DefaultBarWidth)
		bar := ProgressBar(pct, barWidth)
		valueParts = append(valueParts, bar)
	}

	if GetExtraBool(cfg, "show_remaining", true) {
		remaining := w.limits.FiveHourRemaining()
		if remaining > 0 {
			valueParts = append(valueParts, fmt.Sprintf("(%s)", formatDuration(remaining)))
		}
	}

	value := strings.Join(valueParts, " ")

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
