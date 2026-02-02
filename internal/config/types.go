package config

// Config represents the visor configuration.
type Config struct {
	Lines []Line `toml:"line"`
}

// Line represents a single line in the statusline.
type Line struct {
	Widgets []WidgetConfig `toml:"widget"`
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
