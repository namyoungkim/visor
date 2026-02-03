package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/history"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// BlockTimerWidget displays remaining time in the 5-hour Claude Pro rate limit block.
//
// Supported Extra options:
//   - show_label: "true"/"false" - whether to show "Block:" prefix (default: true)
//   - warn_threshold: "80" - percentage elapsed for warning color (default: 80)
//   - critical_threshold: "95" - percentage elapsed for critical/red color (default: 95)
type BlockTimerWidget struct {
	history *history.History
}

func (w *BlockTimerWidget) Name() string {
	return "block_timer"
}

func (w *BlockTimerWidget) SetHistory(h *history.History) {
	w.history = h
}

func (w *BlockTimerWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	if w.history == nil {
		return ""
	}

	remainingMs := w.history.GetBlockRemainingMs()
	if remainingMs <= 0 {
		return ""
	}

	// Convert to hours and minutes
	remainingMinutes := remainingMs / 60000
	hours := remainingMinutes / 60
	minutes := remainingMinutes % 60

	// Format time
	var value string
	if hours > 0 {
		value = fmt.Sprintf("%dh%02dm", hours, minutes)
	} else {
		value = fmt.Sprintf("%dm", minutes)
	}

	// Build output with optional label
	var text string
	if cfg.Format != "" {
		text = FormatOutput(cfg, "", value)
	} else if GetExtraBool(cfg, "show_label", true) {
		text = "Block: " + value
	} else {
		text = value
	}

	// Determine color based on elapsed percentage
	elapsedPct := w.history.GetBlockElapsedPct()
	warnThreshold := GetExtraFloat(cfg, "warn_threshold", BlockTimerWarningPct)
	criticalThreshold := GetExtraFloat(cfg, "critical_threshold", BlockTimerCriticalPct)
	color := ColorByThreshold(elapsedPct, warnThreshold, criticalThreshold)

	return render.Colorize(text, color)
}

func (w *BlockTimerWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return w.history != nil && w.history.BlockStartTime > 0
}
