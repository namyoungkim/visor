package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// ContextWidget displays context window usage percentage.
type ContextWidget struct{}

func (w *ContextWidget) Name() string {
	return "context"
}

func (w *ContextWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	pct := session.ContextWindow.UsedPercentage
	color := ColorByThreshold(pct, ContextWarningPct, ContextDangerPct)
	text := fmt.Sprintf("Ctx: %.0f%%", pct)
	return render.Colorize(text, color)
}

func (w *ContextWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return true
}
