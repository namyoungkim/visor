package widgets

import (
	"fmt"
	"os"
	"strings"

	"github.com/namyoungkim/visor/internal/claudeconfig"
	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// ConfigCountsWidget displays Claude configuration counts.
//
// Supported Extra options:
//   - show_claude_md: "true"/"false" - show CLAUDE.md count (default: true)
//   - show_rules: "true"/"false" - show permission rules count (default: true)
//   - show_mcps: "true"/"false" - show MCP plugins count (default: true)
//   - show_hooks: "true"/"false" - show hooks count (default: true)
//
// Output format: "2 CLAUDE.md | 3 rules | 2 MCPs | 1 hook"
type ConfigCountsWidget struct {
	counts *claudeconfig.Counts
}

func (w *ConfigCountsWidget) Name() string {
	return "config_counts"
}

// SetCounts sets the config counts data for this widget.
func (w *ConfigCountsWidget) SetCounts(counts *claudeconfig.Counts) {
	w.counts = counts
}

func (w *ConfigCountsWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	// Load counts lazily if not set
	if w.counts == nil {
		cwd, _ := os.Getwd()
		w.counts = claudeconfig.LoadCounts(cwd)
	}

	showClaudeMD := GetExtraBool(cfg, "show_claude_md", true)
	showRules := GetExtraBool(cfg, "show_rules", true)
	showMCPs := GetExtraBool(cfg, "show_mcps", true)
	showHooks := GetExtraBool(cfg, "show_hooks", true)

	var parts []string

	if showClaudeMD && w.counts.ClaudeMDCount > 0 {
		label := "CLAUDE.md"
		if w.counts.ClaudeMDCount > 1 {
			label = "CLAUDE.mds"
		}
		parts = append(parts, render.Colorize(fmt.Sprintf("%d %s", w.counts.ClaudeMDCount, label), "cyan"))
	}

	if showRules && w.counts.RulesCount > 0 {
		label := "rule"
		if w.counts.RulesCount > 1 {
			label = "rules"
		}
		parts = append(parts, render.Colorize(fmt.Sprintf("%d %s", w.counts.RulesCount, label), "green"))
	}

	if showMCPs && w.counts.MCPCount > 0 {
		label := "MCP"
		if w.counts.MCPCount > 1 {
			label = "MCPs"
		}
		parts = append(parts, render.Colorize(fmt.Sprintf("%d %s", w.counts.MCPCount, label), "magenta"))
	}

	if showHooks && w.counts.HooksCount > 0 {
		label := "hook"
		if w.counts.HooksCount > 1 {
			label = "hooks"
		}
		parts = append(parts, render.Colorize(fmt.Sprintf("%d %s", w.counts.HooksCount, label), "yellow"))
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, " | ")
}

func (w *ConfigCountsWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	// Load counts lazily if not set
	if w.counts == nil {
		cwd, _ := os.Getwd()
		w.counts = claudeconfig.LoadCounts(cwd)
	}

	// Only render if at least one count is non-zero
	return w.counts.ClaudeMDCount > 0 ||
		w.counts.RulesCount > 0 ||
		w.counts.MCPCount > 0 ||
		w.counts.HooksCount > 0
}
