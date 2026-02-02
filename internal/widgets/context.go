package widgets

import (
	"fmt"

	"github.com/leo/visor/internal/config"
	"github.com/leo/visor/internal/input"
	"github.com/leo/visor/internal/render"
)

// ContextWidget displays context window usage percentage.
type ContextWidget struct{}

func (w *ContextWidget) Name() string {
	return "context"
}

func (w *ContextWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	pct := session.ContextWindow.UsedPercentage

	// Choose color based on usage level
	color := "green"
	if pct >= 80 {
		color = "red"
	} else if pct >= 60 {
		color = "yellow"
	}

	text := fmt.Sprintf("Ctx: %.0f%%", pct)
	return render.Colorize(text, color)
}

func (w *ContextWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return true
}
