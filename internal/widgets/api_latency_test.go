package widgets

import (
	"strings"
	"testing"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
)

func TestAPILatencyWidget_ZeroLatency(t *testing.T) {
	w := &APILatencyWidget{}
	session := &input.Session{}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "—") {
		t.Errorf("Expected dash for zero latency, got '%s'", result)
	}
}

func TestAPILatencyWidget_Milliseconds(t *testing.T) {
	w := &APILatencyWidget{}
	session := &input.Session{
		Cost: input.Cost{
			TotalAPIDurationMs: 500,
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "500ms") {
		t.Errorf("Expected 500ms, got '%s'", result)
	}
}

func TestAPILatencyWidget_Seconds(t *testing.T) {
	w := &APILatencyWidget{}
	session := &input.Session{
		Cost: input.Cost{
			TotalAPIDurationMs: 1500,
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "1.5s") {
		t.Errorf("Expected 1.5s, got '%s'", result)
	}
}

func TestAPILatencyWidget_HighLatency(t *testing.T) {
	w := &APILatencyWidget{}
	session := &input.Session{
		Cost: input.Cost{
			TotalAPIDurationMs: 6000,
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	// Should be red for high latency
	if !strings.Contains(result, "\033[31m") {
		t.Errorf("Expected red color for high latency")
	}
}
