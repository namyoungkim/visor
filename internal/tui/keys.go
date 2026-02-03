package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the TUI
type KeyMap struct {
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	Toggle    key.Binding
	Edit      key.Binding
	Add       key.Binding
	Delete    key.Binding
	MoveUp    key.Binding
	MoveDown  key.Binding
	Layout    key.Binding
	NewLine   key.Binding
	TabSwitch key.Binding
	Save      key.Binding
	Help      key.Binding
	Quit      key.Binding
	Back      key.Binding
	Confirm   key.Binding
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		Toggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e", "enter"),
			key.WithHelp("e/enter", "edit"),
		),
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d", "backspace"),
			key.WithHelp("d", "delete"),
		),
		MoveUp: key.NewBinding(
			key.WithKeys("K"),
			key.WithHelp("K", "move up"),
		),
		MoveDown: key.NewBinding(
			key.WithKeys("J"),
			key.WithHelp("J", "move down"),
		),
		Layout: key.NewBinding(
			key.WithKeys("L"),
			key.WithHelp("L", "layout"),
		),
		NewLine: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new line"),
		),
		TabSwitch: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "left/right"),
		),
		Save: key.NewBinding(
			key.WithKeys("s", "ctrl+s"),
			key.WithHelp("s", "save"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c", "ctrl+q"),
			key.WithHelp("q", "quit"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
	}
}

// ShortHelp returns the short help text for the main view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Toggle, k.Edit, k.Add, k.Delete, k.Save, k.Quit}
}

// FullHelp returns the full help text
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Toggle, k.Edit, k.Add, k.Delete},
		{k.MoveUp, k.MoveDown, k.Layout, k.NewLine},
		{k.Save, k.Help, k.Quit},
	}
}
