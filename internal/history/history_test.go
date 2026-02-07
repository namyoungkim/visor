package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestHistory_UpdateBlockStartTime(t *testing.T) {
	h := &History{SessionID: "test"}

	// First call should set the block start time
	h.UpdateBlockStartTime()
	if h.BlockStartTime == 0 {
		t.Error("expected BlockStartTime to be set")
	}

	firstTime := h.BlockStartTime

	// Immediate second call should not change it (block not expired)
	h.UpdateBlockStartTime()
	if h.BlockStartTime != firstTime {
		t.Error("expected BlockStartTime to remain unchanged")
	}

	// Simulate expired block (set start time to 6 hours ago)
	expiredTime := time.Now().UnixMilli() - (6 * 60 * 60 * 1000)
	h.BlockStartTime = expiredTime
	h.UpdateBlockStartTime()

	// After update, block start time should be recent (not the expired time)
	if h.BlockStartTime == expiredTime {
		t.Error("expected BlockStartTime to be updated after expiry")
	}

	// New block start time should be recent (within last second)
	if time.Now().UnixMilli()-h.BlockStartTime > 1000 {
		t.Error("expected BlockStartTime to be set to current time")
	}
}

func TestHistory_GetBlockRemainingMs(t *testing.T) {
	h := &History{SessionID: "test"}

	// No block start time
	if remaining := h.GetBlockRemainingMs(); remaining != 0 {
		t.Errorf("expected 0 for no block start, got %d", remaining)
	}

	// Just started block (should be close to 5 hours)
	h.BlockStartTime = time.Now().UnixMilli()
	remaining := h.GetBlockRemainingMs()

	// Should be within 1 second of full duration
	expectedMin := int64(BlockDurationMs - 1000)
	if remaining < expectedMin {
		t.Errorf("expected remaining >= %d, got %d", expectedMin, remaining)
	}

	// 2 hours elapsed (3 hours remaining)
	h.BlockStartTime = time.Now().UnixMilli() - (2 * 60 * 60 * 1000)
	remaining = h.GetBlockRemainingMs()

	// Should be around 3 hours (±1 second tolerance)
	expected := int64(3 * 60 * 60 * 1000)
	tolerance := int64(1000)
	if remaining < expected-tolerance || remaining > expected+tolerance {
		t.Errorf("expected ~%d, got %d", expected, remaining)
	}

	// Expired block
	h.BlockStartTime = time.Now().UnixMilli() - (6 * 60 * 60 * 1000)
	if remaining := h.GetBlockRemainingMs(); remaining != 0 {
		t.Errorf("expected 0 for expired block, got %d", remaining)
	}
}

func TestHistory_GetBlockElapsedPct(t *testing.T) {
	h := &History{SessionID: "test"}

	// No block start time
	if pct := h.GetBlockElapsedPct(); pct != 0 {
		t.Errorf("expected 0 for no block start, got %.1f", pct)
	}

	// Just started (0% elapsed)
	h.BlockStartTime = time.Now().UnixMilli()
	pct := h.GetBlockElapsedPct()
	if pct > 1.0 { // Allow small tolerance for test execution time
		t.Errorf("expected ~0%% elapsed, got %.1f%%", pct)
	}

	// 50% elapsed (2.5 hours)
	h.BlockStartTime = time.Now().UnixMilli() - (int64(BlockDurationMs) / 2)
	pct = h.GetBlockElapsedPct()
	if pct < 49.0 || pct > 51.0 {
		t.Errorf("expected ~50%% elapsed, got %.1f%%", pct)
	}

	// Fully expired (100%)
	h.BlockStartTime = time.Now().UnixMilli() - (6 * 60 * 60 * 1000)
	pct = h.GetBlockElapsedPct()
	if pct != 100.0 {
		t.Errorf("expected 100%% for expired block, got %.1f%%", pct)
	}
}

func TestLoadGlobalBlockStart_NoFile(t *testing.T) {
	tmpDir := t.TempDir()
	origDir := HistoryDirFunc
	HistoryDirFunc = func() string { return tmpDir }
	defer func() { HistoryDirFunc = origDir }()

	ts := LoadGlobalBlockStart()
	if ts != 0 {
		t.Errorf("expected 0 for missing file, got %d", ts)
	}
}

func TestSaveAndLoadGlobalBlockStart(t *testing.T) {
	tmpDir := t.TempDir()
	origDir := HistoryDirFunc
	HistoryDirFunc = func() string { return tmpDir }
	defer func() { HistoryDirFunc = origDir }()

	want := time.Now().UnixMilli()
	if err := SaveGlobalBlockStart(want); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	got := LoadGlobalBlockStart()
	if got != want {
		t.Errorf("expected %d, got %d", want, got)
	}
}

func TestLoadGlobalBlockStart_CorruptedJSON(t *testing.T) {
	tmpDir := t.TempDir()
	origDir := HistoryDirFunc
	HistoryDirFunc = func() string { return tmpDir }
	defer func() { HistoryDirFunc = origDir }()

	// Write invalid JSON to the block state file
	path := filepath.Join(tmpDir, "block_state.json")
	if err := os.WriteFile(path, []byte("{invalid json}"), 0600); err != nil {
		t.Fatalf("failed to write corrupted file: %v", err)
	}

	ts := LoadGlobalBlockStart()
	if ts != 0 {
		t.Errorf("expected 0 for corrupted file, got %d", ts)
	}
}

func TestGlobalBlockStart_CrossSession(t *testing.T) {
	tmpDir := t.TempDir()
	origDir := HistoryDirFunc
	HistoryDirFunc = func() string { return tmpDir }
	defer func() { HistoryDirFunc = origDir }()

	// Session A sets a block start time
	histA := &History{SessionID: "session_a"}
	histA.UpdateBlockStartTime()
	blockStart := histA.BlockStartTime

	// Save globally (as main.go does)
	if err := SaveGlobalBlockStart(blockStart); err != nil {
		t.Fatalf("failed to save global block start: %v", err)
	}

	// Session B is brand new — no block start yet
	histB, err := Load("session_b")
	if err != nil {
		t.Fatalf("failed to load session_b: %v", err)
	}
	if histB.BlockStartTime != 0 {
		t.Fatal("expected new session to have zero BlockStartTime")
	}

	// Inherit from global (as main.go does)
	histB.BlockStartTime = LoadGlobalBlockStart()
	histB.UpdateBlockStartTime()

	if histB.BlockStartTime != blockStart {
		t.Errorf("session B should inherit session A's block start: want %d, got %d",
			blockStart, histB.BlockStartTime)
	}
}

func TestHistory_BlockStartTime_Persistence(t *testing.T) {
	tmpDir := t.TempDir()
	sessionID := "test_block_persistence"

	origDir := HistoryDirFunc
	HistoryDirFunc = func() string { return tmpDir }
	defer func() { HistoryDirFunc = origDir }()

	// Create history with block start time
	h := &History{SessionID: sessionID}
	h.UpdateBlockStartTime()
	originalBlockStart := h.BlockStartTime

	if err := h.Save(); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	// Load and verify block start time persisted
	loaded, err := Load(sessionID)
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	if loaded.BlockStartTime != originalBlockStart {
		t.Errorf("expected BlockStartTime %d, got %d", originalBlockStart, loaded.BlockStartTime)
	}
}
