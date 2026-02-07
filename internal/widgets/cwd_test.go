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

	if !strings.Contains(result, "…") {
		t.Errorf("Expected truncation with '…', got '%s'", result)
	}
}

func TestCWDWidget_AutoTruncate(t *testing.T) {
	w := &CWDWidget{}
	longPath := "/very/very/very/long/path/that/should/be/auto/truncated"
	session := &input.Session{CWD: longPath}
	cfg := &config.WidgetConfig{} // no max_length set

	// Set terminal width to 90 → auto max = 30
	t.Setenv("COLUMNS", "90")

	result := w.Render(session, cfg)
	if !strings.Contains(result, "…") {
		t.Errorf("Expected auto-truncation with 90-col terminal, got '%s'", result)
	}
}

func TestCWDWidget_AutoTruncateShortPath(t *testing.T) {
	w := &CWDWidget{}
	session := &input.Session{CWD: "/tmp/test"}
	cfg := &config.WidgetConfig{} // no max_length set

	// Set terminal width to 120 → auto max = 40, path is shorter
	t.Setenv("COLUMNS", "120")

	result := w.Render(session, cfg)
	if strings.Contains(result, "…") {
		t.Errorf("Short path should not be truncated, got '%s'", result)
	}
	if !strings.Contains(result, "/tmp/test") {
		t.Errorf("Expected full path '/tmp/test', got '%s'", result)
	}
}

func TestCWDWidget_ExplicitMaxLengthOverridesAuto(t *testing.T) {
	w := &CWDWidget{}
	session := &input.Session{CWD: "/very/very/very/long/path/to/project"}
	cfg := &config.WidgetConfig{
		Extra: map[string]string{"max_length": "10"},
	}

	// Auto would be 40, but explicit 10 should win
	t.Setenv("COLUMNS", "120")

	result := w.Render(session, cfg)
	if !strings.Contains(result, "…") {
		t.Errorf("Expected truncation with explicit max_length=10, got '%s'", result)
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

func TestCWDWidget_ShowBasenameWithMaxLength(t *testing.T) {
	w := &CWDWidget{}
	session := &input.Session{CWD: "/home/user/very-long-directory-name-here"}
	cfg := &config.WidgetConfig{
		Extra: map[string]string{"show_basename": "true", "max_length": "10"},
	}

	result := w.Render(session, cfg)

	if !strings.Contains(result, "…") {
		t.Errorf("Expected truncation for long basename, got '%s'", result)
	}
}

func TestCWDWidget_ShowLabelWithMaxLength(t *testing.T) {
	w := &CWDWidget{}
	session := &input.Session{CWD: "/very/long/path/to/directory"}
	cfg := &config.WidgetConfig{
		Extra: map[string]string{"show_label": "true", "max_length": "10"},
	}

	result := w.Render(session, cfg)

	// max_length applies to path only, label is prepended after
	if !strings.Contains(result, "CWD:") {
		t.Errorf("Expected 'CWD:' prefix, got '%s'", result)
	}
	if !strings.Contains(result, "…") {
		t.Errorf("Expected truncation, got '%s'", result)
	}
}

func TestCWDWidget_NonASCIIPath(t *testing.T) {
	w := &CWDWidget{}
	session := &input.Session{CWD: "/home/user/프로젝트/visor"}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "프로젝트") {
		t.Errorf("Expected Korean characters preserved, got '%s'", result)
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
		expected string
	}{
		{"short enough", "/tmp", 10, "/tmp"},
		{"exact length", "/tmp/a", 6, "/tmp/a"},
		{"truncate with slash boundary", "/very/long/path/to/dir", 10, "…/to/dir"},
		{"truncate no slash boundary", "abcdefghij", 5, "…ghij"},
		{"maxLen 1", "/a/b/c", 1, "/"},
		{"maxLen 2", "/a/b/c", 2, "…c"},
		{"maxLen 3", "/a/b/c", 3, "…/c"},
		{"non-ascii path", "/home/프로젝트/visor", 8, "…/visor"},
		{"non-ascii truncate boundary", "/프/로/젝/트", 5, "…/젝/트"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncatePath(tt.path, tt.maxLen)
			runes := []rune(result)
			if len(runes) > tt.maxLen {
				t.Errorf("Result rune length %d exceeds maxLen %d: '%s'", len(runes), tt.maxLen, result)
			}
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
