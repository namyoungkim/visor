package widgets

import (
	"strings"
	"testing"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
)

func TestModelWidget_WithDisplayName(t *testing.T) {
	w := &ModelWidget{}
	session := &input.Session{
		Model: input.Model{
			DisplayName: "Opus",
			ID:          "claude-opus-4",
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "Opus") {
		t.Errorf("Expected 'Opus', got '%s'", result)
	}
}

func TestModelWidget_WithIDOnly(t *testing.T) {
	w := &ModelWidget{}
	session := &input.Session{
		Model: input.Model{
			DisplayName: "",
			ID:          "claude-opus-4",
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "claude-opus-4") {
		t.Errorf("Expected 'claude-opus-4', got '%s'", result)
	}
}

func TestModelWidget_Empty(t *testing.T) {
	w := &ModelWidget{}
	session := &input.Session{}

	result := w.Render(session, &config.WidgetConfig{})

	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestModelWidget_ShouldRender(t *testing.T) {
	w := &ModelWidget{}

	tests := []struct {
		name     string
		session  *input.Session
		expected bool
	}{
		{"with display name", &input.Session{Model: input.Model{DisplayName: "Opus"}}, true},
		{"with id only", &input.Session{Model: input.Model{ID: "claude-4"}}, true},
		{"empty", &input.Session{}, false},
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
