package widgets

import (
	"strings"
	"testing"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
)

func TestCodeChangesWidget_NoChanges(t *testing.T) {
	w := &CodeChangesWidget{}
	session := &input.Session{}

	result := w.Render(session, &config.WidgetConfig{})

	if result != "" {
		t.Errorf("Expected empty string for no changes, got '%s'", result)
	}
}

func TestCodeChangesWidget_OnlyAdded(t *testing.T) {
	w := &CodeChangesWidget{}
	session := &input.Session{
		Workspace: input.Workspace{
			LinesAdded:   10,
			LinesRemoved: 0,
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "+10") {
		t.Errorf("Expected +10, got '%s'", result)
	}
	if !strings.Contains(result, "-0") {
		t.Errorf("Expected -0, got '%s'", result)
	}
}

func TestCodeChangesWidget_BothAddedRemoved(t *testing.T) {
	w := &CodeChangesWidget{}
	session := &input.Session{
		Workspace: input.Workspace{
			LinesAdded:   25,
			LinesRemoved: 10,
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "+25") {
		t.Errorf("Expected +25, got '%s'", result)
	}
	if !strings.Contains(result, "-10") {
		t.Errorf("Expected -10, got '%s'", result)
	}
	// Check colors
	if !strings.Contains(result, "\033[32m") { // green
		t.Error("Expected green color for added")
	}
	if !strings.Contains(result, "\033[31m") { // red
		t.Error("Expected red color for removed")
	}
}

func TestCodeChangesWidget_ShouldRender(t *testing.T) {
	w := &CodeChangesWidget{}

	tests := []struct {
		name     string
		session  *input.Session
		expected bool
	}{
		{"no changes", &input.Session{}, false},
		{"only added", &input.Session{Workspace: input.Workspace{LinesAdded: 1}}, true},
		{"only removed", &input.Session{Workspace: input.Workspace{LinesRemoved: 1}}, true},
		{"both", &input.Session{Workspace: input.Workspace{LinesAdded: 1, LinesRemoved: 1}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := w.ShouldRender(tt.session, &config.WidgetConfig{})
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
