package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_DefaultConfig(t *testing.T) {
	// When config file doesn't exist, should return default config
	cfg, err := Load("/nonexistent/path/config.toml")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(cfg.Lines) == 0 {
		t.Error("Expected default config to have lines")
	}

	// Check default widgets are present (11 widgets including burn_rate, compact_eta, context_spark, block_timer)
	if len(cfg.Lines[0].Widgets) != 11 {
		t.Errorf("Expected 11 default widgets, got %d", len(cfg.Lines[0].Widgets))
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	// Create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	content := `
[[line]]
  [[line.widget]]
  name = "model"

  [[line.widget]]
  name = "cost"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(cfg.Lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(cfg.Lines))
	}

	if len(cfg.Lines[0].Widgets) != 2 {
		t.Errorf("Expected 2 widgets, got %d", len(cfg.Lines[0].Widgets))
	}

	if cfg.Lines[0].Widgets[0].Name != "model" {
		t.Errorf("Expected first widget 'model', got '%s'", cfg.Lines[0].Widgets[0].Name)
	}
}

func TestLoad_InvalidTOML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	content := `invalid toml [[[`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Expected error for invalid TOML")
	}
}

func TestLoad_EmptyConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	content := ``
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Empty config should fall back to defaults
	if len(cfg.Lines) == 0 {
		t.Error("Expected default config for empty file")
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	content := `[[line]]
  [[line.widget]]
  name = "model"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}

	err := Validate(configPath)
	if err != nil {
		t.Errorf("Expected valid config, got error: %v", err)
	}
}

func TestValidate_InvalidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	content := `invalid [[[`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}

	err := Validate(configPath)
	if err == nil {
		t.Error("Expected error for invalid config")
	}
}

func TestValidate_NonexistentFile(t *testing.T) {
	err := Validate("/nonexistent/path/config.toml")
	if err != nil {
		t.Error("Expected no error for nonexistent file (uses defaults)")
	}
}

func TestInit_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", "config.toml")

	err := Init(configPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Expected config file to be created")
	}

	// Check content is valid TOML
	err = Validate(configPath)
	if err != nil {
		t.Errorf("Expected valid TOML, got error: %v", err)
	}
}

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()
	if path == "" {
		t.Error("Expected non-empty default config path")
	}

	if !filepath.IsAbs(path) {
		t.Error("Expected absolute path")
	}
}

func TestLoad_DefaultSeparator(t *testing.T) {
	// Default config should have separator set
	cfg, err := Load("/nonexistent/path/config.toml")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.General.Separator != DefaultSeparator {
		t.Errorf("Expected default separator '%s', got '%s'", DefaultSeparator, cfg.General.Separator)
	}
}

func TestLoad_MissingSeparator(t *testing.T) {
	// Config without [general] section should get default separator
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	content := `[[line]]
  [[line.widget]]
  name = "model"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.General.Separator != DefaultSeparator {
		t.Errorf("Expected default separator '%s', got '%s'", DefaultSeparator, cfg.General.Separator)
	}
}

func TestLoad_CustomSeparator(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	content := `[general]
separator = " :: "

[[line]]
  [[line.widget]]
  name = "model"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.General.Separator != " :: " {
		t.Errorf("Expected separator ' :: ', got '%s'", cfg.General.Separator)
	}
}

func TestLoad_CustomTheme(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	content := `[theme]
name = "gruvbox"
powerline = true

[theme.colors]
warning = "#ff00ff"
critical = "#ff0000"
backgrounds = ["#111111", "#222222"]

[theme.separators]
left = " :: "
right = " :: "

[[line]]
  [[line.widget]]
  name = "model"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify theme settings
	if cfg.Theme.Name != "gruvbox" {
		t.Errorf("Expected theme name 'gruvbox', got '%s'", cfg.Theme.Name)
	}

	if !cfg.Theme.Powerline {
		t.Error("Expected powerline to be true")
	}

	// Verify color overrides
	if cfg.Theme.Colors == nil {
		t.Fatal("Expected colors to be parsed")
	}

	if cfg.Theme.Colors.Warning != "#ff00ff" {
		t.Errorf("Expected warning '#ff00ff', got '%s'", cfg.Theme.Colors.Warning)
	}

	if cfg.Theme.Colors.Critical != "#ff0000" {
		t.Errorf("Expected critical '#ff0000', got '%s'", cfg.Theme.Colors.Critical)
	}

	if len(cfg.Theme.Colors.Backgrounds) != 2 {
		t.Fatalf("Expected 2 backgrounds, got %d", len(cfg.Theme.Colors.Backgrounds))
	}

	if cfg.Theme.Colors.Backgrounds[0] != "#111111" {
		t.Errorf("Expected background '#111111', got '%s'", cfg.Theme.Colors.Backgrounds[0])
	}

	// Verify separator overrides
	if cfg.Theme.Separators == nil {
		t.Fatal("Expected separators to be parsed")
	}

	if cfg.Theme.Separators.Left != " :: " {
		t.Errorf("Expected left separator ' :: ', got '%s'", cfg.Theme.Separators.Left)
	}
}

func TestLoad_ThemeWithoutOverrides(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	content := `[theme]
name = "nord"

[[line]]
  [[line.widget]]
  name = "model"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Theme.Name != "nord" {
		t.Errorf("Expected theme name 'nord', got '%s'", cfg.Theme.Name)
	}

	// Colors and Separators should be nil when not specified
	if cfg.Theme.Colors != nil {
		t.Error("Expected colors to be nil when not specified")
	}

	if cfg.Theme.Separators != nil {
		t.Error("Expected separators to be nil when not specified")
	}
}

func TestValidate_InvalidColor(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	content := `[theme]
name = "gruvbox"

[theme.colors]
warning = "not-a-color"

[[line]]
  [[line.widget]]
  name = "model"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}

	err := Validate(configPath)
	if err == nil {
		t.Error("Expected error for invalid color")
	}
}

func TestValidate_InvalidBackgroundColor(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	content := `[theme]
name = "gruvbox"

[theme.colors]
backgrounds = ["#111111", "invalid"]

[[line]]
  [[line.widget]]
  name = "model"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}

	err := Validate(configPath)
	if err == nil {
		t.Error("Expected error for invalid background color")
	}
}

func TestValidate_ValidColors(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	content := `[theme]
name = "gruvbox"

[theme.colors]
warning = "#ff00ff"
critical = "red"
normal = "#abc"
backgrounds = ["#111111", "#222222"]

[[line]]
  [[line.widget]]
  name = "model"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}

	err := Validate(configPath)
	if err != nil {
		t.Errorf("Expected no error for valid colors, got: %v", err)
	}
}
