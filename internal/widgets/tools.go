package widgets

import (
	"strings"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
	"github.com/namyoungkim/visor/internal/transcript"
)

// ToolsWidget displays recent tool invocations with their status.
//
// Supported Extra options:
//   - max_display: maximum number of tools to show (default: "3")
//   - show_label: "true"/"false" - show prefix (default: false)
//
// Output format: "✓Read ✓Write ◐Bash" (completed Read, Write; running Bash)
// Status icons: ✓ (completed), ✗ (error), ◐ (running)
type ToolsWidget struct {
	transcript *transcript.Data
}

func (w *ToolsWidget) Name() string {
	return "tools"
}

// SetTranscript sets the transcript data for this widget.
func (w *ToolsWidget) SetTranscript(t *transcript.Data) {
	w.transcript = t
}

func (w *ToolsWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	if w.transcript == nil || len(w.transcript.Tools) == 0 {
		return ""
	}

	maxDisplay := GetExtraInt(cfg, "max_display", 3)
	tools := w.transcript.Tools

	// Show only the last N tools
	start := 0
	if len(tools) > maxDisplay {
		start = len(tools) - maxDisplay
	}

	var parts []string
	for _, tool := range tools[start:] {
		icon, color := toolStatusIcon(tool.Status)
		parts = append(parts, render.Colorize(icon+tool.Name, color))
	}

	text := strings.Join(parts, " ")

	if GetExtraBool(cfg, "show_label", false) {
		text = "Tools: " + text
	}

	return text
}

func (w *ToolsWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return w.transcript != nil && len(w.transcript.Tools) > 0
}

// toolStatusIcon returns the icon and color for a tool status.
func toolStatusIcon(status transcript.ToolStatus) (string, string) {
	switch status {
	case transcript.ToolCompleted:
		return "✓", "green"
	case transcript.ToolError:
		return "✗", "red"
	case transcript.ToolRunning:
		return "◐", "yellow"
	default:
		return "?", "dim"
	}
}

