package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Save writes the configuration to the given path
func Save(cfg *Config, path string) error {
	if path == "" {
		path = DefaultConfigPath()
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Encode to TOML
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	encoder.Indent = "  "
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// DeepCopy creates a deep copy of the configuration
func DeepCopy(cfg *Config) *Config {
	if cfg == nil {
		return nil
	}

	newCfg := &Config{
		General: GeneralConfig{
			Separator: cfg.General.Separator,
		},
		Theme: ThemeConfig{
			Name:      cfg.Theme.Name,
			Powerline: cfg.Theme.Powerline,
		},
		Usage: UsageConfig{
			Enabled:     cfg.Usage.Enabled,
			Provider:    cfg.Usage.Provider,
			ProjectsDir: cfg.Usage.ProjectsDir,
		},
		Lines: make([]Line, len(cfg.Lines)),
	}

	// Deep copy color overrides
	if cfg.Theme.Colors != nil {
		newCfg.Theme.Colors = &ColorOverrides{
			Normal:    cfg.Theme.Colors.Normal,
			Warning:   cfg.Theme.Colors.Warning,
			Critical:  cfg.Theme.Colors.Critical,
			Good:      cfg.Theme.Colors.Good,
			Primary:   cfg.Theme.Colors.Primary,
			Secondary: cfg.Theme.Colors.Secondary,
			Muted:     cfg.Theme.Colors.Muted,
		}
		if len(cfg.Theme.Colors.Backgrounds) > 0 {
			newCfg.Theme.Colors.Backgrounds = make([]string, len(cfg.Theme.Colors.Backgrounds))
			copy(newCfg.Theme.Colors.Backgrounds, cfg.Theme.Colors.Backgrounds)
		}
	}

	// Deep copy separator overrides
	if cfg.Theme.Separators != nil {
		newCfg.Theme.Separators = &SeparatorOverrides{
			Left:      cfg.Theme.Separators.Left,
			Right:     cfg.Theme.Separators.Right,
			LeftSoft:  cfg.Theme.Separators.LeftSoft,
			RightSoft: cfg.Theme.Separators.RightSoft,
			LeftHard:  cfg.Theme.Separators.LeftHard,
			RightHard: cfg.Theme.Separators.RightHard,
		}
	}

	for i, line := range cfg.Lines {
		newCfg.Lines[i] = deepCopyLine(line)
	}

	return newCfg
}

func deepCopyLine(line Line) Line {
	newLine := Line{}

	if line.Widgets != nil {
		newLine.Widgets = make([]WidgetConfig, len(line.Widgets))
		for i, w := range line.Widgets {
			newLine.Widgets[i] = deepCopyWidgetConfig(w)
		}
	}

	if line.Left != nil {
		newLine.Left = make([]WidgetConfig, len(line.Left))
		for i, w := range line.Left {
			newLine.Left[i] = deepCopyWidgetConfig(w)
		}
	}

	if line.Right != nil {
		newLine.Right = make([]WidgetConfig, len(line.Right))
		for i, w := range line.Right {
			newLine.Right[i] = deepCopyWidgetConfig(w)
		}
	}

	return newLine
}

func deepCopyWidgetConfig(w WidgetConfig) WidgetConfig {
	newW := WidgetConfig{
		Name:   w.Name,
		Format: w.Format,
		Style: StyleConfig{
			Fg:   w.Style.Fg,
			Bg:   w.Style.Bg,
			Bold: w.Style.Bold,
		},
	}

	if w.Extra != nil {
		newW.Extra = make(map[string]string, len(w.Extra))
		for k, v := range w.Extra {
			newW.Extra[k] = v
		}
	}

	return newW
}
