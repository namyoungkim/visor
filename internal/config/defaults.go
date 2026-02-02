package config

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
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

[[line]]
  [[line.widget]]
  name = "model"

  [[line.widget]]
  name = "context"

  [[line.widget]]
  name = "cache_hit"

  [[line.widget]]
  name = "api_latency"

  [[line.widget]]
  name = "cost"

  [[line.widget]]
  name = "code_changes"

  [[line.widget]]
  name = "git"
`
}
