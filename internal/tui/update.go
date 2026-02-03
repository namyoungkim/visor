package tui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/theme"
)

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		return m, nil

	case tea.KeyMsg:
		// Clear status on any key press
		m.clearStatus()

		switch m.view {
		case ViewMain:
			return m.updateMain(msg)
		case ViewAddWidget:
			return m.updateAddWidget(msg)
		case ViewEditOptions:
			return m.updateEditOptions(msg)
		case ViewLayoutPicker:
			return m.updateLayoutPicker(msg)
		case ViewThemePicker:
			return m.updateThemePicker(msg)
		case ViewHelp:
			return m.updateHelp(msg)
		case ViewConfirmQuit:
			return m.updateConfirmQuit(msg)
		}
	}

	return m, nil
}

// updateMain handles input in the main view
func (m Model) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		if m.dirty {
			m.prevView = m.view
			m.view = ViewConfirmQuit
			return m, nil
		}
		return m, tea.Quit

	case key.Matches(msg, m.keys.Help):
		m.prevView = m.view
		m.view = ViewHelp
		return m, nil

	case key.Matches(msg, m.keys.Up):
		if m.widgetIndex > 0 {
			m.widgetIndex--
		}
		return m, nil

	case key.Matches(msg, m.keys.Down):
		if m.widgetIndex < m.totalItems()-1 {
			m.widgetIndex++
		}
		return m, nil

	case key.Matches(msg, m.keys.Left):
		// Move to previous line or switch side
		if m.isSplitLayout() && m.side == SideRight {
			m.side = SideLeft
			m.widgetIndex = 0
		} else if m.lineIndex > 0 {
			m.lineIndex--
			m.widgetIndex = 0
			m.side = SideWidgets
			if m.isSplitLayout() {
				m.side = SideLeft
			}
		}
		return m, nil

	case key.Matches(msg, m.keys.Right):
		// Move to next line or switch side
		if m.isSplitLayout() && m.side == SideLeft {
			m.side = SideRight
			m.widgetIndex = 0
		} else if m.lineIndex < len(m.config.Lines)-1 {
			m.lineIndex++
			m.widgetIndex = 0
			m.side = SideWidgets
			if m.isSplitLayout() {
				m.side = SideLeft
			}
		}
		return m, nil

	case key.Matches(msg, m.keys.TabSwitch):
		// Toggle side in split layout
		if m.isSplitLayout() {
			if m.side == SideLeft {
				m.side = SideRight
			} else {
				m.side = SideLeft
			}
			m.widgetIndex = 0
		}
		return m, nil

	case key.Matches(msg, m.keys.Toggle):
		// Toggle not supported - widgets are always enabled if in config
		return m, nil

	case key.Matches(msg, m.keys.Add):
		m.prevView = m.view
		m.view = ViewAddWidget
		m.addWidgetCursor = 0
		return m, nil

	case key.Matches(msg, m.keys.Edit):
		if m.isAddButtonSelected() {
			// Open add widget view
			m.prevView = m.view
			m.view = ViewAddWidget
			m.addWidgetCursor = 0
			return m, nil
		}
		widget := m.currentWidget()
		if widget != nil {
			meta := GetWidgetMeta(widget.Name)
			if meta != nil && len(meta.Options) > 0 {
				m.prevView = m.view
				m.view = ViewEditOptions
				m.editWidget = widget
				m.editWidgetMeta = meta
				m.editOptionIndex = 0
				m.initEditInputs()
				return m, nil
			}
		}
		return m, nil

	case key.Matches(msg, m.keys.Delete):
		if !m.isAddButtonSelected() {
			m.deleteCurrentWidget()
		}
		return m, nil

	case key.Matches(msg, m.keys.MoveUp):
		if !m.isAddButtonSelected() && m.widgetIndex > 0 {
			m.swapWidgets(m.widgetIndex, m.widgetIndex-1)
			m.widgetIndex--
		}
		return m, nil

	case key.Matches(msg, m.keys.MoveDown):
		if !m.isAddButtonSelected() && m.widgetIndex < m.widgetCount()-1 {
			m.swapWidgets(m.widgetIndex, m.widgetIndex+1)
			m.widgetIndex++
		}
		return m, nil

	case key.Matches(msg, m.keys.NewLine):
		m.addNewLine()
		return m, nil

	case key.Matches(msg, m.keys.Layout):
		m.prevView = m.view
		m.view = ViewLayoutPicker
		if m.isSplitLayout() {
			m.layoutChoice = 1
		} else {
			m.layoutChoice = 0
		}
		return m, nil

	case key.Matches(msg, m.keys.Theme):
		m.prevView = m.view
		m.view = ViewThemePicker
		// Find current theme index
		themes := theme.List()
		m.themeChoice = 0
		for i, t := range themes {
			if t == m.config.Theme.Name {
				m.themeChoice = i
				break
			}
		}
		return m, nil

	case key.Matches(msg, m.keys.Save):
		if err := config.Save(m.config, m.configPath); err != nil {
			m.setStatus("Error: "+err.Error(), true)
		} else {
			m.dirty = false
			m.setStatus("Saved!", false)
		}
		return m, nil
	}

	return m, nil
}

