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
		contains string
	}{
		{
			name: "single running agent",
			agents: []transcript.Agent{
				{ID: "1", Type: "Explore", Status: "running"},
			},
			contains: "◐ 1 agent",
		},
		{
			name: "multiple running agents (plural)",
			agents: []transcript.Agent{
				{ID: "1", Type: "Explore", Status: "running"},
				{ID: "2", Type: "Plan", Status: "running"},
			},
			contains: "◐ 2 agents",
		},
		{
			name: "single completed agent",
			agents: []transcript.Agent{
				{ID: "1", Type: "Explore", Status: "completed"},
			},
			contains: "✓ 1 done",
		},
		{
			name: "multiple completed agents",
			agents: []transcript.Agent{
				{ID: "1", Type: "Explore", Status: "completed"},
				{ID: "2", Type: "Plan", Status: "completed"},
			},
			contains: "✓ 2 done",
		},
		{
			name: "mixed running and completed",
			agents: []transcript.Agent{
				{ID: "1", Type: "Explore", Status: "running"},
				{ID: "2", Type: "Plan", Status: "completed"},
				{ID: "3", Type: "Bash", Status: "completed"},
			},
			contains: "◐ 1 | ✓ 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &AgentsWidget{}
			w.SetTranscript(&transcript.Data{Agents: tt.agents})

			result := w.Render(&input.Session{}, &config.WidgetConfig{})

			if !strings.Contains(result, tt.contains) {
				t.Errorf("Expected result to contain '%s', got '%s'", tt.contains, result)
			}
		})
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
