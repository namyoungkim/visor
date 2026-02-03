package widgets

import (
	"strings"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
	"github.com/namyoungkim/visor/internal/transcript"
)

// AgentsWidget displays the status of spawned sub-agents.
//
// Supported Extra options:
//   - max_display: maximum number of agents to show (default: "3")
//   - show_label: "true"/"false" - show prefix (default: false)
//
// Output format: "✓Plan ✓Bash ◐Explore" (completed Plan/Bash, running Explore)
// Status icons: ✓ (completed), ◐ (running)
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

	maxDisplay := GetExtraInt(cfg, "max_display", 3)
	agents := w.transcript.Agents

	// Show only the last N agents
	start := 0
	if len(agents) > maxDisplay {
		start = len(agents) - maxDisplay
	}

	var parts []string
	for _, agent := range agents[start:] {
		icon, color := agentStatusIcon(agent.Status)
		parts = append(parts, render.Colorize(icon+agent.Type, color))
	}

	text := strings.Join(parts, " ")

	if GetExtraBool(cfg, "show_label", false) {
		text = "Agents: " + text
	}

	return text
}

// agentStatusIcon returns the icon and color for an agent status.
func agentStatusIcon(status string) (string, string) {
	if status == "running" {
		return "◐", "yellow"
	}
	return "✓", "green"
}

func (w *AgentsWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return w.transcript != nil && len(w.transcript.Agents) > 0
}
