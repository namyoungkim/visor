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

// PowerlineSegment represents a styled powerline segment.
type PowerlineSegment struct {
	Text string // Widget text (without ANSI codes if possible)
	Bg   string // Background color (hex or name)
	Fg   string // Foreground color (hex or name)
}

// PowerlineLayout renders widgets in powerline style with arrows and backgrounds.
// Each widget gets a colored background, with arrow separators between them.
func PowerlineLayout(segments []PowerlineSegment, separator string) string {
	if separator == "" {
		separator = "" // Default powerline separator (U+E0B0)
	}

	// Filter out segments with empty text
	var nonEmpty []PowerlineSegment
	for _, s := range segments {
		if s.Text != "" {
			nonEmpty = append(nonEmpty, s)
		}
	}

	if len(nonEmpty) == 0 {
		return ""
	}

	var result strings.Builder

	for i, seg := range nonEmpty {
		// Apply background and foreground colors
		bg := ResolveBgColor(seg.Bg)
		fg := ResolveColor(seg.Fg)

		// Write segment text with styling
		result.WriteString(bg)
		result.WriteString(fg)
		result.WriteString(" ")
		result.WriteString(seg.Text)
		result.WriteString(" ")

		// Write separator arrow
		if i < len(nonEmpty)-1 {
			nextBg := ResolveBgColor(nonEmpty[i+1].Bg)
			// Arrow: current bg as foreground, next bg as background
			arrowFg := seg.Bg
			if arrowFg == "" {
				arrowFg = "black"
			}
			result.WriteString(Reset)
			result.WriteString(nextBg)
			result.WriteString(ResolveColor(arrowFg))
			result.WriteString(separator)
		}
	}

	result.WriteString(Reset)

	// Add final arrow to default background
	if len(nonEmpty) > 0 {
		lastBg := nonEmpty[len(nonEmpty)-1].Bg
		if lastBg != "" {
			result.WriteString(ResolveColor(lastBg))
			result.WriteString(separator)
			result.WriteString(Reset)
		}
	}

	// Truncate to terminal width
	width := TerminalWidth()
	return Truncate(result.String(), width)
}

// PowerlineSplitLayout renders left and right aligned powerline segments.
func PowerlineSplitLayout(left, right []PowerlineSegment, leftSep, rightSep string) string {
	if leftSep == "" {
		leftSep = "" // U+E0B0
	}
	if rightSep == "" {
		rightSep = "" // U+E0B2
	}

	width := TerminalWidth()

	leftStr := PowerlineLayout(left, leftSep)
	rightStr := powerlineLayoutReverse(right, rightSep)

	if leftStr == "" && rightStr == "" {
		return ""
	}
	if rightStr == "" {
		return Truncate(leftStr, width)
	}
	if leftStr == "" {
		return Truncate(rightStr, width)
	}

	leftLen := VisibleLength(leftStr)
	rightLen := VisibleLength(rightStr)

	minGap := 2
	totalContentLen := leftLen + rightLen + minGap

	if totalContentLen > width {
		return Truncate(leftStr+" "+rightStr, width)
	}

	padding := width - leftLen - rightLen
	if padding < minGap {
		padding = minGap
	}

	return leftStr + strings.Repeat(" ", padding) + rightStr
}

// powerlineLayoutReverse renders segments in reverse order for right alignment.
func powerlineLayoutReverse(segments []PowerlineSegment, separator string) string {
	if separator == "" {
		separator = "" // U+E0B2 (right arrow)
	}

	// Filter non-empty
	var nonEmpty []PowerlineSegment
	for _, s := range segments {
		if s.Text != "" {
			nonEmpty = append(nonEmpty, s)
		}
	}

	if len(nonEmpty) == 0 {
		return ""
	}

	var result strings.Builder

	// First segment gets a leading arrow
	if len(nonEmpty) > 0 {
		firstBg := nonEmpty[0].Bg
		if firstBg != "" {
			result.WriteString(ResolveColor(firstBg))
			result.WriteString(separator)
		}
	}

	for i, seg := range nonEmpty {
		bg := ResolveBgColor(seg.Bg)
		fg := ResolveColor(seg.Fg)

		result.WriteString(bg)
		result.WriteString(fg)
		result.WriteString(" ")
		result.WriteString(seg.Text)
		result.WriteString(" ")

		if i < len(nonEmpty)-1 {
			arrowFg := nonEmpty[i+1].Bg
			if arrowFg == "" {
				arrowFg = "black"
			}
			result.WriteString(Reset)
			result.WriteString(bg)
			result.WriteString(ResolveColor(arrowFg))
			result.WriteString(separator)
		}
	}

	result.WriteString(Reset)
	return result.String()
}