// updateAddWidget handles input in the add widget view
func (m Model) updateAddWidget(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	allWidgets := AllWidgets()

	switch {
	case key.Matches(msg, m.keys.Back):
		m.view = m.prevView
		return m, nil

	case key.Matches(msg, m.keys.Up):
		if m.addWidgetCursor > 0 {
			m.addWidgetCursor--
		}
		return m, nil

	case key.Matches(msg, m.keys.Down):
		if m.addWidgetCursor < len(allWidgets)-1 {
			m.addWidgetCursor++
		}
		return m, nil

	case key.Matches(msg, m.keys.Confirm):
		if m.addWidgetCursor >= 0 && m.addWidgetCursor < len(allWidgets) {
			widgetName := allWidgets[m.addWidgetCursor].Name
			m.addWidget(widgetName)
			m.view = m.prevView
		}
		return m, nil
	}

	return m, nil
}

// updateEditOptions handles input in the edit options view
func (m Model) updateEditOptions(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		m.view = m.prevView
		return m, nil

	case key.Matches(msg, m.keys.Up):
		if m.editFocusedInput > 0 {
			m.editInputs[m.editFocusedInput].Blur()
			m.editFocusedInput--
			m.editInputs[m.editFocusedInput].Focus()
		}
		return m, nil

	case key.Matches(msg, m.keys.Down), key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
		if m.editFocusedInput < len(m.editInputs)-1 {
			m.editInputs[m.editFocusedInput].Blur()
			m.editFocusedInput++
			m.editInputs[m.editFocusedInput].Focus()
		}
		return m, nil

	case key.Matches(msg, m.keys.Save), key.Matches(msg, m.keys.Confirm):
		m.applyEditedOptions()
		m.view = m.prevView
		return m, nil

	default:
		// Pass to focused input
		if m.editFocusedInput >= 0 && m.editFocusedInput < len(m.editInputs) {
			var cmd tea.Cmd
			m.editInputs[m.editFocusedInput], cmd = m.editInputs[m.editFocusedInput].Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// updateLayoutPicker handles input in layout picker view
func (m Model) updateLayoutPicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		m.view = m.prevView
		return m, nil

	case key.Matches(msg, m.keys.Up), key.Matches(msg, m.keys.Left):
		if m.layoutChoice > 0 {
			m.layoutChoice--
		}
		return m, nil

	case key.Matches(msg, m.keys.Down), key.Matches(msg, m.keys.Right):
		if m.layoutChoice < 1 {
			m.layoutChoice++
		}
		return m, nil

	case key.Matches(msg, m.keys.Confirm):
		m.applyLayout(m.layoutChoice == 1)
		m.view = m.prevView
		return m, nil
	}

	return m, nil
}

// updateThemePicker handles input in theme picker view
func (m Model) updateThemePicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	themes := theme.List()

	switch {
	case key.Matches(msg, m.keys.Back):
		m.view = m.prevView
		return m, nil

	case key.Matches(msg, m.keys.Up):
		if m.themeChoice > 0 {
			m.themeChoice--
		}
		return m, nil

	case key.Matches(msg, m.keys.Down):
		if m.themeChoice < len(themes)-1 {
			m.themeChoice++
		}
		return m, nil

	case key.Matches(msg, m.keys.Confirm):
		if m.themeChoice >= 0 && m.themeChoice < len(themes) {
			themeName := themes[m.themeChoice]
			t := theme.Get(themeName)
			m.config.Theme.Name = themeName
			m.config.Theme.Powerline = t.Powerline
			m.markDirty()
		}
		m.view = m.prevView
		return m, nil
	}

	return m, nil
}

// updateHelp handles input in help view
func (m Model) updateHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back), key.Matches(msg, m.keys.Help):
		m.view = m.prevView
		return m, nil
	}
	return m, nil
}

// updateConfirmQuit handles input in confirm quit view
func (m Model) updateConfirmQuit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		return m, tea.Quit
	case "n", "N", "esc":
		m.view = m.prevView
		return m, nil
	case "s", "S":
		if err := config.Save(m.config, m.configPath); err != nil {
			m.setStatus("Error: "+err.Error(), true)
			m.view = m.prevView
		} else {
			return m, tea.Quit
		}
	}
	return m, nil
}

