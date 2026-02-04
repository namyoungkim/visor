package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

// hexColorRegex matches hex color codes (#RGB, #RRGGBB, #RRGGBBAA).
var hexColorRegex = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)

// namedColors contains standard terminal color names.
var namedColors = map[string]bool{
	"black":   true,
	"red":     true,
	"green":   true,
	"yellow":  true,
	"blue":    true,
	"magenta": true,
	"cyan":    true,
	"white":   true,
	"gray":    true,
	"grey":    true,
	// Bright variants
	"brightblack":   true,
	"brightred":     true,
	"brightgreen":   true,
	"brightyellow":  true,
	"brightblue":    true,
	"brightmagenta": true,
	"brightcyan":    true,
	"brightwhite":   true,
}

// validateColor checks if a color string is valid.
func validateColor(color string) bool {
	if color == "" {
		return true
	}
	if strings.HasPrefix(color, "#") {
		return hexColorRegex.MatchString(color)
	}
	return namedColors[strings.ToLower(color)]
}

// DefaultConfigPath returns the default config file path.
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "visor", "config.toml")
}

// Load loads configuration from the given path.
// Falls back to default config if the file doesn't exist.
func Load(path string) (*Config, error) {
	if path == "" {
		path = DefaultConfigPath()
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}

	if len(cfg.Lines) == 0 {
		return DefaultConfig(), nil
	}

	// Apply default separator if not set
	if cfg.General.Separator == "" {
		cfg.General.Separator = DefaultSeparator
	}

	return &cfg, nil
}

// Validate checks if the configuration is valid.
func Validate(path string) error {
	if path == "" {
		path = DefaultConfigPath()
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // No config file is valid (uses defaults)
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return err
	}

	// Validate theme colors if specified
	if cfg.Theme.Colors != nil {
		if err := validateColors(cfg.Theme.Colors); err != nil {
			return err
		}
	}

	return nil
}

// validateColors validates all color overrides.
func validateColors(colors *ColorOverrides) error {
	colorFields := map[string]string{
		"normal":    colors.Normal,
		"warning":   colors.Warning,
		"critical":  colors.Critical,
		"good":      colors.Good,
		"primary":   colors.Primary,
		"secondary": colors.Secondary,
		"muted":     colors.Muted,
	}

	for name, color := range colorFields {
		if !validateColor(color) {
			return fmt.Errorf("invalid color for %s: %q", name, color)
		}
	}

	for i, bg := range colors.Backgrounds {
		if !validateColor(bg) {
			return fmt.Errorf("invalid background color at index %d: %q", i, bg)
		}
	}

	return nil
}

// Init creates a default configuration file at the given path.
func Init(path string) error {
	if path == "" {
		path = DefaultConfigPath()
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(DefaultConfigTOML()), 0644)
}
