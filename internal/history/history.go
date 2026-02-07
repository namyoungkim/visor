package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
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
	SessionID      string  `json:"session_id"`
	Entries        []Entry `json:"entries"`
	BlockStartTime int64   `json:"block_start_ts,omitempty"`
}

// BlockDurationMs is the duration of a Claude Pro rate limit block (5 hours).
const BlockDurationMs = 5 * 60 * 60 * 1000 // 5 hours in milliseconds

// UpdateBlockStartTime sets the block start time if not already set or if expired.
func (h *History) UpdateBlockStartTime() {
	now := time.Now().UnixMilli()

	// If no block start time or block has expired, start a new block
	if h.BlockStartTime == 0 || now-h.BlockStartTime >= BlockDurationMs {
		h.BlockStartTime = now
	}
}

// GetBlockRemainingMs returns the remaining milliseconds in the current block.
// Returns 0 if block has expired.
func (h *History) GetBlockRemainingMs() int64 {
	if h.BlockStartTime == 0 {
		return 0
	}

	now := time.Now().UnixMilli()
	elapsed := now - h.BlockStartTime
	remaining := BlockDurationMs - elapsed

	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetBlockElapsedPct returns the percentage of the block that has elapsed.
func (h *History) GetBlockElapsedPct() float64 {
	if h.BlockStartTime == 0 {
		return 0
	}

	now := time.Now().UnixMilli()
	elapsed := now - h.BlockStartTime

	if elapsed >= BlockDurationMs {
		return 100.0
	}

	return float64(elapsed) / float64(BlockDurationMs) * 100
}

// GetBlockStartTime returns the block start time as time.Time.
// Returns zero time if no block has been started.
func (h *History) GetBlockStartTime() time.Time {
	if h.BlockStartTime == 0 {
		return time.Time{}
	}
	return time.UnixMilli(h.BlockStartTime)
}

// globalBlockState is the structure for the global block state file.
type globalBlockState struct {
	BlockStartTs int64 `json:"block_start_ts"`
}

// globalBlockStatePath returns the path to the global block state file.
func globalBlockStatePath() string {
	return filepath.Join(HistoryDirFunc(), "block_state.json")
}

// LoadGlobalBlockStart reads the global block start timestamp from disk.
// Returns 0 if the file doesn't exist or is unreadable.
func LoadGlobalBlockStart() int64 {
	data, err := os.ReadFile(globalBlockStatePath())
	if err != nil {
		return 0
	}
	var state globalBlockState
	if err := json.Unmarshal(data, &state); err != nil {
		return 0
	}
	return state.BlockStartTs
}

// SaveGlobalBlockStart persists the block start timestamp globally.
// Uses atomic rename to prevent corruption from concurrent writes.
func SaveGlobalBlockStart(ts int64) error {
	dir := HistoryDirFunc()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.Marshal(globalBlockState{BlockStartTs: ts})
	if err != nil {
		return err
	}
	tmpPath := globalBlockStatePath() + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmpPath, globalBlockStatePath())
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

// sanitizeSessionID removes potentially dangerous characters from session ID.
// Only allows alphanumeric, dash, and underscore characters.
var sessionIDRegex = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

func sanitizeSessionID(sessionID string) string {
	if sessionID == "" {
		return "default"
	}
	// Replace invalid characters with underscore
	safe := sessionIDRegex.ReplaceAllString(sessionID, "_")
	// Limit length to prevent overly long filenames
	if len(safe) > 64 {
		safe = safe[:64]
	}
	return safe
}

// historyPath returns the path for a session's history file.
func historyPath(sessionID string) string {
	safe := sanitizeSessionID(sessionID)
	return filepath.Join(HistoryDirFunc(), "history_"+safe+".json")
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
// TODO(v0.3): Use for cost_spark widget - cost trend sparkline
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
