package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

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
	_, err := toml.DecodeFile(path, &cfg)
	return err
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
