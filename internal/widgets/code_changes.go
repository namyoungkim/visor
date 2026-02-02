package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// CodeChangesWidget displays lines added/removed in the session.
// This is a unique metric that no other statusline exposes.
type CodeChangesWidget struct{}

func (w *CodeChangesWidget) Name() string {
	return "code_changes"
}

func (w *CodeChangesWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	added := session.Workspace.LinesAdded
	removed := session.Workspace.LinesRemoved

	if added == 0 && removed == 0 {
		return ""
	}

	addedStr := render.Colorize(fmt.Sprintf("+%d", added), "green")
	removedStr := render.Colorize(fmt.Sprintf("-%d", removed), "red")

	return addedStr + "/" + removedStr
}

func (w *CodeChangesWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return session.Workspace.LinesAdded > 0 || session.Workspace.LinesRemoved > 0
}
