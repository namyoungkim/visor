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
					{Name: "context_spark"},
					{Name: "compact_eta"},
					{Name: "cache_hit"},
					{Name: "api_latency"},
					{Name: "cost"},
					{Name: "burn_rate"},
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

# === Single-line layout (default) ===
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
  name = "context_spark"
  # [line.widget.extra]
  # width = "8"           # Sparkline width (default: 8)
  # show_label = "true"   # Show "Ctx:" prefix

  [[line.widget]]
  name = "compact_eta"
  # [line.widget.extra]
  # show_when_above = "40"  # Show only when context > 40% (default)
  # show_label = "true"     # Show "ETA:" prefix

  [[line.widget]]
  name = "cache_hit"

  [[line.widget]]
  name = "api_latency"

  [[line.widget]]
  name = "cost"
  # [line.widget.extra]
  # show_label = "true"  # Show "Cost:" prefix

  [[line.widget]]
  name = "burn_rate"
  # [line.widget.extra]
  # show_label = "true"  # Show "Burn:" prefix

  [[line.widget]]
  name = "code_changes"

  [[line.widget]]
  name = "git"

# === Split layout example (left/right aligned) ===
# Uncomment to use split layout instead of single-line:
#
# [[line]]
#   [[line.left]]
#   name = "model"
#
#   [[line.left]]
#   name = "git"
#
#   [[line.left]]
#   name = "context"
#
#   [[line.right]]
#   name = "cost"
#
#   [[line.right]]
#   name = "burn_rate"
#
#   [[line.right]]
#   name = "cache_hit"
`
}
