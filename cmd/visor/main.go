package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/history"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
	"github.com/namyoungkim/visor/internal/transcript"
	"github.com/namyoungkim/visor/internal/widgets"
)

// version is set by ldflags during build, defaults to "dev" for local builds
var version = "dev"

func main() {
	// Define flags
	versionFlag := flag.Bool("version", false, "Print version information")
	initFlag := flag.Bool("init", false, "Generate default config at ~/.config/visor/config.toml")
	setupFlag := flag.Bool("setup", false, "Configure Claude Code to use visor statusline")
	checkFlag := flag.Bool("check", false, "Validate configuration file")
	debugFlag := flag.Bool("debug", false, "Enable debug output to stderr")

	flag.Parse()

	// Handle flags
	if *versionFlag {
		fmt.Printf("visor %s\n", version)
		return
	}

	if *initFlag {
		if err := config.Init(""); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created config at %s\n", config.DefaultConfigPath())
		return
	}

	if *setupFlag {
		printSetupInstructions()
		return
	}

	if *checkFlag {
		if err := config.Validate(""); err != nil {
			fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Config is valid")
		return
	}

	// Main pipeline: stdin → parse → render → stdout
	session := input.Parse(os.Stdin)
	cfg, err := config.Load("")
	if err != nil {
		if *debugFlag {
			fmt.Fprintf(os.Stderr, "[visor] config error: %v\n", err)
		}
		return
	}

	if *debugFlag {
		fmt.Fprintf(os.Stderr, "[visor] config loaded from: %s\n", config.DefaultConfigPath())
		fmt.Fprintf(os.Stderr, "[visor] session model: %s\n", session.Model.DisplayName)
	}

	// Load history for this session
	hist, err := history.Load(session.SessionID)
	if err != nil && *debugFlag {
		fmt.Fprintf(os.Stderr, "[visor] history error: %v\n", err)
	}

	// Update block timer (for Claude Pro rate limit tracking)
	hist.UpdateBlockStartTime()

	// Calculate cache hit rate for history
	var cacheHitPct float64
	if session.CurrentUsage != nil {
		total := session.CurrentUsage.InputTokens + session.CurrentUsage.CacheReadTokens
		if total > 0 {
			cacheHitPct = float64(session.CurrentUsage.CacheReadTokens) / float64(total) * 100
		}
	}

	// Add current session data to history
	hist.Add(history.Entry{
		ContextPct:   session.ContextWindow.UsedPercentage,
		CostUSD:      session.Cost.TotalCostUSD,
		DurationMs:   session.Cost.TotalDurationMs,
		CacheHitPct:  cacheHitPct,
		APILatencyMs: session.Cost.TotalAPIDurationMs,
	})

	// Set history on context_spark widget
	widgets.SetHistory(hist)

	// Load transcript for tools/agents widgets
	transcriptData := transcript.Parse(session.TranscriptPath)
	widgets.SetTranscript(transcriptData)

	if *debugFlag {
		fmt.Fprintf(os.Stderr, "[visor] transcript path: %s\n", session.TranscriptPath)
		fmt.Fprintf(os.Stderr, "[visor] transcript tools: %d, agents: %d\n", len(transcriptData.Tools), len(transcriptData.Agents))
	}

	output := renderSession(session, cfg)
	if output != "" {
		fmt.Print(output)
	}

	// Save history
	if err := hist.Save(); err != nil && *debugFlag {
		fmt.Fprintf(os.Stderr, "[visor] failed to save history: %v\n", err)
	}
}

func renderSession(session *input.Session, cfg *config.Config) string {
	var result []string

	for _, line := range cfg.Lines {
		var lineOutput string

		// Check if this is a split layout (left/right defined)
		if len(line.Left) > 0 || len(line.Right) > 0 {
			leftRendered := widgets.RenderAll(session, line.Left)
			rightRendered := widgets.RenderAll(session, line.Right)
			lineOutput = render.SplitLayout(leftRendered, rightRendered, cfg.General.Separator)
		} else {
			// Regular layout
			rendered := widgets.RenderAll(session, line.Widgets)
			lineOutput = render.Layout(rendered, cfg.General.Separator)
		}

		if lineOutput != "" {
			result = append(result, lineOutput)
		}
	}

	if len(result) == 0 {
		return ""
	}

	return render.JoinLines(result)
}

func printSetupInstructions() {
	fmt.Println(`To configure Claude Code to use visor:

1. Add to your Claude Code settings (~/.claude/settings.json):

{
  "statusline": {
    "command": "visor"
  }
}

2. Or set the environment variable:

export CLAUDE_STATUSLINE_COMMAND="visor"

3. Optionally customize with: visor --init`)
}