// Helper methods for model manipulation

func (m *Model) deleteCurrentWidget() {
	widgets := m.currentWidgets()
	if m.widgetIndex < 0 || m.widgetIndex >= len(widgets) {
		return
	}

	// Create a new slice to avoid modifying the original backing array
	newWidgets := make([]config.WidgetConfig, 0, len(widgets)-1)
	newWidgets = append(newWidgets, widgets[:m.widgetIndex]...)
	newWidgets = append(newWidgets, widgets[m.widgetIndex+1:]...)
	m.setCurrentWidgets(newWidgets)
	m.markDirty()

	if m.widgetIndex >= len(newWidgets) && m.widgetIndex > 0 {
		m.widgetIndex--
	}
}

func (m *Model) swapWidgets(i, j int) {
	widgets := m.currentWidgets()
	if i < 0 || i >= len(widgets) || j < 0 || j >= len(widgets) {
		return
	}

	widgets[i], widgets[j] = widgets[j], widgets[i]
	m.setCurrentWidgets(widgets)
	m.markDirty()
}

func (m *Model) addWidget(name string) {
	widgets := m.currentWidgets()
	newWidget := config.WidgetConfig{Name: name}
	widgets = append(widgets, newWidget)
	m.setCurrentWidgets(widgets)
	m.markDirty()
	m.widgetIndex = len(widgets) - 1
}

func (m *Model) addNewLine() {
	newLine := config.Line{
		Widgets: []config.WidgetConfig{},
	}
	m.config.Lines = append(m.config.Lines, newLine)
	m.lineIndex = len(m.config.Lines) - 1
	m.widgetIndex = 0
	m.side = SideWidgets
	m.markDirty()
}

func (m *Model) applyLayout(split bool) {
	line := m.currentLine()
	if line == nil {
		return
	}

	if split && !m.isSplitLayout() {
		// Convert to split layout
		line.Left = line.Widgets
		line.Right = []config.WidgetConfig{}
		line.Widgets = nil
		m.side = SideLeft
		m.markDirty()
	} else if !split && m.isSplitLayout() {
		// Convert to single layout
		line.Widgets = append(line.Left, line.Right...)
		line.Left = nil
		line.Right = nil
		m.side = SideWidgets
		m.markDirty()
	}
}

func (m *Model) initEditInputs() {
	if m.editWidgetMeta == nil {
		return
	}

	m.editInputs = make([]textinput.Model, len(m.editWidgetMeta.Options))
	for i, opt := range m.editWidgetMeta.Options {
		ti := textinput.New()
		ti.Placeholder = opt.DefaultValue
		ti.CharLimit = 20

		// Get current value or default
		if m.editWidget.Extra != nil {
			if v, ok := m.editWidget.Extra[opt.Key]; ok {
				ti.SetValue(v)
			}
		}
		if ti.Value() == "" {
			ti.SetValue(opt.DefaultValue)
		}

		if i == 0 {
			ti.Focus()
		}
		m.editInputs[i] = ti
	}
	m.editFocusedInput = 0
}

func (m *Model) applyEditedOptions() {
	if m.editWidget == nil || m.editWidgetMeta == nil {
		return
	}

	if m.editWidget.Extra == nil {
		m.editWidget.Extra = make(map[string]string)
	}

	for i, opt := range m.editWidgetMeta.Options {
		if i < len(m.editInputs) {
			value := m.editInputs[i].Value()

			// Validate based on option type
			if !m.validateOptionValue(value, opt.Type) {
				continue // Skip invalid values
			}

			if value != opt.DefaultValue && value != "" {
				m.editWidget.Extra[opt.Key] = value
			} else {
				delete(m.editWidget.Extra, opt.Key)
			}
		}
	}

	// Clean up empty Extra map
	if len(m.editWidget.Extra) == 0 {
		m.editWidget.Extra = nil
	}

	m.markDirty()
}

// validateOptionValue validates the value based on the option type
func (m *Model) validateOptionValue(value string, optType OptionType) bool {
	if value == "" {
		return true // Empty is always valid (uses default)
	}

	switch optType {
	case OptionTypeBool:
		return value == "true" || value == "false"
	case OptionTypeInt:
		_, err := strconv.Atoi(value)
		return err == nil
	case OptionTypeFloat:
		_, err := strconv.ParseFloat(value, 64)
		return err == nil
	case OptionTypeString:
		return true
	}
	return true
}
