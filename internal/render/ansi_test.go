package render

import (
	"strings"
	"testing"
)

func TestColorize(t *testing.T) {
	tests := []struct {
		text     string
		fg       string
		contains string
	}{
		{"hello", "red", FgRed},
		{"hello", "green", FgGreen},
		{"hello", "yellow", FgYellow},
		{"hello", "blue", FgBlue},
		{"hello", "cyan", FgCyan},
		{"hello", "magenta", FgMagenta},
		{"hello", "gray", FgBrightBlack},
		{"hello", "grey", FgBrightBlack},
	}

	for _, tt := range tests {
		result := Colorize(tt.text, tt.fg)
		if !strings.Contains(result, tt.contains) {
			t.Errorf("Colorize(%q, %q) should contain %q, got %q",
				tt.text, tt.fg, tt.contains, result)
		}
		if !strings.Contains(result, tt.text) {
			t.Errorf("Colorize(%q, %q) should contain text %q, got %q",
				tt.text, tt.fg, tt.text, result)
		}
		if !strings.HasSuffix(result, Reset) {
			t.Errorf("Colorize(%q, %q) should end with Reset, got %q",
				tt.text, tt.fg, result)
		}
	}
}

func TestColorize_UnknownColor(t *testing.T) {
	result := Colorize("hello", "unknown")
	if result != "hello" {
		t.Errorf("Colorize with unknown color should return plain text, got %q", result)
	}
}

func TestColorize_EmptyColor(t *testing.T) {
	result := Colorize("hello", "")
	if result != "hello" {
		t.Errorf("Colorize with empty color should return plain text, got %q", result)
	}
}

func TestStyle(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		fg       string
		bg       string
		bold     bool
		contains []string
	}{
		{"fg only", "hello", "red", "", false, []string{FgRed}},
		{"bg only", "hello", "", "blue", false, []string{BgBlue}},
		{"bold only", "hello", "", "", true, []string{Bold}},
		{"fg and bg", "hello", "green", "black", false, []string{FgGreen, BgBlack}},
		{"all options", "hello", "yellow", "red", true, []string{Bold, FgYellow, BgRed}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Style(tt.text, tt.fg, tt.bg, tt.bold)
			for _, c := range tt.contains {
				if !strings.Contains(result, c) {
					t.Errorf("Style(%q, %q, %q, %v) should contain %q, got %q",
						tt.text, tt.fg, tt.bg, tt.bold, c, result)
				}
			}
			if !strings.Contains(result, tt.text) {
				t.Errorf("Style should contain text %q, got %q", tt.text, result)
			}
			if !strings.HasSuffix(result, Reset) {
				t.Errorf("Style should end with Reset, got %q", result)
			}
		})
	}
}

func TestStyle_NoOptions(t *testing.T) {
	result := Style("hello", "", "", false)
	if result != "hello" {
		t.Errorf("Style with no options should return plain text, got %q", result)
	}
}

func TestStyle_UnknownColors(t *testing.T) {
	result := Style("hello", "unknown", "unknown", false)
	if result != "hello" {
		t.Errorf("Style with unknown colors should return plain text, got %q", result)
	}
}

func TestRGB(t *testing.T) {
	tests := []struct {
		r, g, b  int
		expected string
	}{
		{255, 0, 0, "\033[38;2;255;0;0m"},
		{0, 255, 0, "\033[38;2;0;255;0m"},
		{0, 0, 255, "\033[38;2;0;0;255m"},
		{128, 128, 128, "\033[38;2;128;128;128m"},
	}

	for _, tt := range tests {
		result := RGB(tt.r, tt.g, tt.b)
		if result != tt.expected {
			t.Errorf("RGB(%d, %d, %d) = %q, expected %q",
				tt.r, tt.g, tt.b, result, tt.expected)
		}
	}
}

func TestColorMap_AllColors(t *testing.T) {
	expectedColors := []string{
		"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white",
		"bright_black", "bright_red", "bright_green", "bright_yellow",
		"bright_blue", "bright_magenta", "bright_cyan", "bright_white",
		"gray", "grey",
	}

	for _, color := range expectedColors {
		if _, ok := ColorMap[color]; !ok {
			t.Errorf("ColorMap should contain %q", color)
		}
	}
}

func TestBgColorMap_AllColors(t *testing.T) {
	expectedColors := []string{
		"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white",
	}

	for _, color := range expectedColors {
		if _, ok := BgColorMap[color]; !ok {
			t.Errorf("BgColorMap should contain %q", color)
		}
	}
}
