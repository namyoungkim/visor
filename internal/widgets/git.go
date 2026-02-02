package widgets

import (
	"fmt"
	"strings"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/git"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// GitWidget displays git branch and status.
type GitWidget struct{}

func (w *GitWidget) Name() string {
	return "git"
}

func (w *GitWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	status := git.GetStatus()
	if !status.IsRepo {
		return ""
	}

	var parts []string

	// Branch name
	branch := status.Branch
	if branch == "" {
		branch = "HEAD"
	}
	parts = append(parts, render.Colorize(branch, "magenta"))

	// Status indicators
	var indicators []string

	if status.Staged > 0 {
		indicators = append(indicators, render.Colorize(fmt.Sprintf("+%d", status.Staged), "green"))
	}
	if status.Modified > 0 {
		indicators = append(indicators, render.Colorize(fmt.Sprintf("~%d", status.Modified), "yellow"))
	}
	if status.Ahead > 0 {
		indicators = append(indicators, render.Colorize(fmt.Sprintf("↑%d", status.Ahead), "cyan"))
	}
	if status.Behind > 0 {
		indicators = append(indicators, render.Colorize(fmt.Sprintf("↓%d", status.Behind), "red"))
	}

	if len(indicators) > 0 {
		parts = append(parts, strings.Join(indicators, ""))
	}

	return strings.Join(parts, " ")
}

func (w *GitWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return git.GetStatus().IsRepo
}
