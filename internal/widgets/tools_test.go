package widgets

import (
	"strings"
	"testing"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/transcript"
)

func TestToolsWidget_Render(t *testing.T) {
	tests := []struct {
		name     string
		tools    []transcript.Tool
		contains []string
	}{
		{
			name: "completed tools show full name",
			tools: []transcript.Tool{
				{ID: "1", Name: "Read", Status: transcript.ToolCompleted},
				{ID: "2", Name: "Write", Status: transcript.ToolCompleted},
			},
			contains: []string{"✓Read", "✓Write"},
		},
		{
			name: "running tool shows spinner",
			tools: []transcript.Tool{
				{ID: "1", Name: "Bash", Status: transcript.ToolRunning},
			},
			contains: []string{"◐Bash"},
		},
		{
			name: "error tool shows X",
			tools: []transcript.Tool{
				{ID: "1", Name: "Edit", Status: transcript.ToolError},
			},
			contains: []string{"✗Edit"},
		},
		{
			name: "mixed status",
			tools: []transcript.Tool{
				{ID: "1", Name: "Read", Status: transcript.ToolCompleted},
				{ID: "2", Name: "Write", Status: transcript.ToolError},
				{ID: "3", Name: "Bash", Status: transcript.ToolRunning},
			},
			contains: []string{"✓Read", "✗Write", "◐Bash"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &ToolsWidget{}
			w.SetTranscript(&transcript.Data{Tools: tt.tools})

			result := w.Render(&input.Session{}, &config.WidgetConfig{})

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain '%s', got '%s'", expected, result)
				}
			}
		})
	}
}

func TestToolsWidget_MaxDisplay(t *testing.T) {
	w := &ToolsWidget{}
	w.SetTranscript(&transcript.Data{
		Tools: []transcript.Tool{
			{ID: "1", Name: "Read", Status: transcript.ToolCompleted},
			{ID: "2", Name: "Write", Status: transcript.ToolCompleted},
			{ID: "3", Name: "Edit", Status: transcript.ToolCompleted},
			{ID: "4", Name: "Bash", Status: transcript.ToolCompleted},
			{ID: "5", Name: "Glob", Status: transcript.ToolCompleted},
		},
	})

	cfg := &config.WidgetConfig{
		Extra: map[string]string{"max_display": "2"},
	}

	result := w.Render(&input.Session{}, cfg)

	// Should only show last 2 tools
	if strings.Contains(result, "Read") {
		t.Error("Expected Read to be hidden (max_display=2)")
	}
	if !strings.Contains(result, "Bash") || !strings.Contains(result, "Glob") {
		t.Errorf("Expected last 2 tools (Bash, Glob), got '%s'", result)
	}
}

func TestToolsWidget_ShowLabel(t *testing.T) {
	w := &ToolsWidget{}
	w.SetTranscript(&transcript.Data{
		Tools: []transcript.Tool{
			{ID: "1", Name: "Read", Status: transcript.ToolCompleted},
		},
	})

	cfg := &config.WidgetConfig{
		Extra: map[string]string{"show_label": "true"},
	}

	result := w.Render(&input.Session{}, cfg)

	if !strings.Contains(result, "Tools:") {
		t.Errorf("Expected 'Tools:' label, got '%s'", result)
	}
}

func TestToolsWidget_Empty(t *testing.T) {
	w := &ToolsWidget{}

	// nil transcript
	result := w.Render(&input.Session{}, &config.WidgetConfig{})
	if result != "" {
		t.Errorf("Expected empty string for nil transcript, got '%s'", result)
	}

	// empty tools
	w.SetTranscript(&transcript.Data{Tools: []transcript.Tool{}})
	result = w.Render(&input.Session{}, &config.WidgetConfig{})
	if result != "" {
		t.Errorf("Expected empty string for empty tools, got '%s'", result)
	}
}

func TestToolsWidget_ShouldRender(t *testing.T) {
	w := &ToolsWidget{}

	tests := []struct {
		name       string
		transcript *transcript.Data
		expected   bool
	}{
		{"nil transcript", nil, false},
		{"empty tools", &transcript.Data{Tools: []transcript.Tool{}}, false},
		{"with tools", &transcript.Data{Tools: []transcript.Tool{{ID: "1", Name: "Read", Count: 1}}}, true},
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

func TestToolsWidget_ShowCount(t *testing.T) {
	w := &ToolsWidget{}
	w.SetTranscript(&transcript.Data{
		Tools: []transcript.Tool{
			{ID: "1", Name: "Read", Status: transcript.ToolCompleted, Count: 7},
			{ID: "2", Name: "Edit", Status: transcript.ToolCompleted, Count: 4},
			{ID: "3", Name: "Bash", Status: transcript.ToolRunning, Count: 1},
		},
	})

	// Default: show_count = true
	result := w.Render(&input.Session{}, &config.WidgetConfig{})

	// Should show counts for tools with Count > 1
	if !strings.Contains(result, "×7") {
		t.Errorf("Expected '×7' for Read (count=7), got '%s'", result)
	}
	if !strings.Contains(result, "×4") {
		t.Errorf("Expected '×4' for Edit (count=4), got '%s'", result)
	}
	// Should NOT show count for Bash (count=1)
	if strings.Contains(result, "×1") {
		t.Errorf("Should not show '×1' for count=1, got '%s'", result)
	}
}

func TestToolsWidget_HideCount(t *testing.T) {
	w := &ToolsWidget{}
	w.SetTranscript(&transcript.Data{
		Tools: []transcript.Tool{
			{ID: "1", Name: "Read", Status: transcript.ToolCompleted, Count: 7},
		},
	})

	cfg := &config.WidgetConfig{
		Extra: map[string]string{"show_count": "false"},
	}

	result := w.Render(&input.Session{}, cfg)

	// Should NOT show count when show_count=false
	if strings.Contains(result, "×") {
		t.Errorf("Should not show count when show_count=false, got '%s'", result)
	}
}

func TestToolsWidget_PipeSeparator(t *testing.T) {
	w := &ToolsWidget{}
	w.SetTranscript(&transcript.Data{
		Tools: []transcript.Tool{
			{ID: "1", Name: "Read", Status: transcript.ToolCompleted, Count: 1},
			{ID: "2", Name: "Write", Status: transcript.ToolCompleted, Count: 1},
		},
	})

	result := w.Render(&input.Session{}, &config.WidgetConfig{})

	// Should use " | " separator
	if !strings.Contains(result, " | ") {
		t.Errorf("Expected ' | ' separator, got '%s'", result)
	}
}
