package config

// DefaultSeparator is the default separator between widgets.
const DefaultSeparator = " | "

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		General: GeneralConfig{
			Separator: DefaultSeparator,
		},
		Lines: []Line{
			{
				Widgets: []WidgetConfig{
					{Name: "model"},
					{Name: "context"},
					{Name: "cache_hit"},
					{Name: "api_latency"},
					{Name: "cost"},
					{Name: "code_changes"},
					{Name: "git"},
				},
			},
		},
	}
}

// DefaultConfigTOML returns the default configuration as a TOML string.
func DefaultConfigTOML() string {
	return `# visor configuration
# Place at ~/.config/visor/config.toml

[general]
separator = " | "  # Widget separator (default: " | ")

[[line]]
  [[line.widget]]
  name = "model"

  [[line.widget]]
  name = "context"
  # format = "Context: {value}"  # Custom format (optional)
  # [line.widget.extra]
  # show_label = "false"  # Hide "Ctx:" prefix
  # show_bar = "false"    # Hide progress bar
  # bar_width = "10"      # Progress bar width (default: 10)

  [[line.widget]]
  name = "cache_hit"

  [[line.widget]]
  name = "api_latency"

  [[line.widget]]
  name = "cost"
  # [line.widget.extra]
  # show_label = "true"  # Show "Cost:" prefix

  [[line.widget]]
  name = "code_changes"

  [[line.widget]]
  name = "git"
`
}
