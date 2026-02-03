package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
	"github.com/namyoungkim/visor/internal/transcript"
)

// AgentsWidget displays the status of spawned sub-agents.
//
// Supported Extra options:
//   - show_label: "true"/"false" - show prefix (default: false)
//
// Output format:
//   - "◐ 1 agent" (1 running)
//   - "✓ 2 done" (all completed)
//   - "◐ 1 | ✓ 2" (1 running, 2 done)
type AgentsWidget struct {
	transcript *transcript.Data
}

func (w *AgentsWidget) Name() string {
	return "agents"
}

// SetTranscript sets the transcript data for this widget.
func (w *AgentsWidget) SetTranscript(t *transcript.Data) {
	w.transcript = t
}

func (w *AgentsWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	if w.transcript == nil || len(w.transcript.Agents) == 0 {
		return ""
	}

	running := 0
	completed := 0
	for _, agent := range w.transcript.Agents {
		if agent.Status == "running" {
			running++
		} else {
			completed++
		}
	}

	var text string
	if running > 0 && completed > 0 {
		text = fmt.Sprintf("◐ %d | ✓ %d", running, completed)
	} else if running > 0 {
		text = fmt.Sprintf("◐ %d agent", running)
		if running > 1 {
			text += "s"
		}
	} else {
		text = fmt.Sprintf("✓ %d done", completed)
	}

	if GetExtraBool(cfg, "show_label", false) {
		text = "Agents: " + text
	}

	// Color based on status
	color := "green"
	if running > 0 {
		color = "yellow"
	}

	return render.Colorize(text, color)
}

func (w *AgentsWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return w.transcript != nil && len(w.transcript.Agents) > 0
}
