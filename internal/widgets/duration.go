package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// DurationWidget displays the session duration.
//
// Supported Extra options:
//   - show_icon: "true"/"false" - show ⏱️ prefix (default: true)
//
// Output format: "⏱️ 5m" or "45s" or "1h23m"
type DurationWidget struct{}

func (w *DurationWidget) Name() string {
	return "duration"
}

func (w *DurationWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	ms := session.Cost.TotalDurationMs
	if ms <= 0 {
		return ""
	}

	duration := formatDurationMs(ms)
	showIcon := GetExtraBool(cfg, "show_icon", true)

	var text string
	if showIcon {
		text = "⏱️ " + duration
	} else {
		text = duration
	}

	return render.Colorize(text, "cyan")
}

func (w *DurationWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return session.Cost.TotalDurationMs > 0
}

// formatDurationMs converts milliseconds to a human-readable duration string.
// Returns formats like: "45s", "5m", "1h23m", "2h"
func formatDurationMs(ms int64) string {
	totalSeconds := ms / 1000

	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	switch {
	case hours > 0 && minutes > 0:
		return fmt.Sprintf("%dh%dm", hours, minutes)
	case hours > 0:
		return fmt.Sprintf("%dh", hours)
	case minutes > 0 && seconds > 0 && minutes < 5:
		// Show seconds only for short durations
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	case minutes > 0:
		return fmt.Sprintf("%dm", minutes)
	default:
		return fmt.Sprintf("%ds", seconds)
	}
}
