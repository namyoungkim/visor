package widgets

import (
	"os"
	"strings"
	"testing"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
)

func TestCWDWidget_Empty(t *testing.T) {
	w := &CWDWidget{}
	session := &input.Session{}

	result := w.Render(session, &config.WidgetConfig{})
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestCWDWidget_ShouldRender(t *testing.T) {
	w := &CWDWidget{}

	tests := []struct {
		name     string
		session  *input.Session
		expected bool
	}{
		{"with cwd", &input.Session{CWD: "/home/user/project"}, true},
		{"empty cwd", &input.Session{}, false},
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

func TestCWDWidget_HomeAbbreviation(t *testing.T) {
	w := &CWDWidget{}
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get home directory")
	}

	session := &input.Session{CWD: home + "/project/visor"}
	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "~/project/visor") {
		t.Errorf("Expected home abbreviation with '~/project/visor', got '%s'", result)
	}
}

func TestCWDWidget_ShowBasename(t *testing.T) {
	w := &CWDWidget{}
	session := &input.Session{CWD: "/home/user/project/visor"}
	cfg := &config.WidgetConfig{
		Extra: map[string]string{"show_basename": "true"},
	}

	result := w.Render(session, cfg)

	if !strings.Contains(result, "visor") {
		t.Errorf("Expected 'visor', got '%s'", result)
	}
	if strings.Contains(result, "/") {
		t.Errorf("Expected basename only (no '/'), got '%s'", result)
	}
}

func TestCWDWidget_ShowLabel(t *testing.T) {
	w := &CWDWidget{}
	session := &input.Session{CWD: "/tmp/test"}
	cfg := &config.WidgetConfig{
		Extra: map[string]string{"show_label": "true"},
	}

	result := w.Render(session, cfg)

	if !strings.Contains(result, "CWD:") {
		t.Errorf("Expected 'CWD:' prefix, got '%s'", result)
	}
}

func TestCWDWidget_MaxLength(t *testing.T) {
	w := &CWDWidget{}
	session := &input.Session{CWD: "/very/long/path/to/some/deep/directory"}
	cfg := &config.WidgetConfig{
		Extra: map[string]string{"max_length": "15"},
	}

	result := w.Render(session, cfg)

	// Should contain truncation indicator
	if !strings.Contains(result, "…") {
		t.Errorf("Expected truncation with '…', got '%s'", result)
	}
}

func TestCWDWidget_NonHomePath(t *testing.T) {
	w := &CWDWidget{}
	session := &input.Session{CWD: "/tmp/test"}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "/tmp/test") {
		t.Errorf("Expected '/tmp/test', got '%s'", result)
	}
}

func TestAbbreviateHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get home directory")
	}

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"home dir itself", home, "~"},
		{"under home", home + "/project", "~/project"},
		{"not under home", "/tmp/test", "/tmp/test"},
		{"home prefix but not dir", home + "extra", home + "extra"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := abbreviateHome(tt.path)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestTruncatePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		maxLen   int
		contains string
	}{
		{"short enough", "/tmp", 10, "/tmp"},
		{"needs truncation", "/very/long/path/to/dir", 10, "…"},
		{"very short max", "/a/b/c", 2, "/a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncatePath(tt.path, tt.maxLen)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("Expected result to contain '%s', got '%s'", tt.contains, result)
			}
		})
	}
}
