package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/namyoungkim/visor/internal/config"
)

// View renders the current state
func (m Model) View() string {
	switch m.view {
	case ViewAddWidget:
		return m.viewAddWidget()
	case ViewEditOptions:
		return m.viewEditOptions()
	case ViewLayoutPicker:
		return m.viewLayoutPicker()
	case ViewHelp:
		return m.viewHelp()
	case ViewConfirmQuit:
		return m.viewConfirmQuit()
	default:
		return m.viewMain()
	}
}

func (m Model) viewMain() string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("visor Configuration")
	if m.dirty {
		title += warningStyle.Render(" [modified]")
	}
	b.WriteString(title + "\n\n")

	// Lines and widgets
	for lineIdx, line := range m.config.Lines {
		isCurrentLine := lineIdx == m.lineIndex

		// Line header
		lineHeader := fmt.Sprintf("Line %d", lineIdx+1)
		if m.isSplitLayoutForLine(&line) {
			lineHeader += " (split)"
		}
		if isCurrentLine {
			b.WriteString(lineHeaderStyle.Render(lineHeader) + "\n")
		} else {
			b.WriteString(disabledStyle.Render(lineHeader) + "\n")
		}

		// Render widgets based on layout
		if m.isSplitLayoutForLine(&line) {
			// Split layout - show left and right sections
			b.WriteString(m.renderSplitLine(&line, lineIdx, isCurrentLine))
		} else {
			// Single layout
			b.WriteString(m.renderWidgetList(line.Widgets, lineIdx, isCurrentLine, SideWidgets))
		}

		b.WriteString("\n")
	}

	// Preview section
	b.WriteString(sectionStyle.Render("Preview") + "\n")
	preview := RenderPreview(m.config)
	b.WriteString(previewBoxStyle.Render(preview) + "\n\n")

	// Status message
	if m.statusMsg != "" {
		if m.statusIsErr {
			b.WriteString(errorStyle.Render(m.statusMsg) + "\n")
		} else {
			b.WriteString(successStyle.Render(m.statusMsg) + "\n")
		}
	}

	// Help
	helpText := "[j/k] Move  [e] Edit  [a] Add  [d] Delete  [J/K] Reorder  [l] Layout  [n] New line  [s] Save  [?] Help  [q] Quit"
	b.WriteString(helpStyle.Render(helpText))

	return b.String()
}

func (m Model) isSplitLayoutForLine(line *config.Line) bool {
	return len(line.Left) > 0 || len(line.Right) > 0
}

func (m Model) renderSplitLine(line *config.Line, lineIdx int, isCurrentLine bool) string {
	var b strings.Builder

	// Left side
	leftHeader := "  Left:"
	if isCurrentLine && m.side == SideLeft {
		b.WriteString(selectedStyle.Render(leftHeader) + "\n")
	} else {
		b.WriteString(disabledStyle.Render(leftHeader) + "\n")
	}
	b.WriteString(m.renderWidgetList(line.Left, lineIdx, isCurrentLine && m.side == SideLeft, SideLeft))

	// Right side
	rightHeader := "  Right:"
	if isCurrentLine && m.side == SideRight {
		b.WriteString(selectedStyle.Render(rightHeader) + "\n")
	} else {
		b.WriteString(disabledStyle.Render(rightHeader) + "\n")
	}
	b.WriteString(m.renderWidgetList(line.Right, lineIdx, isCurrentLine && m.side == SideRight, SideRight))

	return b.String()
}

func (m Model) renderWidgetList(widgets []config.WidgetConfig, lineIdx int, isActive bool, side Side) string {
	var b strings.Builder

	indent := "    "

	for i, widget := range widgets {
		isSelected := isActive && m.widgetIndex == i && m.side == side

		prefix := indent
		if isSelected {
			prefix = indent[:2] + cursorStyle.Render() + " "
		}

		name := widget.Name
		if isSelected {
			name = selectedStyle.Render(name)
		} else if isActive {
			name = normalStyle.Render(name)
		} else {
			name = disabledStyle.Render(name)
		}

		// Add option indicator if widget has custom options
		if widget.Extra != nil && len(widget.Extra) > 0 {
			name += disabledStyle.Render(" *")
		}

		b.WriteString(prefix + name + "\n")
	}

	// Add widget button
	isAddSelected := isActive && m.widgetIndex >= len(widgets) && m.side == side
	addPrefix := indent
	if isAddSelected {
		addPrefix = indent[:2] + cursorStyle.Render() + " "
	}

	addText := "[+] Add widget..."
	if isAddSelected {
		addText = selectedStyle.Render(addText)
	} else if isActive {
		addText = normalStyle.Render(addText)
	} else {
		addText = disabledStyle.Render(addText)
	}
	b.WriteString(addPrefix + addText + "\n")

	return b.String()
}

