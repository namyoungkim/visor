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

	// Running agents are prioritized, so general-purpose (running) comes first
	if !strings.Contains(result, "general-purpose") {
		t.Errorf("Expected running agent 'general-purpose' to be shown, got '%s'", result)
	}
	// First completed agent (Plan) should fill the second slot
	if !strings.Contains(result, "Plan") {
		t.Errorf("Expected first completed agent 'Plan' to fill second slot, got '%s'", result)
	}
	// Should only show 2 agents total (1 running + 1 completed)
	if strings.Contains(result, "Bash") {
		t.Error("Expected Bash to be hidden (max_display=2)")
	}
	if strings.Contains(result, "Explore") {
		t.Error("Expected Explore to be hidden (max_display=2)")
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

func TestAgentsWidget_Description(t *testing.T) {
	w := &AgentsWidget{}
	w.SetTranscript(&transcript.Data{
		Agents: []transcript.Agent{
			{
				ID:          "1",
				Type:        "Explore",
				Status:      "completed",
				Description: "Analyze widget structure",
				StartTime:   1000,
				EndTime:     43000, // 42 seconds
			},
		},
	})

	result := w.Render(&input.Session{}, &config.WidgetConfig{})

	// Should show description (truncated to default max_description_len=15)
	if !strings.Contains(result, "Analyze widg") {
		t.Errorf("Expected truncated description containing 'Analyze widg', got '%s'", result)
	}
	// Should show duration
	if !strings.Contains(result, "42s") {
		t.Errorf("Expected duration '42s', got '%s'", result)
	}
}

func TestAgentsWidget_RunningShowsElapsedTime(t *testing.T) {
	// Mock nowUnixMilli to return a fixed time
	originalNow := nowUnixMilli
	nowUnixMilli = func() int64 {
		return 43000 // 42 seconds after StartTime (1000)
	}
	defer func() { nowUnixMilli = originalNow }()

	w := &AgentsWidget{}
	w.SetTranscript(&transcript.Data{
		Agents: []transcript.Agent{
			{
				ID:          "1",
				Type:        "Plan",
				Status:      "running",
				Description: "Planning implementation",
				StartTime:   1000,
			},
		},
	})

	result := w.Render(&input.Session{}, &config.WidgetConfig{})

	// Running agents should show elapsed time with "..." suffix
	if !strings.Contains(result, "42s...") {
		t.Errorf("Expected '42s...' for running agent, got '%s'", result)
	}
}

func TestAgentsWidget_HideDescription(t *testing.T) {
	w := &AgentsWidget{}
	w.SetTranscript(&transcript.Data{
		Agents: []transcript.Agent{
			{
				ID:          "1",
				Type:        "Explore",
				Status:      "completed",
				Description: "Analyze widget structure",
				StartTime:   1000,
				EndTime:     43000,
			},
		},
	})

	cfg := &config.WidgetConfig{
		Extra: map[string]string{"show_description": "false"},
	}

	result := w.Render(&input.Session{}, cfg)

	// Should NOT show description
	if strings.Contains(result, "Analyze") {
		t.Errorf("Should not show description when show_description=false, got '%s'", result)
	}
	// But should still show duration
	if !strings.Contains(result, "42s") {
		t.Errorf("Expected duration '42s', got '%s'", result)
	}
}

func TestAgentsWidget_HideDuration(t *testing.T) {
	w := &AgentsWidget{}
	w.SetTranscript(&transcript.Data{
		Agents: []transcript.Agent{
			{
				ID:          "1",
				Type:        "Explore",
				Status:      "completed",
				Description: "Analyze widget structure",
				StartTime:   1000,
				EndTime:     43000,
			},
		},
	})

	cfg := &config.WidgetConfig{
		Extra: map[string]string{"show_duration": "false"},
	}

	result := w.Render(&input.Session{}, cfg)

	// Should show description
	if !strings.Contains(result, "Analyze") {
		t.Errorf("Expected description, got '%s'", result)
	}
	// Should NOT show duration
	if strings.Contains(result, "42s") {
		t.Errorf("Should not show duration when show_duration=false, got '%s'", result)
	}
}

func TestAgentsWidget_DescriptionTruncation(t *testing.T) {
	w := &AgentsWidget{}
	w.SetTranscript(&transcript.Data{
		Agents: []transcript.Agent{
			{
				ID:          "1",
				Type:        "Explore",
				Status:      "completed",
				Description: "This is a very long description that should be truncated",
				StartTime:   1000,
				EndTime:     2000,
			},
		},
	})

	cfg := &config.WidgetConfig{
		Extra: map[string]string{"max_description_len": "10"},
	}

	result := w.Render(&input.Session{}, cfg)

	// Description should be truncated with "..."
	if strings.Contains(result, "truncated") {
		t.Errorf("Description should be truncated, got '%s'", result)
	}
	if !strings.Contains(result, "...") {
		t.Errorf("Expected '...' for truncation, got '%s'", result)
	}
}

func TestAgentsWidget_InternalSeparator(t *testing.T) {
	w := &AgentsWidget{}
	w.SetTranscript(&transcript.Data{
		Agents: []transcript.Agent{
			{ID: "1", Type: "Plan", Status: "completed"},
			{ID: "2", Type: "Explore", Status: "completed"},
		},
	})

	cfg := &config.WidgetConfig{
		Extra: map[string]string{"show_description": "false", "show_duration": "false"},
	}

	result := w.Render(&input.Session{}, cfg)

	// Should use " · " separator (not " | " to avoid confusion with visor widget separator)
	if !strings.Contains(result, " · ") {
		t.Errorf("Expected ' · ' separator, got '%s'", result)
	}
}

func TestAgentsWidget_RunningPriority(t *testing.T) {
	w := &AgentsWidget{}
	w.SetTranscript(&transcript.Data{
		Agents: []transcript.Agent{
			{ID: "1", Type: "Plan", Status: "completed"},
			{ID: "2", Type: "Explore", Status: "completed"},
			{ID: "3", Type: "Bash", Status: "running"},
		},
	})

	cfg := &config.WidgetConfig{
		Extra: map[string]string{
			"max_display":      "0", // unlimited
			"show_description": "false",
			"show_duration":    "false",
		},
	}

	result := w.Render(&input.Session{}, cfg)

	// Running agent (Bash) should appear before completed agents
	bashIdx := strings.Index(result, "Bash")
	planIdx := strings.Index(result, "Plan")
	exploreIdx := strings.Index(result, "Explore")

	if bashIdx < 0 || planIdx < 0 || exploreIdx < 0 {
		t.Fatalf("Expected all agents in output, got '%s'", result)
	}
	if bashIdx > planIdx || bashIdx > exploreIdx {
		t.Errorf("Expected running agent (Bash) before completed agents, got '%s'", result)
	}
}

func TestAgentsWidget_RunningPriorityWithMaxDisplay(t *testing.T) {
	w := &AgentsWidget{}
	w.SetTranscript(&transcript.Data{
		Agents: []transcript.Agent{
			{ID: "1", Type: "Plan", Status: "completed"},
			{ID: "2", Type: "Explore", Status: "completed"},
			{ID: "3", Type: "Bash", Status: "running"},
			{ID: "4", Type: "general-purpose", Status: "completed"},
		},
	})

	cfg := &config.WidgetConfig{
		Extra: map[string]string{
			"max_display":      "2",
			"show_description": "false",
			"show_duration":    "false",
		},
	}

	result := w.Render(&input.Session{}, cfg)

	// Running agent should always be included even with max_display
	if !strings.Contains(result, "Bash") {
		t.Errorf("Expected running agent 'Bash' to be included, got '%s'", result)
	}
	// First completed (Plan) fills the remaining slot
	if !strings.Contains(result, "Plan") {
		t.Errorf("Expected first completed agent 'Plan' to be shown, got '%s'", result)
	}
	// Other completed agents should be hidden
	if strings.Contains(result, "Explore") {
		t.Error("Expected Explore to be hidden (max_display=2)")
	}
	if strings.Contains(result, "general-purpose") {
		t.Error("Expected general-purpose to be hidden (max_display=2)")
	}
}

func TestAgentsWidget_MultipleRunningExceedMaxDisplay(t *testing.T) {
	w := &AgentsWidget{}
	w.SetTranscript(&transcript.Data{
		Agents: []transcript.Agent{
			{ID: "1", Type: "Plan", Status: "completed"},
			{ID: "2", Type: "Bash", Status: "running"},
			{ID: "3", Type: "Explore", Status: "running"},
			{ID: "4", Type: "general-purpose", Status: "running"},
		},
	})

	cfg := &config.WidgetConfig{
		Extra: map[string]string{
			"max_display":      "2",
			"show_description": "false",
			"show_duration":    "false",
		},
	}

	result := w.Render(&input.Session{}, cfg)

	// First 2 running agents should be shown (Bash, Explore)
	if !strings.Contains(result, "Bash") {
		t.Errorf("Expected first running agent 'Bash', got '%s'", result)
	}
	if !strings.Contains(result, "Explore") {
		t.Errorf("Expected second running agent 'Explore', got '%s'", result)
	}
	// Third running agent and completed should be hidden
	if strings.Contains(result, "general-purpose") {
		t.Error("Expected general-purpose to be hidden (max_display=2)")
	}
	if strings.Contains(result, "Plan") {
		t.Error("Expected Plan to be hidden (max_display=2)")
	}
}

func TestFormatDurationSec(t *testing.T) {
	tests := []struct {
		seconds  int64
		expected string
	}{
		{0, "0s"},
		{42, "42s"},
		{59, "59s"},
		{60, "1m"},
		{90, "1m"},
		{120, "2m"},
		{3600, "1h"},
		{3660, "1h1m"},
		{7200, "2h"},
		{7320, "2h2m"},
	}

	for _, tt := range tests {
		result := formatDurationSec(tt.seconds)
		if result != tt.expected {
			t.Errorf("formatDurationSec(%d) = %s, want %s", tt.seconds, result, tt.expected)
		}
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is longer", 10, "this is..."},
		{"abc", 3, "abc"},      // Fits exactly, no truncation
		{"abcd", 3, "..."},     // maxLen <= 3 returns "..."
		{"abcdef", 5, "ab..."}, // Truncates normally
		// Multibyte character tests (rune-based truncation)
		{"한글테스트", 5, "한글테스트"},      // 5 runes, fits exactly
		{"한글테스트입니다", 7, "한글테스..."},  // 8 runes, truncate to 7
		{"분석하기", 10, "분석하기"},        // 4 runes, fits in 10
		{"코드리뷰중", 5, "코드리뷰중"},       // 5 runes, fits exactly
	}

	for _, tt := range tests {
		result := truncateString(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("truncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
		}
	}
}
