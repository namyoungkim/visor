package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/namyoungkim/visor/internal/config"
)

// View represents the current view mode
type View int

const (
	ViewMain View = iota
	ViewAddWidget
	ViewEditOptions
	ViewLayoutPicker
	ViewHelp
	ViewConfirmQuit
)

// Side represents left or right side in split layout
type Side int

const (
	SideWidgets Side = iota // For single-line layout
	SideLeft
	SideRight
)

// Model holds all TUI state
type Model struct {
	// Configuration
	config     *config.Config
	configPath string
	dirty      bool // Has unsaved changes

	// Current view
	view     View
	prevView View // For going back

	// Main view state
	lineIndex   int  // Current line
	widgetIndex int  // Current widget within line
	side        Side // Current side (for split layout)

	// Add widget state
	addWidgetCursor int

	// Edit options state
	editWidget       *config.WidgetConfig
	editWidgetMeta   *WidgetMeta
	editOptionIndex  int
	editInputs       []textinput.Model
	editFocusedInput int

	// Layout picker state
	layoutChoice int // 0 = single, 1 = split

	// Help
	keys KeyMap
	help help.Model

	// Terminal size
	width  int
	height int

	// Status message
	statusMsg   string
	statusIsErr bool
}

// NewModel creates a new TUI model
func NewModel(cfg *config.Config, configPath string) Model {
	h := help.New()
	h.ShowAll = false

	return Model{
		config:     cfg,
		configPath: configPath,
		dirty:      false,
		view:       ViewMain,
		keys:       DefaultKeyMap(),
		help:       h,
		width:      80,
		height:     24,
	}
}

// currentLine returns the current line or nil
func (m *Model) currentLine() *config.Line {
	if m.lineIndex < 0 || m.lineIndex >= len(m.config.Lines) {
		return nil
	}
	return &m.config.Lines[m.lineIndex]
}

// currentWidgets returns the widgets for the current side
func (m *Model) currentWidgets() []config.WidgetConfig {
	line := m.currentLine()
	if line == nil {
		return nil
	}

	// Check if split layout
	if len(line.Left) > 0 || len(line.Right) > 0 {
		switch m.side {
		case SideLeft:
			return line.Left
		case SideRight:
			return line.Right
		}
	}

	return line.Widgets
}

// setCurrentWidgets sets widgets for the current side
func (m *Model) setCurrentWidgets(widgets []config.WidgetConfig) {
	line := m.currentLine()
	if line == nil {
		return
	}

	// Check if split layout
	if len(line.Left) > 0 || len(line.Right) > 0 {
		switch m.side {
		case SideLeft:
			line.Left = widgets
		case SideRight:
			line.Right = widgets
		}
		return
	}

	line.Widgets = widgets
}

// currentWidget returns the current widget or nil
func (m *Model) currentWidget() *config.WidgetConfig {
	widgets := m.currentWidgets()
	if m.widgetIndex < 0 || m.widgetIndex >= len(widgets) {
		return nil
	}
	return &widgets[m.widgetIndex]
}

// isSplitLayout returns true if current line uses split layout
func (m *Model) isSplitLayout() bool {
	line := m.currentLine()
	if line == nil {
		return false
	}
	return len(line.Left) > 0 || len(line.Right) > 0
}

// widgetCount returns the number of widgets on current side
func (m *Model) widgetCount() int {
	return len(m.currentWidgets())
}

// totalItems returns total selectable items (widgets + add button)
func (m *Model) totalItems() int {
	return m.widgetCount() + 1 // +1 for "Add widget" button
}

// isAddButtonSelected returns true if the add button is selected
func (m *Model) isAddButtonSelected() bool {
	return m.widgetIndex >= m.widgetCount()
}

// markDirty marks the config as having unsaved changes
func (m *Model) markDirty() {
	m.dirty = true
}

// setStatus sets a status message
func (m *Model) setStatus(msg string, isErr bool) {
	m.statusMsg = msg
	m.statusIsErr = isErr
}

// clearStatus clears the status message
func (m *Model) clearStatus() {
	m.statusMsg = ""
	m.statusIsErr = false
}
