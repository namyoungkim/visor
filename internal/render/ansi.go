package render

import "fmt"

// ANSI color codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Underline = "\033[4m"
)

// Foreground colors
const (
	FgBlack   = "\033[30m"
	FgRed     = "\033[31m"
	FgGreen   = "\033[32m"
	FgYellow  = "\033[33m"
	FgBlue    = "\033[34m"
	FgMagenta = "\033[35m"
	FgCyan    = "\033[36m"
	FgWhite   = "\033[37m"
	FgDefault = "\033[39m"
)

// Bright foreground colors
const (
	FgBrightBlack   = "\033[90m"
	FgBrightRed     = "\033[91m"
	FgBrightGreen   = "\033[92m"
	FgBrightYellow  = "\033[93m"
	FgBrightBlue    = "\033[94m"
	FgBrightMagenta = "\033[95m"
	FgBrightCyan    = "\033[96m"
	FgBrightWhite   = "\033[97m"
)

// Background colors
const (
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
	BgDefault = "\033[49m"
)

// ColorMap maps color names to ANSI codes
var ColorMap = map[string]string{
	"black":          FgBlack,
	"red":            FgRed,
	"green":          FgGreen,
	"yellow":         FgYellow,
	"blue":           FgBlue,
	"magenta":        FgMagenta,
	"cyan":           FgCyan,
	"white":          FgWhite,
	"bright_black":   FgBrightBlack,
	"bright_red":     FgBrightRed,
	"bright_green":   FgBrightGreen,
	"bright_yellow":  FgBrightYellow,
	"bright_blue":    FgBrightBlue,
	"bright_magenta": FgBrightMagenta,
	"bright_cyan":    FgBrightCyan,
	"bright_white":   FgBrightWhite,
	"gray":           FgBrightBlack,
	"grey":           FgBrightBlack,
}

// BgColorMap maps color names to background ANSI codes
var BgColorMap = map[string]string{
	"black":   BgBlack,
	"red":     BgRed,
	"green":   BgGreen,
	"yellow":  BgYellow,
	"blue":    BgBlue,
	"magenta": BgMagenta,
	"cyan":    BgCyan,
	"white":   BgWhite,
}

// Colorize applies ANSI color to text.
func Colorize(text, fg string) string {
	if code, ok := ColorMap[fg]; ok {
		return code + text + Reset
	}
	return text
}

// Style applies multiple style options to text.
func Style(text string, fg, bg string, bold bool) string {
	var codes string

	if bold {
		codes += Bold
	}
	if code, ok := BgColorMap[bg]; ok {
		codes += code
	}
	if code, ok := ColorMap[fg]; ok {
		codes += code
	}

	if codes == "" {
		return text
	}
	return codes + text + Reset
}

// RGB returns an ANSI 256-color or true color code.
func RGB(r, g, b int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

// RGBBg returns an ANSI true color background code.
func RGBBg(r, g, b int) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
}

// HexToRGB converts a hex color string (#RRGGBB) to RGB values.
func HexToRGB(hex string) (r, g, b int) {
	if len(hex) == 0 {
		return 0, 0, 0
	}
	if hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) != 6 {
		return 0, 0, 0
	}
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return r, g, b
}

// ColorizeHex applies a hex color to text.
func ColorizeHex(text, hex string) string {
	if hex == "" {
		return text
	}
	r, g, b := HexToRGB(hex)
	return RGB(r, g, b) + text + Reset
}

// StyleHex applies hex colors and styling to text.
func StyleHex(text, fgHex, bgHex string, bold bool) string {
	var codes string

	if bold {
		codes += Bold
	}
	if bgHex != "" {
		r, g, b := HexToRGB(bgHex)
		codes += RGBBg(r, g, b)
	}
	if fgHex != "" {
		r, g, b := HexToRGB(fgHex)
		codes += RGB(r, g, b)
	}

	if codes == "" {
		return text
	}
	return codes + text + Reset
}

// ResolveColor returns the ANSI code for a color (name or hex).
func ResolveColor(color string) string {
	if color == "" {
		return ""
	}

	// Check if it's a hex color
	if color[0] == '#' {
		r, g, b := HexToRGB(color)
		return RGB(r, g, b)
	}

	// Check named colors
	if code, ok := ColorMap[color]; ok {
		return code
	}

	return ""
}

// ResolveBgColor returns the ANSI background code for a color (name or hex).
func ResolveBgColor(color string) string {
	if color == "" {
		return ""
	}

	// Check if it's a hex color
	if color[0] == '#' {
		r, g, b := HexToRGB(color)
		return RGBBg(r, g, b)
	}

	// Check named colors
	if code, ok := BgColorMap[color]; ok {
		return code
	}

	return ""
}
