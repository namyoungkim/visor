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
					{Name: "block_timer"},
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
  # show_label = "false"         # Hide "Ctx:" prefix
  # show_bar = "false"           # Hide progress bar
  # bar_width = "10"             # Progress bar width (default: 10)
  # warn_threshold = "60"        # Warning color threshold % (default: 60)
  # critical_threshold = "80"    # Critical color threshold % (default: 80)

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
  name = "block_timer"
  # [line.widget.extra]
  # show_label = "true"          # Show "Block:" prefix (default: true)
  # warn_threshold = "80"        # Warning at 80% elapsed (default: 80)
  # critical_threshold = "95"    # Critical at 95% elapsed (default: 95)

  [[line.widget]]
  name = "cache_hit"
  # [line.widget.extra]
  # good_threshold = "80"   # Good/green threshold % (default: 80)
  # warn_threshold = "50"   # Warning threshold % (default: 50)

  [[line.widget]]
  name = "api_latency"
  # [line.widget.extra]
  # warn_threshold = "2000"      # Warning threshold ms (default: 2000)
  # critical_threshold = "5000"  # Critical threshold ms (default: 5000)

  [[line.widget]]
  name = "cost"
  # [line.widget.extra]
  # show_label = "true"          # Show "Cost:" prefix
  # warn_threshold = "0.5"       # Warning threshold USD (default: 0.5)
  # critical_threshold = "1.0"   # Critical threshold USD (default: 1.0)

  [[line.widget]]
  name = "burn_rate"
  # [line.widget.extra]
  # show_label = "true"          # Show "Burn:" prefix
  # warn_threshold = "10"        # Warning threshold cents/min (default: 10)
  # critical_threshold = "25"    # Critical threshold cents/min (default: 25)

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
