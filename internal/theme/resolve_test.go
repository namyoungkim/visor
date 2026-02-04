package theme

import (
	"testing"

	"github.com/namyoungkim/visor/internal/config"
)

func TestResolve_NilConfig(t *testing.T) {
	resolved := Resolve(nil)
	if resolved == nil {
		t.Fatal("expected non-nil theme")
	}
	if resolved.Name != "default" {
		t.Errorf("expected default theme, got %s", resolved.Name)
	}
}

func TestResolve_EmptyName(t *testing.T) {
	cfg := &config.ThemeConfig{Name: ""}
	resolved := Resolve(cfg)
	if resolved.Name != "default" {
		t.Errorf("expected default theme, got %s", resolved.Name)
	}
}

func TestResolve_PresetOnly(t *testing.T) {
	cfg := &config.ThemeConfig{Name: "gruvbox"}
	resolved := Resolve(cfg)

	if resolved.Name != "gruvbox" {
		t.Errorf("expected gruvbox theme, got %s", resolved.Name)
	}
	if resolved.Colors.Warning != "#fabd2f" {
		t.Errorf("expected gruvbox warning color, got %s", resolved.Colors.Warning)
	}
}

func TestResolve_ColorOverrides(t *testing.T) {
	cfg := &config.ThemeConfig{
		Name: "gruvbox",
		Colors: &config.ColorOverrides{
			Warning:  "#ff00ff",
			Critical: "#ff0000",
		},
	}
	resolved := Resolve(cfg)

	// Overridden colors
	if resolved.Colors.Warning != "#ff00ff" {
		t.Errorf("expected warning override #ff00ff, got %s", resolved.Colors.Warning)
	}
	if resolved.Colors.Critical != "#ff0000" {
		t.Errorf("expected critical override #ff0000, got %s", resolved.Colors.Critical)
	}

	// Non-overridden colors should remain from preset
	if resolved.Colors.Good != "#b8bb26" {
		t.Errorf("expected gruvbox good color #b8bb26, got %s", resolved.Colors.Good)
	}
	if resolved.Colors.Normal != "#ebdbb2" {
		t.Errorf("expected gruvbox normal color #ebdbb2, got %s", resolved.Colors.Normal)
	}
}

func TestResolve_BackgroundOverride(t *testing.T) {
	cfg := &config.ThemeConfig{
		Name: "gruvbox",
		Colors: &config.ColorOverrides{
			Backgrounds: []string{"#111111", "#222222"},
		},
	}
	resolved := Resolve(cfg)

	if len(resolved.Colors.Backgrounds) != 2 {
		t.Fatalf("expected 2 backgrounds, got %d", len(resolved.Colors.Backgrounds))
	}
	if resolved.Colors.Backgrounds[0] != "#111111" {
		t.Errorf("expected background #111111, got %s", resolved.Colors.Backgrounds[0])
	}
}

func TestResolve_SeparatorOverrides(t *testing.T) {
	cfg := &config.ThemeConfig{
		Name: "default",
		Separators: &config.SeparatorOverrides{
			Left:  " :: ",
			Right: " :: ",
		},
	}
	resolved := Resolve(cfg)

	if resolved.Separators.Left != " :: " {
		t.Errorf("expected left separator ' :: ', got %s", resolved.Separators.Left)
	}
	if resolved.Separators.Right != " :: " {
		t.Errorf("expected right separator ' :: ', got %s", resolved.Separators.Right)
	}

	// Non-overridden separators should remain from preset
	if resolved.Separators.LeftSoft != " | " {
		t.Errorf("expected default left soft separator, got %s", resolved.Separators.LeftSoft)
	}
}

func TestResolve_PowerlineOverride(t *testing.T) {
	cfg := &config.ThemeConfig{
		Name:      "gruvbox",
		Powerline: true,
	}
	resolved := Resolve(cfg)

	if !resolved.Powerline {
		t.Error("expected powerline to be enabled")
	}
	// When powerline is enabled, separators should be powerline style
	if resolved.Separators.LeftHard != "" {
		t.Errorf("expected powerline left hard separator, got %s", resolved.Separators.LeftHard)
	}
}

func TestResolve_DoesNotModifyPreset(t *testing.T) {
	// Get original preset
	original := Get("gruvbox")
	originalWarning := original.Colors.Warning

	// Resolve with override
	cfg := &config.ThemeConfig{
		Name: "gruvbox",
		Colors: &config.ColorOverrides{
			Warning: "#ffffff",
		},
	}
	Resolve(cfg)

	// Original preset should be unchanged
	if original.Colors.Warning != originalWarning {
		t.Error("preset was modified by Resolve")
	}
}

func TestClone(t *testing.T) {
	original := Get("powerline")
	cloned := clone(original)

	// Modify clone
	cloned.Colors.Warning = "#000000"
	cloned.Colors.Backgrounds[0] = "#ffffff"

	// Original should be unchanged
	if original.Colors.Warning == "#000000" {
		t.Error("clone modified original colors")
	}
	if original.Colors.Backgrounds[0] == "#ffffff" {
		t.Error("clone modified original backgrounds")
	}
}

func TestClone_Nil(t *testing.T) {
	cloned := clone(nil)
	if cloned != nil {
		t.Error("clone of nil should return nil")
	}
}

func TestResolve_InvalidPreset(t *testing.T) {
	cfg := &config.ThemeConfig{Name: "nonexistent"}
	resolved := Resolve(cfg)

	// Invalid preset should fallback to default theme
	if resolved == nil {
		t.Fatal("expected non-nil theme")
	}
	if resolved.Name != "default" {
		t.Errorf("expected fallback to default, got %s", resolved.Name)
	}
}
