package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// MaxEntries is the maximum number of history entries to keep per session.
const MaxEntries = 20

// Entry represents a single history entry.
type Entry struct {
	Timestamp      int64   `json:"ts"`
	ContextPct     float64 `json:"ctx_pct"`
	CostUSD        float64 `json:"cost"`
	DurationMs     int64   `json:"dur_ms"`
	CacheHitPct    float64 `json:"cache_pct"`
	APILatencyMs   int64   `json:"api_ms"`
}

// History manages session history data.
type History struct {
	SessionID string  `json:"session_id"`
	Entries   []Entry `json:"entries"`
}

// HistoryDirFunc is the function used to get the history directory.
// Can be overridden in tests.
var HistoryDirFunc = defaultHistoryDir

// defaultHistoryDir returns the default directory for storing history files.
func defaultHistoryDir() string {
	// Use ~/.cache/visor for history
	home, err := os.UserHomeDir()
	if err != nil {
		return "/tmp"
	}
	return filepath.Join(home, ".cache", "visor")
}

// historyPath returns the path for a session's history file.
func historyPath(sessionID string) string {
	if sessionID == "" {
		sessionID = "default"
	}
	return filepath.Join(HistoryDirFunc(), "history_"+sessionID+".json")
}

// Load loads history for a session from disk.
func Load(sessionID string) (*History, error) {
	path := historyPath(sessionID)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &History{SessionID: sessionID, Entries: []Entry{}}, nil
		}
		return nil, err
	}

	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		// Corrupted file, start fresh
		return &History{SessionID: sessionID, Entries: []Entry{}}, nil
	}

	return &h, nil
}

// Save writes history to disk.
func (h *History) Save() error {
	dir := HistoryDirFunc()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Trim to max entries
	if len(h.Entries) > MaxEntries {
		h.Entries = h.Entries[len(h.Entries)-MaxEntries:]
	}

	data, err := json.Marshal(h)
	if err != nil {
		return err
	}

	return os.WriteFile(historyPath(h.SessionID), data, 0644)
}

// Add adds a new entry to the history.
func (h *History) Add(entry Entry) {
	if entry.Timestamp == 0 {
		entry.Timestamp = time.Now().Unix()
	}
	h.Entries = append(h.Entries, entry)

	// Trim in memory
	if len(h.Entries) > MaxEntries {
		h.Entries = h.Entries[len(h.Entries)-MaxEntries:]
	}
}

// GetContextHistory returns the last n context percentage values.
func (h *History) GetContextHistory(n int) []float64 {
	if n <= 0 || len(h.Entries) == 0 {
		return nil
	}

	start := 0
	if len(h.Entries) > n {
		start = len(h.Entries) - n
	}

	result := make([]float64, len(h.Entries)-start)
	for i, e := range h.Entries[start:] {
		result[i] = e.ContextPct
	}
	return result
}

// GetCostHistory returns the last n cost values.
func (h *History) GetCostHistory(n int) []float64 {
	if n <= 0 || len(h.Entries) == 0 {
		return nil
	}

	start := 0
	if len(h.Entries) > n {
		start = len(h.Entries) - n
	}

	result := make([]float64, len(h.Entries)-start)
	for i, e := range h.Entries[start:] {
		result[i] = e.CostUSD
	}
	return result
}

// Latest returns the most recent entry, or nil if none.
func (h *History) Latest() *Entry {
	if len(h.Entries) == 0 {
		return nil
	}
	return &h.Entries[len(h.Entries)-1]
}

// Count returns the number of entries.
func (h *History) Count() int {
	return len(h.Entries)
}
