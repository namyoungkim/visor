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
	Debug     bool   `toml:"debug"`
}

// ThemeConfig contains theme settings.
type ThemeConfig struct {
	// Name is the theme preset name (default, powerline, gruvbox, nord, etc.)
	Name string `toml:"name"`

	// Powerline enables powerline-style rendering with arrows and backgrounds.
	// This is automatically true for themes ending in "-powerline".
	Powerline bool `toml:"powerline"`

	// Colors allows overriding individual colors from the preset.
	Colors *ColorOverrides `toml:"colors,omitempty"`

	// Separators allows overriding separator characters from the preset.
	Separators *SeparatorOverrides `toml:"separators,omitempty"`
}

// ColorOverrides allows overriding individual colors from a theme preset.
// Empty strings mean "use preset value".
type ColorOverrides struct {
	Normal      string   `toml:"normal,omitempty"`
	Warning     string   `toml:"warning,omitempty"`
	Critical    string   `toml:"critical,omitempty"`
	Good        string   `toml:"good,omitempty"`
	Primary     string   `toml:"primary,omitempty"`
	Secondary   string   `toml:"secondary,omitempty"`
	Muted       string   `toml:"muted,omitempty"`
	Backgrounds []string `toml:"backgrounds,omitempty"`
}

// SeparatorOverrides allows overriding separator characters from a theme preset.
// Empty strings mean "use preset value".
type SeparatorOverrides struct {
	Left      string `toml:"left,omitempty"`
	Right     string `toml:"right,omitempty"`
	LeftSoft  string `toml:"left_soft,omitempty"`
	RightSoft string `toml:"right_soft,omitempty"`
	LeftHard  string `toml:"left_hard,omitempty"`
	RightHard string `toml:"right_hard,omitempty"`
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
