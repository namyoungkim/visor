package theme

import (
	"regexp"
	"strings"
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

// ValidateColor checks if a color string is valid.
// Valid formats:
//   - Hex colors: #RGB, #RRGGBB, #RRGGBBAA
//   - Named colors: black, red, green, yellow, blue, magenta, cyan, white, gray
//   - Empty string (treated as "use preset value")
func ValidateColor(color string) bool {
	// Empty string is valid (means "use preset value")
	if color == "" {
		return true
	}

	// Check hex format
	if strings.HasPrefix(color, "#") {
		return hexColorRegex.MatchString(color)
	}

	// Check named colors (case-insensitive)
	return namedColors[strings.ToLower(color)]
}
