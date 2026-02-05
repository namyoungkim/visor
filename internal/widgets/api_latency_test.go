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

func TestAPILatencyWidget_ZeroCalls(t *testing.T) {
	w := &APILatencyWidget{}
	session := &input.Session{
		Cost: input.Cost{
			TotalAPIDurationMs: 5000,
			TotalAPICalls:      0,
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "—") {
		t.Errorf("Expected dash for zero calls, got '%s'", result)
	}
}

func TestAPILatencyWidget_PerCallMilliseconds(t *testing.T) {
	w := &APILatencyWidget{}
	session := &input.Session{
		Cost: input.Cost{
			TotalAPIDurationMs: 2500,
			TotalAPICalls:      5, // 500ms per call
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "500ms") {
		t.Errorf("Expected 500ms per call, got '%s'", result)
	}
}

func TestAPILatencyWidget_PerCallSeconds(t *testing.T) {
	w := &APILatencyWidget{}
	session := &input.Session{
		Cost: input.Cost{
			TotalAPIDurationMs: 3000,
			TotalAPICalls:      2, // 1500ms = 1.5s per call
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	if !strings.Contains(result, "1.5s") {
		t.Errorf("Expected 1.5s per call, got '%s'", result)
	}
}

func TestAPILatencyWidget_HighLatency(t *testing.T) {
	w := &APILatencyWidget{}
	session := &input.Session{
		Cost: input.Cost{
			TotalAPIDurationMs: 6000,
			TotalAPICalls:      1, // 6000ms per call
		},
	}

	result := w.Render(session, &config.WidgetConfig{})

	// Should be red for high latency
	if !strings.Contains(result, "\033[31m") {
		t.Errorf("Expected red color for high latency")
	}
}
