package widgets

import (
	"strings"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
	"github.com/namyoungkim/visor/internal/transcript"
)

// ToolsWidget displays recent tool invocations with their status and count.
//
// Supported Extra options:
//   - max_display: maximum number of tools to show, 0 = unlimited (default: "0")
//   - show_label: "true"/"false" - show prefix (default: false)
//   - show_count: "true"/"false" - show invocation count (default: true)
//
// Output format: "✓Bash ×7 | ✓Edit ×4 | ✓Read ×6" (with counts)
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

	maxDisplay := GetExtraInt(cfg, "max_display", 0) // 0 = unlimited
	showCount := GetExtraBool(cfg, "show_count", true)
	tools := w.transcript.Tools

	// Show only the last N tools (0 = show all)
	start := 0
	if maxDisplay > 0 && len(tools) > maxDisplay {
		start = len(tools) - maxDisplay
	}

	var parts []string
	for _, tool := range tools[start:] {
		icon, color := toolStatusIcon(tool.Status)
		part := render.Colorize(icon+tool.Name, color) + countSuffix(showCount, tool.Count)
		parts = append(parts, part)
	}

	text := strings.Join(parts, " | ")

	if GetExtraBool(cfg, "show_label", false) {
		text = "Tools: " + text
	}

	return text
}

// countSuffix returns the count suffix (e.g., " ×7") if show_count is enabled and count > 1.
func countSuffix(showCount bool, count int) string {
	if showCount && count > 1 {
		return render.Colorize(" ×"+itoa(count), "dim")
	}
	return ""
}

// itoa converts an int to a string without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
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

