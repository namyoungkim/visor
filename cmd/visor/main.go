package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/namyoungkim/visor/internal/auth"
	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/cost"
	"github.com/namyoungkim/visor/internal/history"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
	"github.com/namyoungkim/visor/internal/transcript"
	"github.com/namyoungkim/visor/internal/tui"
	"github.com/namyoungkim/visor/internal/usage"
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
	tuiFlag := flag.Bool("tui", false, "Open interactive configuration editor")

	flag.Parse()

	// Handle flags
	if *versionFlag {
		fmt.Printf("visor %s\n", version)
		return
	}

	if *initFlag {
		// Get preset name from positional argument
		presetName := "default"
		if len(flag.Args()) > 0 {
			presetName = flag.Args()[0]
		}

		// Handle help command
		if presetName == "help" {
			fmt.Print(config.ListPresets())
			return
		}

		// Validate preset exists
		if _, ok := config.GetPreset(presetName); !ok {
			fmt.Fprintf(os.Stderr, "Unknown preset: %s\n\n", presetName)
			fmt.Print(config.ListPresets())
			os.Exit(1)
		}

		if err := config.InitWithPreset(presetName, ""); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created config at %s (preset: %s)\n", config.DefaultConfigPath(), presetName)
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

	if *tuiFlag {
		if err := tui.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
			os.Exit(1)
		}
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

	// Load cost data for daily/weekly/block cost widgets (v0.6)
	if cfg.Usage.Enabled {
		costData := loadCostData(session, hist, cfg, *debugFlag)
		widgets.SetCostData(costData)

		// Try to load usage limits from OAuth API
		limits := loadUsageLimits(*debugFlag)
		widgets.SetUsageLimits(limits)
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

// loadCostData loads aggregated cost data from JSONL transcripts.
func loadCostData(session *input.Session, hist *history.History, cfg *config.Config, debug bool) *cost.CostData {
	// Parse cost entries from the current session's transcript
	entries, err := cost.ParseSession(session.TranscriptPath)
	if err != nil && debug {
		fmt.Fprintf(os.Stderr, "[visor] cost parsing error: %v\n", err)
	}

	// Get block start time from history
	blockStart := hist.GetBlockStartTime()

	// Aggregate the data
	data := cost.Aggregate(entries, blockStart)

	// Set provider from config or auto-detect
	if cfg.Usage.Provider != "" {
		data.Provider = cost.Provider(cfg.Usage.Provider)
	}

	if debug {
		fmt.Fprintf(os.Stderr, "[visor] cost data: today=$%.2f, week=$%.2f, provider=%s\n",
			data.Today, data.Week, data.Provider)
	}

	return data
}

// loadUsageLimits loads usage limits from the OAuth API.
func loadUsageLimits(debug bool) *usage.Limits {
	provider := auth.DefaultProvider()
	client := usage.NewClient(provider)

	limits, err := client.GetLimits()
	if err != nil {
		if debug {
			fmt.Fprintf(os.Stderr, "[visor] usage API error: %v\n", err)
		}
		return nil
	}

	if debug {
		fmt.Fprintf(os.Stderr, "[visor] usage limits: 5h=%.0f%%, 7d=%.0f%%\n",
			limits.FiveHour.Utilization, limits.SevenDay.Utilization)
	}

	return limits
}
