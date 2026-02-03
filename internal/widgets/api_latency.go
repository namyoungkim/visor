package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// APILatencyWidget displays the total API latency.
// This is a unique metric that no other statusline exposes.
//
// Supported Extra options:
//   - warn_threshold: "2000" - milliseconds for warning color (default: 2000)
//   - critical_threshold: "5000" - milliseconds for critical/red color (default: 5000)
type APILatencyWidget struct{}

func (w *APILatencyWidget) Name() string {
	return "api_latency"
}

func (w *APILatencyWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	ms := session.Cost.TotalAPIDurationMs

	if ms == 0 {
		return render.Colorize("API: —", "gray")
	}

	var text string
	if ms >= 1000 {
		// Convert to seconds
		secs := float64(ms) / 1000.0
		text = fmt.Sprintf("API: %.1fs", secs)
	} else {
		text = fmt.Sprintf("API: %dms", ms)
	}

	warnThreshold := GetExtraFloat(cfg, "warn_threshold", LatencyWarningMs)
	criticalThreshold := GetExtraFloat(cfg, "critical_threshold", LatencyDangerMs)
	color := ColorByThreshold(float64(ms), warnThreshold, criticalThreshold)
	return render.Colorize(text, color)
}

func (w *APILatencyWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return true
}
