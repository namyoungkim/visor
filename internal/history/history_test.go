package history

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHistory_AddAndCount(t *testing.T) {
	h := &History{SessionID: "test"}

	if h.Count() != 0 {
		t.Errorf("expected 0 entries, got %d", h.Count())
	}

	h.Add(Entry{ContextPct: 10.0})
	if h.Count() != 1 {
		t.Errorf("expected 1 entry, got %d", h.Count())
	}

	h.Add(Entry{ContextPct: 20.0})
	if h.Count() != 2 {
		t.Errorf("expected 2 entries, got %d", h.Count())
	}
}

func TestHistory_MaxEntries(t *testing.T) {
	h := &History{SessionID: "test"}

	// Add more than MaxEntries
	for i := 0; i < MaxEntries+10; i++ {
		h.Add(Entry{ContextPct: float64(i)})
	}

	if h.Count() != MaxEntries {
		t.Errorf("expected %d entries (max), got %d", MaxEntries, h.Count())
	}

	// Oldest entries should be trimmed
	first := h.Entries[0].ContextPct
	if first != 10.0 {
		t.Errorf("expected first entry to be 10.0, got %.1f", first)
	}
}

func TestHistory_GetContextHistory(t *testing.T) {
	h := &History{SessionID: "test"}

	h.Add(Entry{ContextPct: 10.0})
	h.Add(Entry{ContextPct: 20.0})
	h.Add(Entry{ContextPct: 30.0})
	h.Add(Entry{ContextPct: 40.0})
	h.Add(Entry{ContextPct: 50.0})

	// Get last 3
	values := h.GetContextHistory(3)
	if len(values) != 3 {
		t.Errorf("expected 3 values, got %d", len(values))
	}
	if values[0] != 30.0 || values[1] != 40.0 || values[2] != 50.0 {
		t.Errorf("expected [30, 40, 50], got %v", values)
	}

	// Get more than available
	values = h.GetContextHistory(10)
	if len(values) != 5 {
		t.Errorf("expected 5 values (all), got %d", len(values))
	}
}

func TestHistory_Latest(t *testing.T) {
	h := &History{SessionID: "test"}

	if h.Latest() != nil {
		t.Error("expected nil for empty history")
	}

	h.Add(Entry{ContextPct: 10.0})
	h.Add(Entry{ContextPct: 20.0})

	latest := h.Latest()
	if latest == nil {
		t.Fatal("expected non-nil latest")
	}
	if latest.ContextPct != 20.0 {
		t.Errorf("expected 20.0, got %.1f", latest.ContextPct)
	}
}

func TestHistory_SaveAndLoad(t *testing.T) {
	// Use temp directory
	tmpDir := t.TempDir()
	sessionID := "test_save_load"

	// Override history dir for testing
	origDir := HistoryDirFunc
	HistoryDirFunc = func() string { return tmpDir }
	defer func() { HistoryDirFunc = origDir }()

	// Create and save history
	h := &History{SessionID: sessionID}
	h.Add(Entry{ContextPct: 42.5, CostUSD: 0.48})
	h.Add(Entry{ContextPct: 55.0, CostUSD: 0.75})

	if err := h.Save(); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	// Check file exists
	path := filepath.Join(tmpDir, "history_"+sessionID+".json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("history file not created")
	}

	// Load and verify
	loaded, err := Load(sessionID)
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	if loaded.Count() != 2 {
		t.Errorf("expected 2 entries, got %d", loaded.Count())
	}

	if loaded.Entries[0].ContextPct != 42.5 {
		t.Errorf("expected 42.5, got %.1f", loaded.Entries[0].ContextPct)
	}
}

func TestHistory_LoadNonexistent(t *testing.T) {
	tmpDir := t.TempDir()
	origDir := HistoryDirFunc
	HistoryDirFunc = func() string { return tmpDir }
	defer func() { HistoryDirFunc = origDir }()

	h, err := Load("nonexistent_session")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if h.Count() != 0 {
		t.Errorf("expected empty history, got %d entries", h.Count())
	}
}

func TestSanitizeSessionID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "default"},
		{"normal_session-123", "normal_session-123"},
		{"../../../etc/passwd", "_________etc_passwd"},
		{"session/with/slashes", "session_with_slashes"},
		{"session with spaces", "session_with_spaces"},
		{"a@b#c$d%e", "a_b_c_d_e"},
	}

	for _, tt := range tests {
		result := sanitizeSessionID(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeSessionID(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestSanitizeSessionID_LongInput(t *testing.T) {
	// Create a long session ID
	longID := ""
	for i := 0; i < 100; i++ {
		longID += "a"
	}

	result := sanitizeSessionID(longID)
	if len(result) > 64 {
		t.Errorf("expected max length 64, got %d", len(result))
	}
}
