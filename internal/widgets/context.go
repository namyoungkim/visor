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
type ContextWidget struct{}

func (w *ContextWidget) Name() string {
	return "context"
}

func (w *ContextWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	pct := session.ContextWindow.UsedPercentage
	color := ColorByThreshold(pct, ContextWarningPct, ContextDangerPct)

	value := fmt.Sprintf("%.0f%%", pct)

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
