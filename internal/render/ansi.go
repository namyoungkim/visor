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
