package theme

import "sort"

// Theme represents a complete theme configuration.
type Theme struct {
	Name       string
	Colors     ColorPalette
	Separators SeparatorSet
	Powerline  bool
}

// ColorPalette defines colors for various states and elements.
type ColorPalette struct {
	// Status colors
	Normal   string // Default/neutral state
	Warning  string // Warning state (yellow-ish)
	Critical string // Critical/error state (red-ish)
	Good     string // Success/good state (green-ish)

	// UI element colors
	Primary   string // Primary accent color
	Secondary string // Secondary accent color
	Muted     string // Muted/dimmed text

	// Powerline-specific background colors (for segment backgrounds)
	Backgrounds []string
}

// SeparatorSet defines the characters used between widgets.
type SeparatorSet struct {
	Left       string // Regular left separator
	Right      string // Regular right separator
	LeftSoft   string // Soft separator (same background)
	RightSoft  string // Soft separator (same background)
	LeftHard   string // Hard separator (different background) - Powerline
	RightHard  string // Hard separator (different background) - Powerline
}

// DefaultSeparators returns standard ASCII separators.
func DefaultSeparators() SeparatorSet {
	return SeparatorSet{
		Left:      " | ",
		Right:     " | ",
		LeftSoft:  " | ",
		RightSoft: " | ",
		LeftHard:  " | ",
		RightHard: " | ",
	}
}

// PowerlineSeparators returns Powerline font separators.
func PowerlineSeparators() SeparatorSet {
	return SeparatorSet{
		Left:      "",         // U+E0B0
		Right:     "",         // U+E0B2
		LeftSoft:  "",         // U+E0B1
		RightSoft: "",         // U+E0B3
		LeftHard:  "",         // U+E0B0
		RightHard: "",         // U+E0B2
	}
}

// Get returns a theme by name.
// Returns the default theme if name is not found.
func Get(name string) *Theme {
	if t, ok := presets[name]; ok {
		return t
	}
	return presets["default"]
}

// List returns all available theme names in sorted order.
func List() []string {
	names := make([]string, 0, len(presets))
	for name := range presets {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
