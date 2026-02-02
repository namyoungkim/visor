package render

import (
	"os"
	"strconv"
	"unicode/utf8"
)

// TerminalWidth returns the current terminal width.
// Falls back to 80 if unable to determine.
func TerminalWidth() int {
	// Try COLUMNS env var first (faster than ioctl)
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if n, err := strconv.Atoi(cols); err == nil && n > 0 {
			return n
		}
	}
	return 80
}

// Truncate truncates a string to fit within maxWidth.
// Accounts for ANSI escape codes (doesn't count them toward width).
func Truncate(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	visible := 0
	inEscape := false
	result := make([]byte, 0, len(s))

	for i := 0; i < len(s); {
		if s[i] == '\033' {
			// Start of ANSI escape sequence
			inEscape = true
			result = append(result, s[i])
			i++
			continue
		}

		if inEscape {
			result = append(result, s[i])
			// End of escape sequence
			if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') {
				inEscape = false
			}
			i++
			continue
		}

		// Regular character
		r, size := utf8.DecodeRuneInString(s[i:])
		charWidth := runeWidth(r)

		if visible+charWidth > maxWidth {
			// Add ellipsis if there's room
			if visible < maxWidth {
				result = append(result, "..."[:maxWidth-visible]...)
			}
			break
		}

		result = append(result, s[i:i+size]...)
		visible += charWidth
		i += size
	}

	return string(result)
}

// VisibleLength returns the visible length of a string, excluding ANSI codes.
func VisibleLength(s string) int {
	visible := 0
	inEscape := false

	for i := 0; i < len(s); {
		if s[i] == '\033' {
			inEscape = true
			i++
			continue
		}

		if inEscape {
			if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') {
				inEscape = false
			}
			i++
			continue
		}

		r, size := utf8.DecodeRuneInString(s[i:])
		visible += runeWidth(r)
		i += size
	}

	return visible
}

// runeWidth returns the display width of a rune.
// CJK characters are typically double-width.
func runeWidth(r rune) int {
	// CJK ranges (simplified)
	if r >= 0x1100 && r <= 0x115F { // Hangul Jamo
		return 2
	}
	if r >= 0x2E80 && r <= 0x9FFF { // CJK
		return 2
	}
	if r >= 0xAC00 && r <= 0xD7A3 { // Hangul
		return 2
	}
	if r >= 0xF900 && r <= 0xFAFF { // CJK Compatibility
		return 2
	}
	if r >= 0xFF00 && r <= 0xFFEF { // Fullwidth forms
		return 2
	}
	return 1
}
