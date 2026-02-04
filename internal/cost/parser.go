package cost

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Entry represents a single cost entry from JSONL parsing.
type Entry struct {
	Timestamp    time.Time
	ModelID      string
	InputTokens  int
	OutputTokens int
	CacheRead    int
	CacheWrite   int
	CostUSD      float64
	SessionID    string
}

// jsonlMessage represents a message from the Claude transcript JSONL.
type jsonlMessage struct {
	Type         string           `json:"type"`
	Timestamp    string           `json:"timestamp"`
	Message      *assistantMessage `json:"message,omitempty"`
	CostUSD      float64          `json:"costUsd,omitempty"`
	DurationMs   int64            `json:"durationMs,omitempty"`
	SessionID    string           `json:"sessionId,omitempty"`
}

type assistantMessage struct {
	Model string `json:"model"`
	Usage *usage `json:"usage,omitempty"`
}

type usage struct {
	InputTokens                int `json:"input_tokens"`
	OutputTokens               int `json:"output_tokens"`
	CacheCreationInputTokens   int `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens       int `json:"cache_read_input_tokens,omitempty"`
}

// ParseJSONL parses a Claude transcript JSONL file and extracts cost entries.
func ParseJSONL(path string) ([]Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []Entry
	scanner := bufio.NewScanner(file)
	// Set a larger buffer for potentially large JSONL lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		entry, ok := parseJSONLLine(line)
		if ok {
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return entries, err
	}

	return entries, nil
}

// parseJSONLLine parses a single JSONL line into a cost entry.
func parseJSONLLine(line string) (Entry, bool) {
	var msg jsonlMessage
	if err := json.Unmarshal([]byte(line), &msg); err != nil {
		return Entry{}, false
	}

	// Only process assistant messages with usage data
	if msg.Type != "assistant" || msg.Message == nil || msg.Message.Usage == nil {
		return Entry{}, false
	}

	u := msg.Message.Usage
	entry := Entry{
		ModelID:      msg.Message.Model,
		InputTokens:  u.InputTokens,
		OutputTokens: u.OutputTokens,
		CacheRead:    u.CacheReadInputTokens,
		CacheWrite:   u.CacheCreationInputTokens,
		SessionID:    msg.SessionID,
	}

	// Parse timestamp
	if msg.Timestamp != "" {
		if t, err := time.Parse(time.RFC3339, msg.Timestamp); err == nil {
			entry.Timestamp = t
		}
	}

	// If costUsd is directly provided, use it
	if msg.CostUSD > 0 {
		entry.CostUSD = msg.CostUSD
	} else {
		// Calculate from tokens
		entry.CostUSD = CalculateCost(entry.ModelID, entry.InputTokens, entry.OutputTokens, entry.CacheRead, entry.CacheWrite)
	}

	return entry, true
}

// ParseAllSessions parses all JSONL files in the projects directory.
func ParseAllSessions(projectsDir string) ([]Entry, error) {
	var allEntries []Entry

	if projectsDir == "" {
		projectsDir = GetProjectsDir()
	}

	if projectsDir == "" {
		return allEntries, nil
	}

	// Walk all subdirectories
	err := filepath.Walk(projectsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if info.IsDir() {
			return nil
		}

		// Only process .jsonl files
		if !strings.HasSuffix(path, ".jsonl") {
			return nil
		}

		entries, err := ParseJSONL(path)
		if err != nil {
			return nil // Skip errors
		}

		allEntries = append(allEntries, entries...)
		return nil
	})

	if err != nil {
		return allEntries, err
	}

	return allEntries, nil
}

// ParseSession parses a single session's JSONL file.
func ParseSession(transcriptPath string) ([]Entry, error) {
	if transcriptPath == "" {
		return nil, nil
	}

	return ParseJSONL(transcriptPath)
}