func (m Model) viewAddWidget() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Add Widget") + "\n\n")

	allWidgets := AllWidgets()
	for i, w := range allWidgets {
		prefix := "  "
		if i == m.addWidgetCursor {
			prefix = cursorStyle.Render() + " "
		}

		name := w.Name
		desc := disabledStyle.Render(" - " + w.Description)

		if i == m.addWidgetCursor {
			name = selectedStyle.Render(name)
		} else {
			name = normalStyle.Render(name)
		}

		b.WriteString(prefix + name + desc + "\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("[j/k] Move  [enter] Add  [esc] Cancel"))

	return b.String()
}

func (m Model) viewEditOptions() string {
	var b strings.Builder

	if m.editWidget == nil || m.editWidgetMeta == nil {
		return "No widget selected"
	}

	b.WriteString(titleStyle.Render("Edit: "+m.editWidget.Name) + "\n\n")

	for i, opt := range m.editWidgetMeta.Options {
		if i >= len(m.editInputs) {
			break
		}

		label := fmt.Sprintf("%-20s", opt.Key+":")
		if i == m.editFocusedInput {
			label = selectedStyle.Render(label)
		} else {
			label = normalStyle.Render(label)
		}

		var inputView string
		if i == m.editFocusedInput {
			inputView = focusedInputStyle.Render(m.editInputs[i].View())
		} else {
			inputView = inputStyle.Render(m.editInputs[i].View())
		}

		desc := disabledStyle.Render("  " + opt.Description)

		b.WriteString(label + inputView + desc + "\n")
	}

	// Preview
	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("Preview") + "\n")
	preview := RenderWidgetPreview(m.editWidget)
	b.WriteString(previewBoxStyle.Render(preview) + "\n\n")

	b.WriteString(helpStyle.Render("[tab/↓] Next  [↑] Previous  [enter] Save  [esc] Cancel"))

	return b.String()
}

func (m Model) viewLayoutPicker() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Select Layout") + "\n\n")

	choices := []struct {
		name string
		desc string
	}{
		{"Single Line", "All widgets in a row"},
		{"Split Layout", "Left-aligned and right-aligned widgets"},
	}

	for i, c := range choices {
		prefix := "  "
		if i == m.layoutChoice {
			prefix = cursorStyle.Render() + " "
		}

		name := c.name
		desc := disabledStyle.Render(" - " + c.desc)

		if i == m.layoutChoice {
			name = selectedStyle.Render(name)
		} else {
			name = normalStyle.Render(name)
		}

		b.WriteString(prefix + name + desc + "\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("[j/k] Move  [enter] Select  [esc] Cancel"))

	return b.String()
}

func (m Model) viewHelp() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Help") + "\n\n")

	sections := []struct {
		title string
		keys  []string
	}{
		{
			"Navigation",
			[]string{
				"j/k, ↑/↓    Move cursor up/down",
				"h/l, ←/→    Move between lines/sides",
				"Tab         Switch left/right (split layout)",
			},
		},
		{
			"Widget Management",
			[]string{
				"a           Add widget",
				"e, Enter    Edit widget options",
				"d           Delete widget",
				"J/K         Move widget up/down",
			},
		},
		{
			"Line Management",
			[]string{
				"n           Add new line",
				"l           Change layout (single/split)",
			},
		},
		{
			"File Operations",
			[]string{
				"s, Ctrl+S   Save configuration",
				"q, Ctrl+Q   Quit (prompts if unsaved)",
			},
		},
	}

	for _, sec := range sections {
		b.WriteString(sectionStyle.Render(sec.title) + "\n")
		for _, k := range sec.keys {
			b.WriteString("  " + k + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(helpStyle.Render("[?/esc] Close help"))

	return b.String()
}

func (m Model) viewConfirmQuit() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(warningColor).
		Padding(1, 2)

	content := warningStyle.Render("You have unsaved changes!") + "\n\n"
	content += "[s] Save and quit\n"
	content += "[y] Quit without saving\n"
	content += "[n] Cancel"

	return "\n" + box.Render(content)
}
