package config

// Config represents the visor configuration.
type Config struct {
	General GeneralConfig `toml:"general"`
	Theme   ThemeConfig   `toml:"theme"`
	Usage   UsageConfig   `toml:"usage"`
	Lines   []Line        `toml:"line"`
}

// GeneralConfig contains global settings.
type GeneralConfig struct {
	Separator string `toml:"separator"`
}

// ThemeConfig contains theme settings.
type ThemeConfig struct {
	// Name is the theme preset name (default, powerline, gruvbox, nord, etc.)
	Name string `toml:"name"`

	// Powerline enables powerline-style rendering with arrows and backgrounds.
	// This is automatically true for themes ending in "-powerline".
	Powerline bool `toml:"powerline"`
}

// UsageConfig contains usage tracking settings.
type UsageConfig struct {
	// Enabled enables usage tracking features (cost aggregation, API limits).
	Enabled bool `toml:"enabled"`

	// Provider specifies the billing provider: "anthropic", "claude_pro", "aws", "gcp".
	// Auto-detected if empty.
	Provider string `toml:"provider"`

	// ProjectsDir is the path to Claude projects directory for JSONL parsing.
	// Defaults to ~/.claude/projects if empty.
	ProjectsDir string `toml:"projects_dir"`
}

// Line represents a single line in the statusline.
// Supports both single-side and split layout:
// - Single: widgets = ["model", "cost"]
// - Split: left = ["model", "git"], right = ["cost"]
type Line struct {
	Widgets []WidgetConfig `toml:"widget"`
	Left    []WidgetConfig `toml:"left"`
	Right   []WidgetConfig `toml:"right"`
}

// WidgetConfig represents configuration for a single widget.
type WidgetConfig struct {
	Name   string            `toml:"name"`
	Format string            `toml:"format"`
	Style  StyleConfig       `toml:"style"`
	Extra  map[string]string `toml:"extra"`
}

// StyleConfig contains ANSI styling options.
type StyleConfig struct {
	Fg   string `toml:"fg"`
	Bg   string `toml:"bg"`
	Bold bool   `toml:"bold"`
}
