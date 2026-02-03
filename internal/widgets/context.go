package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// ContextWidget displays context window usage percentage.
//
// Supported Extra options:
//   - show_label: "true"/"false" - whether to show "Ctx:" prefix (default: true)
//   - show_bar: "true"/"false" - whether to show progress bar (default: true)
//   - bar_width: "10" - progress bar width in characters (default: 10)
//   - warn_threshold: "60" - percentage for warning color (default: 60)
//   - critical_threshold: "80" - percentage for critical/red color (default: 80)
type ContextWidget struct{}

func (w *ContextWidget) Name() string {
	return "context"
}

func (w *ContextWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	pct := session.ContextWindow.UsedPercentage
	warnThreshold := GetExtraFloat(cfg, "warn_threshold", ContextWarningPct)
	criticalThreshold := GetExtraFloat(cfg, "critical_threshold", ContextDangerPct)
	color := ColorByThreshold(pct, warnThreshold, criticalThreshold)

	pctStr := fmt.Sprintf("%.0f%%", pct)

	// Build value with optional progress bar
	var value string
	if GetExtraBool(cfg, "show_bar", true) {
		barWidth := GetExtraInt(cfg, "bar_width", DefaultBarWidth)
		bar := ProgressBar(pct, barWidth)
		value = pctStr + " " + bar
	} else {
		value = pctStr
	}

	var text string
	if cfg.Format != "" {
		text = FormatOutput(cfg, "", value)
	} else if GetExtraBool(cfg, "show_label", true) {
		text = "Ctx: " + value
	} else {
		text = value
	}

	return render.Colorize(text, color)
}

func (w *ContextWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return true
}
