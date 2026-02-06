package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WidgetDef represents a widget definition with optional extra settings.
type WidgetDef struct {
	Name  string
	Extra map[string]string
}

// W is a helper to create a WidgetDef with just a name.
func W(name string) WidgetDef {
	return WidgetDef{Name: name}
}

// WL is a helper to create a WidgetDef with show_label=true.
func WL(name string) WidgetDef {
	return WidgetDef{Name: name, Extra: map[string]string{"show_label": "true"}}
}

// Preset represents a configuration preset.
type Preset struct {
	Name        string
	Description string
	Lines       [][]WidgetDef // Each inner slice is a line of widgets
}

// Presets contains all available presets.
var Presets = map[string]Preset{
	"minimal": {
		Name:        "minimal",
		Description: "Essential widgets only (4 widgets)",
		Lines: [][]WidgetDef{
			{W("model"), W("context"), W("cost"), W("git")},
		},
	},
	"default": {
		Name:        "default",
		Description: "Balanced default with visor metrics (6 widgets)",
		Lines: [][]WidgetDef{
			{W("model"), W("context"), W("cache_hit"), W("api_latency"), W("cost"), W("git")},
		},
	},
	"efficiency": {
		Name:        "efficiency",
		Description: "Cost optimization focus (6 widgets)",
		Lines: [][]WidgetDef{
			{W("model"), W("context"), W("burn_rate"), W("cache_hit"), W("compact_eta"), W("cost")},
		},
	},
	"developer": {
		Name:        "developer",
		Description: "Tool/agent monitoring (7 widgets)",
		Lines: [][]WidgetDef{
			{W("model"), W("context"), W("tools"), {Name: "agents", Extra: map[string]string{"show_description": "false"}}, W("todos"), W("code_changes"), W("git")},
		},
	},
	"pro": {
		Name:        "pro",
		Description: "Claude Pro rate limits (6 widgets)",
		Lines: [][]WidgetDef{
			{W("model"), W("context"), W("block_limit"), W("week_limit"), W("daily_cost"), W("cost")},
		},
	},
	"full": {
		Name:        "full",
		Description: "All widgets, multi-line layout (24 widgets)",
		Lines: [][]WidgetDef{
			// Line 1: Session identity
			{W("model"), WL("plan"), W("session_id")},
			// Line 2: Core metrics
			{W("context"), W("duration"), W("cost"), W("git")},
			// Line 3: Tools (dynamic width)
			{W("tools")},
			// Line 4: Agents (dynamic width)
			{W("agents")},
			// Line 5: Efficiency metrics
			{WL("cache_hit"), W("api_latency"), W("token_speed"), WL("burn_rate")},
			// Line 6: Tracking
			{W("context_spark"), WL("compact_eta"), W("todos"), W("code_changes"), W("config_counts")},
			// Line 7: Cost & rate limits
			{WL("block_timer"), {Name: "block_limit", Extra: map[string]string{"show_label": "true", "show_bar": "true"}}, WL("week_limit"), WL("daily_cost"), WL("weekly_cost"), WL("block_cost")},
		},
	},
}

// PresetOrder defines the display order for listing presets.
var PresetOrder = []string{"minimal", "default", "efficiency", "developer", "pro", "full"}

// GetPreset returns a preset by name.
func GetPreset(name string) (Preset, bool) {
	p, ok := Presets[name]
	return p, ok
}

// ListPresets returns a formatted string of all available presets.
func ListPresets() string {
	var sb strings.Builder
	sb.WriteString("Available presets:\n\n")

	for _, name := range PresetOrder {
		p := Presets[name]
		sb.WriteString(fmt.Sprintf("  %-12s %s\n", name, p.Description))
	}

	sb.WriteString("\nUsage:\n")
	sb.WriteString("  visor --init           # Use 'default' preset\n")
	sb.WriteString("  visor --init minimal   # Use specific preset\n")
	sb.WriteString("  visor --init help      # Show this help\n")

	return sb.String()
}

// GetPresetTOML generates a TOML configuration string for a preset.
func GetPresetTOML(name string) (string, error) {
	p, ok := Presets[name]
	if !ok {
		return "", fmt.Errorf("unknown preset: %s", name)
	}

	var sb strings.Builder

	// Header comment
	sb.WriteString(fmt.Sprintf(`# visor configuration
# Preset: %s - %s
# Place at ~/.config/visor/config.toml

[general]
separator = " | "  # Widget separator

[theme]
name = "default"   # Theme: default, powerline, gruvbox, nord, gruvbox-powerline, nord-powerline

[usage]
enabled = true     # Enable usage tracking (daily/weekly cost, rate limits)
# five_hour_limit = 0   # 0 = auto-detect from subscription tier (Pro: 45, Max 5x: 225, Max 20x: 900)
# seven_day_limit = 0   # 0 = auto-detect from subscription tier

`, p.Name, p.Description))

	// Generate widget configuration for each line
	for lineIdx, widgets := range p.Lines {
		if lineIdx > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString("[[line]]\n")
		for _, widget := range widgets {
			sb.WriteString(fmt.Sprintf("  [[line.widget]]\n  name = %q\n", widget.Name))
			if len(widget.Extra) > 0 {
				sb.WriteString("  [line.widget.extra]\n")
				for k, v := range widget.Extra {
					sb.WriteString(fmt.Sprintf("  %s = %q\n", k, v))
				}
			}
		}
	}

	return sb.String(), nil
}

// InitWithPreset creates a configuration file using the specified preset.
// If preset is empty, defaults to "default" for API safety (CLI already validates).
// If path is empty, uses DefaultConfigPath().
func InitWithPreset(preset, path string) error {
	if path == "" {
		path = DefaultConfigPath()
	}

	// Default to "default" preset if empty (defensive check for direct API calls)
	if preset == "" {
		preset = "default"
	}

	content, err := GetPresetTOML(preset)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0644)
}
