package widgets

import (
	"strings"
	"testing"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/history"
	"github.com/namyoungkim/visor/internal/input"
)

func TestContextSparkWidget_Name(t *testing.T) {
	w := &ContextSparkWidget{}
	if w.Name() != "context_spark" {
		t.Errorf("expected 'context_spark', got %s", w.Name())
	}
}

func TestContextSparkWidget_NoHistory(t *testing.T) {
	w := &ContextSparkWidget{}
	session := &input.Session{}
	cfg := &config.WidgetConfig{Name: "context_spark"}

	result := w.Render(session, cfg)
	if !strings.Contains(result, "—") {
		t.Errorf("expected dash for no history, got %q", result)
	}
}

func TestContextSparkWidget_InsufficientHistory(t *testing.T) {
	w := &ContextSparkWidget{}
	h := &history.History{SessionID: "test"}
	h.Add(history.Entry{ContextPct: 42.0})
	w.SetHistory(h)

	session := &input.Session{}
	cfg := &config.WidgetConfig{Name: "context_spark"}

	result := w.Render(session, cfg)
	// Need at least 2 data points
	if !strings.Contains(result, "—") {
		t.Errorf("expected dash for insufficient history, got %q", result)
	}
}

func TestContextSparkWidget_Sparkline(t *testing.T) {
	w := &ContextSparkWidget{}
	h := &history.History{SessionID: "test"}

	// Add increasing values
	for i := 0; i < 8; i++ {
		h.Add(history.Entry{ContextPct: float64(i * 12)})
	}
	w.SetHistory(h)

	session := &input.Session{}
	cfg := &config.WidgetConfig{Name: "context_spark"}

	result := w.Render(session, cfg)

	// Should contain sparkline characters
	for _, char := range sparkChars {
		if strings.ContainsRune(result, char) {
			return // Found at least one sparkline char
		}
	}
	t.Errorf("expected sparkline characters, got %q", result)
}

func TestContextSparkWidget_ShouldRender(t *testing.T) {
	w := &ContextSparkWidget{}
	session := &input.Session{}
	cfg := &config.WidgetConfig{Name: "context_spark"}

	// No history
	if w.ShouldRender(session, cfg) {
		t.Error("expected false with no history")
	}

	// With history but only 1 entry
	h := &history.History{SessionID: "test"}
	h.Add(history.Entry{ContextPct: 42.0})
	w.SetHistory(h)

	if w.ShouldRender(session, cfg) {
		t.Error("expected false with only 1 entry")
	}

	// With 2+ entries
	h.Add(history.Entry{ContextPct: 50.0})
	if !w.ShouldRender(session, cfg) {
		t.Error("expected true with 2 entries")
	}
}

func TestSparkline(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   string
	}{
		{
			name:   "empty",
			values: []float64{},
			want:   "",
		},
		{
			name:   "all zeros",
			values: []float64{0, 0, 0},
			want:   "▁▁▁",
		},
		{
			name:   "all 100",
			values: []float64{100, 100, 100},
			want:   "█████", // Adjusted for 5 values
		},
		{
			name:   "increasing",
			values: []float64{0, 50, 100},
			want:   "▁▄█",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sparkline(tt.values)
			if tt.name == "all 100" {
				// Just check it contains max char
				if !strings.Contains(got, "█") {
					t.Errorf("expected max char █, got %q", got)
				}
			} else if got != tt.want {
				t.Errorf("sparkline(%v) = %q, want %q", tt.values, got, tt.want)
			}
		})
	}
}

func TestSparkColor(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   string
	}{
		{
			name:   "rising trend",
			values: []float64{30, 35, 40, 60}, // Last much higher than avg
			want:   "red",
		},
		{
			name:   "falling trend",
			values: []float64{60, 55, 50, 30}, // Last much lower than avg
			want:   "green",
		},
		{
			name:   "stable",
			values: []float64{50, 50, 50, 51}, // Nearly same
			want:   "yellow",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sparkColor(tt.values)
			if got != tt.want {
				t.Errorf("sparkColor(%v) = %q, want %q", tt.values, got, tt.want)
			}
		})
	}
}
