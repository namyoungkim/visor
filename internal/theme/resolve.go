package theme

import "github.com/namyoungkim/visor/internal/config"

// Resolve returns a theme by resolving the preset and applying any overrides.
// If cfg is nil or has an empty name, returns the default theme.
func Resolve(cfg *config.ThemeConfig) *Theme {
	// Get base theme from preset
	var base *Theme
	if cfg == nil || cfg.Name == "" {
		base = Get("default")
	} else {
		base = Get(cfg.Name)
	}

	// Clone the base theme to avoid modifying presets
	resolved := clone(base)

	// Apply powerline override from config if explicitly set
	if cfg != nil && cfg.Powerline && !resolved.Powerline {
		resolved.Powerline = true
		resolved.Separators = PowerlineSeparators()
	}

	// Apply color overrides
	if cfg != nil && cfg.Colors != nil {
		applyColorOverrides(&resolved.Colors, cfg.Colors)
	}

	// Apply separator overrides
	if cfg != nil && cfg.Separators != nil {
		applySeparatorOverrides(&resolved.Separators, cfg.Separators)
	}

	return resolved
}

// clone creates a deep copy of a theme.
func clone(t *Theme) *Theme {
	if t == nil {
		return nil
	}

	// Copy backgrounds slice
	backgrounds := make([]string, len(t.Colors.Backgrounds))
	copy(backgrounds, t.Colors.Backgrounds)

	return &Theme{
		Name:      t.Name,
		Powerline: t.Powerline,
		Colors: ColorPalette{
			Normal:      t.Colors.Normal,
			Warning:     t.Colors.Warning,
			Critical:    t.Colors.Critical,
			Good:        t.Colors.Good,
			Primary:     t.Colors.Primary,
			Secondary:   t.Colors.Secondary,
			Muted:       t.Colors.Muted,
			Backgrounds: backgrounds,
		},
		Separators: SeparatorSet{
			Left:      t.Separators.Left,
			Right:     t.Separators.Right,
			LeftSoft:  t.Separators.LeftSoft,
			RightSoft: t.Separators.RightSoft,
			LeftHard:  t.Separators.LeftHard,
			RightHard: t.Separators.RightHard,
		},
	}
}

// applyColorOverrides applies non-empty color overrides to the palette.
func applyColorOverrides(palette *ColorPalette, overrides *config.ColorOverrides) {
	if overrides.Normal != "" {
		palette.Normal = overrides.Normal
	}
	if overrides.Warning != "" {
		palette.Warning = overrides.Warning
	}
	if overrides.Critical != "" {
		palette.Critical = overrides.Critical
	}
	if overrides.Good != "" {
		palette.Good = overrides.Good
	}
	if overrides.Primary != "" {
		palette.Primary = overrides.Primary
	}
	if overrides.Secondary != "" {
		palette.Secondary = overrides.Secondary
	}
	if overrides.Muted != "" {
		palette.Muted = overrides.Muted
	}
	if len(overrides.Backgrounds) > 0 {
		palette.Backgrounds = make([]string, len(overrides.Backgrounds))
		copy(palette.Backgrounds, overrides.Backgrounds)
	}
}

// applySeparatorOverrides applies non-empty separator overrides.
func applySeparatorOverrides(seps *SeparatorSet, overrides *config.SeparatorOverrides) {
	if overrides.Left != "" {
		seps.Left = overrides.Left
	}
	if overrides.Right != "" {
		seps.Right = overrides.Right
	}
	if overrides.LeftSoft != "" {
		seps.LeftSoft = overrides.LeftSoft
	}
	if overrides.RightSoft != "" {
		seps.RightSoft = overrides.RightSoft
	}
	if overrides.LeftHard != "" {
		seps.LeftHard = overrides.LeftHard
	}
	if overrides.RightHard != "" {
		seps.RightHard = overrides.RightHard
	}
}
