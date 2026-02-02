package widgets

import (
	"fmt"

	"github.com/leo/visor/internal/config"
	"github.com/leo/visor/internal/input"
	"github.com/leo/visor/internal/render"
)

// APILatencyWidget displays the total API latency.
// This is a unique metric that no other statusline exposes.
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

	// Color based on latency
	color := "green"
	if ms >= 5000 {
		color = "red"
	} else if ms >= 2000 {
		color = "yellow"
	}

	return render.Colorize(text, color)
}

func (w *APILatencyWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return true
}
