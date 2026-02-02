package render

import (
	"strings"
)

// Layout combines rendered widget strings into final output.
func Layout(widgets []string, separator string) string {
	// Filter out empty widgets
	var nonEmpty []string
	for _, w := range widgets {
		if w != "" {
			nonEmpty = append(nonEmpty, w)
		}
	}

	if len(nonEmpty) == 0 {
		return ""
	}

	line := strings.Join(nonEmpty, separator)

	// Truncate to terminal width
	width := TerminalWidth()
	return Truncate(line, width)
}

// SplitLayout renders left and right aligned widgets on a single line.
// Left widgets are joined normally, right widgets are right-aligned.
// Example: "model | git                      cost | cache"
func SplitLayout(left, right []string, separator string) string {
	width := TerminalWidth()

	leftStr := joinNonEmpty(left, separator)
	rightStr := joinNonEmpty(right, separator)

	// If only one side has content, use regular layout
	if leftStr == "" && rightStr == "" {
		return ""
	}
	if rightStr == "" {
		return Truncate(leftStr, width)
	}
	if leftStr == "" {
		return Truncate(rightStr, width)
	}

	// Calculate visual lengths (without ANSI codes)
	leftLen := VisibleLength(leftStr)
	rightLen := VisibleLength(rightStr)

	// Calculate padding needed between left and right
	minGap := 2 // Minimum space between sides
	totalContentLen := leftLen + rightLen + minGap

	if totalContentLen > width {
		// Not enough space - fall back to regular layout
		combined := leftStr + separator + rightStr
		return Truncate(combined, width)
	}

	// Calculate padding to push right side to the right edge
	padding := width - leftLen - rightLen
	if padding < minGap {
		padding = minGap
	}

	return leftStr + strings.Repeat(" ", padding) + rightStr
}

// joinNonEmpty joins non-empty strings with a separator.
func joinNonEmpty(items []string, separator string) string {
	var nonEmpty []string
	for _, s := range items {
		if s != "" {
			nonEmpty = append(nonEmpty, s)
		}
	}
	return strings.Join(nonEmpty, separator)
}

// JoinLines joins multiple lines with newline characters.
func JoinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

// MultiLine renders multiple lines of widgets.
func MultiLine(lines [][]string, separator string) string {
	var result []string

	for _, lineWidgets := range lines {
		line := Layout(lineWidgets, separator)
		if line != "" {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
