package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// APILatencyWidget displays the per-call average API latency.
// This is a unique metric that no other statusline exposes.
//
// Formula: total_api_duration_ms / total_api_calls
//
// Supported Extra options:
//   - warn_threshold: "2000" - milliseconds for warning color (default: 2000)
//   - critical_threshold: "5000" - milliseconds for critical/red color (default: 5000)
type APILatencyWidget struct{}

func (w *APILatencyWidget) Name() string {
	return "api_latency"
}

func (w *APILatencyWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	totalMs := session.Cost.TotalAPIDurationMs
	calls := session.Cost.TotalAPICalls

	if calls <= 0 || totalMs <= 0 {
		return render.Colorize("API: —", "gray")
	}

	ms := totalMs / int64(calls)

	var text string
	if ms >= 1000 {
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
