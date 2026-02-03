package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/namyoungkim/visor/internal/config"
)

// Run starts the TUI configuration editor
func Run() error {
	configPath := config.DefaultConfigPath()

	// Load existing config or create default
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create a deep copy to track changes
	cfg = config.DeepCopy(cfg)

	// Create and run the TUI
	model := NewModel(cfg, configPath)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
