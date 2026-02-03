package theme

import (
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name     string
		want     string
		powerline bool
	}{
		{"default", "default", false},
		{"powerline", "powerline", true},
		{"gruvbox", "gruvbox", false},
		{"nord", "nord", false},
		{"gruvbox-powerline", "gruvbox-powerline", true},
		{"nord-powerline", "nord-powerline", true},
		{"unknown", "default", false}, // fallback to default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := Get(tt.name)
			if theme.Name != tt.want {
				t.Errorf("Get(%q).Name = %q, want %q", tt.name, theme.Name, tt.want)
			}
			if theme.Powerline != tt.powerline {
				t.Errorf("Get(%q).Powerline = %v, want %v", tt.name, theme.Powerline, tt.powerline)
			}
		})
	}
}

func TestList(t *testing.T) {
	themes := List()

	// Should have at least the core themes
	required := []string{"default", "powerline", "gruvbox", "nord"}
	for _, name := range required {
		found := false
		for _, theme := range themes {
			if theme == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("List() missing required theme %q", name)
		}
	}
}

func TestDefaultSeparators(t *testing.T) {
	sep := DefaultSeparators()
	if sep.Left != " | " {
		t.Errorf("DefaultSeparators().Left = %q, want %q", sep.Left, " | ")
	}
}

func TestPowerlineSeparators(t *testing.T) {
	sep := PowerlineSeparators()
	// Powerline separators use special Unicode characters
	if sep.Left != "" {
		t.Errorf("PowerlineSeparators().Left = %q, want Powerline glyph", sep.Left)
	}
}

func TestColorPalette(t *testing.T) {
	theme := Get("default")

	// Check that essential colors are set
	if theme.Colors.Normal == "" {
		t.Error("ColorPalette.Normal should not be empty")
	}
	if theme.Colors.Warning == "" {
		t.Error("ColorPalette.Warning should not be empty")
	}
	if theme.Colors.Critical == "" {
		t.Error("ColorPalette.Critical should not be empty")
	}
	if theme.Colors.Good == "" {
		t.Error("ColorPalette.Good should not be empty")
	}
}

func TestPowerlineBackgrounds(t *testing.T) {
	theme := Get("powerline")

	if len(theme.Colors.Backgrounds) == 0 {
		t.Error("Powerline theme should have background colors")
	}
}

func TestGruvboxColors(t *testing.T) {
	theme := Get("gruvbox")

	// Gruvbox uses hex colors
	if theme.Colors.Normal != "#ebdbb2" {
		t.Errorf("Gruvbox Normal = %q, want #ebdbb2", theme.Colors.Normal)
	}
}

func TestNordColors(t *testing.T) {
	theme := Get("nord")

	// Nord uses hex colors
	if theme.Colors.Normal != "#eceff4" {
		t.Errorf("Nord Normal = %q, want #eceff4", theme.Colors.Normal)
	}
}
