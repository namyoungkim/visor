package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSave_WritesValidTOML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	cfg := &Config{
		General: GeneralConfig{
			Separator: " | ",
		},
		Lines: []Line{
			{
				Widgets: []WidgetConfig{
					{Name: "model"},
					{Name: "context", Extra: map[string]string{"show_bar": "true"}},
				},
			},
		},
	}

	err := Save(cfg, configPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Expected config file to be created")
	}

	// Verify we can load it back
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Expected to load saved config, got %v", err)
	}

	if len(loaded.Lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(loaded.Lines))
	}

	if len(loaded.Lines[0].Widgets) != 2 {
		t.Errorf("Expected 2 widgets, got %d", len(loaded.Lines[0].Widgets))
	}

	if loaded.Lines[0].Widgets[0].Name != "model" {
		t.Errorf("Expected first widget 'model', got '%s'", loaded.Lines[0].Widgets[0].Name)
	}
}

func TestSave_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", "nested", "config.toml")

	cfg := DefaultConfig()
	err := Save(cfg, configPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Expected config file to be created in nested directory")
	}
}

func TestSave_SplitLayout(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	cfg := &Config{
		General: GeneralConfig{
			Separator: " | ",
		},
		Lines: []Line{
			{
				Left: []WidgetConfig{
					{Name: "model"},
					{Name: "git"},
				},
				Right: []WidgetConfig{
					{Name: "cost"},
				},
			},
		},
	}

	err := Save(cfg, configPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify we can load it back
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Expected to load saved config, got %v", err)
	}

	if len(loaded.Lines[0].Left) != 2 {
		t.Errorf("Expected 2 left widgets, got %d", len(loaded.Lines[0].Left))
	}

	if len(loaded.Lines[0].Right) != 1 {
		t.Errorf("Expected 1 right widget, got %d", len(loaded.Lines[0].Right))
	}
}

func TestDeepCopy_PreservesValues(t *testing.T) {
	original := &Config{
		General: GeneralConfig{
			Separator: " :: ",
		},
		Lines: []Line{
			{
				Widgets: []WidgetConfig{
					{
						Name:   "context",
						Format: "Ctx: {value}",
						Extra:  map[string]string{"show_bar": "true"},
					},
				},
			},
		},
	}

	copy := DeepCopy(original)

	// Verify values are preserved
	if copy.General.Separator != original.General.Separator {
		t.Errorf("Expected separator '%s', got '%s'", original.General.Separator, copy.General.Separator)
	}

	if len(copy.Lines) != len(original.Lines) {
		t.Errorf("Expected %d lines, got %d", len(original.Lines), len(copy.Lines))
	}

	if copy.Lines[0].Widgets[0].Name != "context" {
		t.Errorf("Expected widget name 'context', got '%s'", copy.Lines[0].Widgets[0].Name)
	}

	if copy.Lines[0].Widgets[0].Extra["show_bar"] != "true" {
		t.Errorf("Expected extra show_bar='true', got '%s'", copy.Lines[0].Widgets[0].Extra["show_bar"])
	}
}

func TestDeepCopy_IsIndependent(t *testing.T) {
	original := &Config{
		General: GeneralConfig{
			Separator: " | ",
		},
		Lines: []Line{
			{
				Widgets: []WidgetConfig{
					{Name: "model", Extra: map[string]string{"key": "value"}},
				},
			},
		},
	}

	copy := DeepCopy(original)

	// Modify the copy
	copy.General.Separator = " :: "
	copy.Lines[0].Widgets[0].Name = "cost"
	copy.Lines[0].Widgets[0].Extra["key"] = "changed"

	// Verify original is unchanged
	if original.General.Separator != " | " {
		t.Error("Original separator was modified")
	}

	if original.Lines[0].Widgets[0].Name != "model" {
		t.Error("Original widget name was modified")
	}

	if original.Lines[0].Widgets[0].Extra["key"] != "value" {
		t.Error("Original extra map was modified")
	}
}

func TestDeepCopy_NilConfig(t *testing.T) {
	copy := DeepCopy(nil)
	if copy != nil {
		t.Error("Expected nil for nil input")
	}
}

func TestDeepCopy_SplitLayout(t *testing.T) {
	original := &Config{
		Lines: []Line{
			{
				Left:  []WidgetConfig{{Name: "model"}},
				Right: []WidgetConfig{{Name: "cost"}},
			},
		},
	}

	copy := DeepCopy(original)

	if len(copy.Lines[0].Left) != 1 {
		t.Errorf("Expected 1 left widget, got %d", len(copy.Lines[0].Left))
	}

	if len(copy.Lines[0].Right) != 1 {
		t.Errorf("Expected 1 right widget, got %d", len(copy.Lines[0].Right))
	}

	// Verify independence
	copy.Lines[0].Left[0].Name = "changed"
	if original.Lines[0].Left[0].Name != "model" {
		t.Error("Original left widget was modified")
	}
}
