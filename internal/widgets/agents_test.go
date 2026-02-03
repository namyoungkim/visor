package widgets

import (
	"strings"
	"testing"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/transcript"
)

func TestAgentsWidget_Render(t *testing.T) {
	tests := []struct {
		name     string
		agents   []transcript.Agent
		contains []string
	}{
		{
			name: "single running agent shows type",
			agents: []transcript.Agent{
				{ID: "1", Type: "Explore", Status: "running"},
			},
			contains: []string{"◐Explore"},
		},
		{
			name: "multiple running agents show types",
			agents: []transcript.Agent{
				{ID: "1", Type: "Explore", Status: "running"},
				{ID: "2", Type: "Plan", Status: "running"},
			},
			contains: []string{"◐Explore", "◐Plan"},
		},
		{
			name: "completed agent shows checkmark",
			agents: []transcript.Agent{
				{ID: "1", Type: "Explore", Status: "completed"},
			},
			contains: []string{"✓Explore"},
		},
		{
			name: "mixed status shows both icons",
			agents: []transcript.Agent{
				{ID: "1", Type: "Plan", Status: "completed"},
				{ID: "2", Type: "Explore", Status: "running"},
			},
			contains: []string{"✓Plan", "◐Explore"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &AgentsWidget{}
			w.SetTranscript(&transcript.Data{Agents: tt.agents})

			result := w.Render(&input.Session{}, &config.WidgetConfig{})

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain '%s', got '%s'", expected, result)
				}
			}
		})
	}
}

func TestAgentsWidget_MaxDisplay(t *testing.T) {
	w := &AgentsWidget{}
	w.SetTranscript(&transcript.Data{
		Agents: []transcript.Agent{
			{ID: "1", Type: "Plan", Status: "completed"},
			{ID: "2", Type: "Explore", Status: "completed"},
			{ID: "3", Type: "Bash", Status: "completed"},
			{ID: "4", Type: "general-purpose", Status: "running"},
		},
	})

	cfg := &config.WidgetConfig{
		Extra: map[string]string{"max_display": "2"},
	}

	result := w.Render(&input.Session{}, cfg)

	// Should only show last 2 agents
	if strings.Contains(result, "Plan") {
		t.Error("Expected Plan to be hidden (max_display=2)")
	}
	if !strings.Contains(result, "Bash") || !strings.Contains(result, "general-purpose") {
		t.Errorf("Expected last 2 agents (Bash, general-purpose), got '%s'", result)
	}
}

func TestAgentsWidget_ShowLabel(t *testing.T) {
	w := &AgentsWidget{}
	w.SetTranscript(&transcript.Data{
		Agents: []transcript.Agent{
			{ID: "1", Type: "Explore", Status: "running"},
		},
	})

	cfg := &config.WidgetConfig{
		Extra: map[string]string{"show_label": "true"},
	}

	result := w.Render(&input.Session{}, cfg)

	if !strings.Contains(result, "Agents:") {
		t.Errorf("Expected 'Agents:' label, got '%s'", result)
	}
}

func TestAgentsWidget_Empty(t *testing.T) {
	w := &AgentsWidget{}

	// nil transcript
	result := w.Render(&input.Session{}, &config.WidgetConfig{})
	if result != "" {
		t.Errorf("Expected empty string for nil transcript, got '%s'", result)
	}

	// empty agents
	w.SetTranscript(&transcript.Data{Agents: []transcript.Agent{}})
	result = w.Render(&input.Session{}, &config.WidgetConfig{})
	if result != "" {
		t.Errorf("Expected empty string for empty agents, got '%s'", result)
	}
}

func TestAgentsWidget_ShouldRender(t *testing.T) {
	w := &AgentsWidget{}

	tests := []struct {
		name       string
		transcript *transcript.Data
		expected   bool
	}{
		{"nil transcript", nil, false},
		{"empty agents", &transcript.Data{Agents: []transcript.Agent{}}, false},
		{"with agents", &transcript.Data{Agents: []transcript.Agent{{ID: "1", Type: "Explore"}}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w.SetTranscript(tt.transcript)
			result := w.ShouldRender(&input.Session{}, &config.WidgetConfig{})
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
