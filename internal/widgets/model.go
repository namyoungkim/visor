package widgets

import (
	"github.com/leo/visor/internal/config"
	"github.com/leo/visor/internal/input"
	"github.com/leo/visor/internal/render"
)

// ModelWidget displays the current model name.
type ModelWidget struct{}

func (w *ModelWidget) Name() string {
	return "model"
}

func (w *ModelWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	name := session.Model.DisplayName
	if name == "" {
		name = session.Model.ID
	}
	if name == "" {
		return ""
	}

	return render.Colorize(name, "cyan")
}

func (w *ModelWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return session.Model.DisplayName != "" || session.Model.ID != ""
}
