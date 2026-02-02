package config

// Config represents the visor configuration.
type Config struct {
	General GeneralConfig `toml:"general"`
	Lines   []Line        `toml:"line"`
}

// GeneralConfig contains global settings.
type GeneralConfig struct {
	Separator string `toml:"separator"`
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
