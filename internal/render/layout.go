package render

import (
	"strings"
)

// Separator is the default separator between widgets.
const Separator = " "

// Layout combines rendered widget strings into final output.
func Layout(widgets []string) string {
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

	line := strings.Join(nonEmpty, Separator)

	// Truncate to terminal width
	width := TerminalWidth()
	return Truncate(line, width)
}

// MultiLine renders multiple lines of widgets.
func MultiLine(lines [][]string) string {
	var result []string

	for _, lineWidgets := range lines {
		line := Layout(lineWidgets)
		if line != "" {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
