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
