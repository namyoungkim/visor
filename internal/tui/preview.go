package tui

import (
	"strings"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
	"github.com/namyoungkim/visor/internal/widgets"
)

// SampleSession returns a sample session for preview rendering
func SampleSession() *input.Session {
	return &input.Session{
		SessionID: "preview-session",
		Model: input.Model{
			DisplayName: "Opus",
			ID:          "claude-3-opus",
		},
		Cost: input.Cost{
			TotalCostUSD:          0.48,
			TotalDurationMs:       45000,
			TotalAPIDurationMs:    2500,
			TotalAPICalls:         12,
			TotalInputTokens:      15000,
			TotalOutputTokens:     3000,
			TotalCacheReadTokens:  12000,
			TotalCacheWriteTokens: 500,
		},
		ContextWindow: input.ContextWindow{
			UsedPercentage: 42.5,
			UsedTokens:     42500,
			MaxTokens:      100000,
		},
		Workspace: input.Workspace{
			LinesAdded:   156,
			LinesRemoved: 42,
			FilesChanged: 8,
		},
		CurrentUsage: &input.CurrentUsage{
			InputTokens:     5000,
			CacheReadTokens: 4000,
		},
	}
}

// RenderPreview renders a preview of the current configuration
func RenderPreview(cfg *config.Config) string {
	session := SampleSession()

	var lines []string
	for _, line := range cfg.Lines {
		var lineOutput string

		// Check if this is a split layout
		if len(line.Left) > 0 || len(line.Right) > 0 {
			leftRendered := widgets.RenderAll(session, line.Left)
			rightRendered := widgets.RenderAll(session, line.Right)
			lineOutput = render.SplitLayout(leftRendered, rightRendered, cfg.General.Separator)
		} else {
			rendered := widgets.RenderAll(session, line.Widgets)
			lineOutput = render.Layout(rendered, cfg.General.Separator)
		}

		if lineOutput != "" {
			lines = append(lines, lineOutput)
		}
	}

	if len(lines) == 0 {
		return "(no widgets enabled)"
	}

	return strings.Join(lines, "\n")
}

// RenderWidgetPreview renders a single widget for preview
func RenderWidgetPreview(widgetCfg *config.WidgetConfig) string {
	session := SampleSession()

	w, ok := widgets.Get(widgetCfg.Name)
	if !ok {
		return "(unknown widget)"
	}

	if !w.ShouldRender(session, widgetCfg) {
		return "(hidden)"
	}

	result := w.Render(session, widgetCfg)
	if result == "" {
		return "(empty)"
	}

	return result
}
