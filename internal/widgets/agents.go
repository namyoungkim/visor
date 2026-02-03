package widgets

import (
	"strings"
	"time"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
	"github.com/namyoungkim/visor/internal/transcript"
)

// AgentsWidget displays the status of spawned sub-agents with details.
//
// Supported Extra options:
//   - max_display: maximum number of agents to show (default: "3")
//   - show_label: "true"/"false" - show prefix (default: false)
//   - show_description: "true"/"false" - show task description (default: true)
//   - show_duration: "true"/"false" - show elapsed time (default: true)
//   - max_description_len: max length for description (default: "20")
//
// Output format: "Explore: Analyze widgets (42s)" (with description and duration)
// Running agents show "..." instead of duration
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
	showDescription := GetExtraBool(cfg, "show_description", true)
	showDuration := GetExtraBool(cfg, "show_duration", true)
	maxDescLen := GetExtraInt(cfg, "max_description_len", 20)
	agents := w.transcript.Agents

	// Show only the last N agents
	start := 0
	if len(agents) > maxDisplay {
		start = len(agents) - maxDisplay
	}

	var parts []string
	for _, agent := range agents[start:] {
		part := w.renderAgent(agent, showDescription, showDuration, maxDescLen)
		parts = append(parts, part)
	}

	text := strings.Join(parts, " | ")

	if GetExtraBool(cfg, "show_label", false) {
		text = "Agents: " + text
	}

	return text
}

// renderAgent renders a single agent with optional description and duration.
func (w *AgentsWidget) renderAgent(agent transcript.Agent, showDesc, showDur bool, maxDescLen int) string {
	icon, color := agentStatusIcon(agent.Status)

	// If no description or duration options, use simple format
	if !showDesc && !showDur {
		return render.Colorize(icon+agent.Type, color)
	}

	// Build detailed format: "Type: Description (Ns)" or "Type: Description (...)"
	var result string

	// Type with icon
	result = render.Colorize(icon+agent.Type, color)

	// Add description if available and enabled
	if showDesc && agent.Description != "" {
		desc := truncateString(agent.Description, maxDescLen)
		result += ": " + render.Colorize(desc, "dim")
	}

	// Add duration if enabled
	if showDur {
		if agent.Status == "running" {
			result += render.Colorize(" (...)", "dim")
		} else if agent.EndTime > 0 && agent.StartTime > 0 {
			durationSec := (agent.EndTime - agent.StartTime) / 1000
			result += render.Colorize(" ("+formatDurationSec(durationSec)+")", "dim")
		}
	}

	return result
}

// truncateString truncates a string to maxLen, adding "..." if truncated.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return "..."
	}
	return s[:maxLen-3] + "..."
}

// formatDurationSec formats seconds into a human-readable duration string.
func formatDurationSec(seconds int64) string {
	if seconds < 60 {
		return itoa(int(seconds)) + "s"
	}
	minutes := seconds / 60
	if minutes < 60 {
		return itoa(int(minutes)) + "m"
	}
	hours := minutes / 60
	mins := minutes % 60
	if mins == 0 {
		return itoa(int(hours)) + "h"
	}
	return itoa(int(hours)) + "h" + itoa(int(mins)) + "m"
}

// nowUnixMilli returns the current time in Unix milliseconds.
// This is used for calculating running duration.
var nowUnixMilli = func() int64 {
	return time.Now().UnixMilli()
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
