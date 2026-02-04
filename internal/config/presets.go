package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Preset represents a configuration preset.
type Preset struct {
	Name        string
	Description string
	Lines       [][]string // Each inner slice is a line of widgets
}

// Presets contains all available presets.
var Presets = map[string]Preset{
	"minimal": {
		Name:        "minimal",
		Description: "Essential widgets only (4 widgets)",
		Lines: [][]string{
			{"model", "context", "cost", "git"},
		},
	},
	"default": {
		Name:        "default",
		Description: "Balanced default with visor metrics (6 widgets)",
		Lines: [][]string{
			{"model", "context", "cache_hit", "api_latency", "cost", "git"},
		},
	},
	"efficiency": {
		Name:        "efficiency",
		Description: "Cost optimization focus (6 widgets)",
		Lines: [][]string{
			{"model", "context", "burn_rate", "cache_hit", "compact_eta", "cost"},
		},
	},
	"developer": {
		Name:        "developer",
		Description: "Tool/agent monitoring (6 widgets)",
		Lines: [][]string{
			{"model", "context", "tools", "agents", "code_changes", "git"},
		},
	},
	"pro": {
		Name:        "pro",
		Description: "Claude Pro rate limits (6 widgets)",
		Lines: [][]string{
			{"model", "context", "block_limit", "week_limit", "daily_cost", "cost"},
		},
	},
	"full": {
		Name:        "full",
		Description: "All widgets, multi-line layout (18 widgets)",
		Lines: [][]string{
			{"model", "context", "cost", "git"},
			{"cache_hit", "api_latency", "burn_rate", "compact_eta", "context_spark"},
			{"tools", "code_changes"},
			{"agents"},
			{"block_timer", "block_limit", "week_limit", "daily_cost", "weekly_cost", "block_cost"},
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

`, p.Name, p.Description))

	// Generate widget configuration for each line
	for lineIdx, widgets := range p.Lines {
		if lineIdx > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString("[[line]]\n")
		for _, widget := range widgets {
			sb.WriteString(fmt.Sprintf("  [[line.widget]]\n  name = %q\n", widget))
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
