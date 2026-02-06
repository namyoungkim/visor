package cost

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseJSONLLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		want    bool
		wantErr bool
	}{
		{
			name:    "valid assistant message",
			line:    `{"type":"assistant","timestamp":"2024-01-15T10:30:00Z","message":{"model":"claude-3-5-sonnet-20241022","usage":{"input_tokens":1000,"output_tokens":500}}}`,
			want:    true,
			wantErr: false,
		},
		{
			name:    "user message parsed as user turn",
			line:    `{"type":"user","timestamp":"2024-01-15T10:30:00Z","sessionId":"s1"}`,
			want:    true,
			wantErr: false,
		},
		{
			name:    "meta user message ignored",
			line:    `{"type":"user","isMeta":true,"timestamp":"2024-01-15T10:30:00Z"}`,
			want:    false,
			wantErr: false,
		},
		{
			name:    "invalid json",
			line:    `{invalid json}`,
			want:    false,
			wantErr: false,
		},
		{
			name:    "empty line",
			line:    "",
			want:    false,
			wantErr: false,
		},
		{
			name:    "assistant without usage",
			line:    `{"type":"assistant","message":{"model":"claude-3-5-sonnet"}}`,
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := parseJSONLLine(tt.line)
			if ok != tt.want {
				t.Errorf("parseJSONLLine() ok = %v, want %v", ok, tt.want)
			}
		})
	}
}

func TestParseJSONLLineExtractsData(t *testing.T) {
	line := `{"type":"assistant","timestamp":"2024-01-15T10:30:00Z","message":{"model":"claude-3-5-sonnet-20241022","usage":{"input_tokens":1000,"output_tokens":500}},"sessionId":"test-session"}`

	entry, ok := parseJSONLLine(line)
	if !ok {
		t.Fatal("parseJSONLLine() failed")
	}

	if entry.ModelID != "claude-3-5-sonnet-20241022" {
		t.Errorf("ModelID = %q, want claude-3-5-sonnet-20241022", entry.ModelID)
	}

	if entry.InputTokens != 1000 {
		t.Errorf("InputTokens = %d, want 1000", entry.InputTokens)
	}

	if entry.OutputTokens != 500 {
		t.Errorf("OutputTokens = %d, want 500", entry.OutputTokens)
	}

	if entry.SessionID != "test-session" {
		t.Errorf("SessionID = %q, want test-session", entry.SessionID)
	}
}

func TestParseJSONLLine_UserTurn(t *testing.T) {
	line := `{"type":"user","timestamp":"2024-01-15T10:30:00Z","sessionId":"s1"}`
	entry, ok := parseJSONLLine(line)
	if !ok {
		t.Fatal("parseJSONLLine() should parse user turn")
	}
	if !entry.IsUserTurn {
		t.Error("IsUserTurn should be true for user message")
	}
	if entry.SessionID != "s1" {
		t.Errorf("SessionID = %q, want s1", entry.SessionID)
	}
	if entry.CostUSD != 0 {
		t.Errorf("CostUSD = %v, want 0 for user turn", entry.CostUSD)
	}
}

func TestParseJSONLLine_MetaUser(t *testing.T) {
	line := `{"type":"user","isMeta":true,"timestamp":"2024-01-15T10:30:00Z","sessionId":"s1"}`
	_, ok := parseJSONLLine(line)
	if ok {
		t.Error("parseJSONLLine() should not parse meta user message")
	}
}

func TestParseJSONLLine_AssistantIsNotUserTurn(t *testing.T) {
	line := `{"type":"assistant","timestamp":"2024-01-15T10:30:00Z","message":{"model":"claude-3-5-sonnet-20241022","usage":{"input_tokens":1000,"output_tokens":500}}}`
	entry, ok := parseJSONLLine(line)
	if !ok {
		t.Fatal("parseJSONLLine() should parse assistant message")
	}
	if entry.IsUserTurn {
		t.Error("IsUserTurn should be false for assistant message")
	}
}

func TestParseJSONL(t *testing.T) {
	// Create a temporary JSONL file
	dir := t.TempDir()
	path := filepath.Join(dir, "test.jsonl")

	content := `{"type":"user","timestamp":"2024-01-15T10:29:00Z","sessionId":"s1"}
{"type":"assistant","timestamp":"2024-01-15T10:30:00Z","message":{"model":"claude-3-5-sonnet-20241022","usage":{"input_tokens":1000,"output_tokens":500}}}
{"type":"assistant","timestamp":"2024-01-15T10:31:00Z","message":{"model":"claude-3-5-sonnet-20241022","usage":{"input_tokens":2000,"output_tokens":1000}}}
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	entries, err := ParseJSONL(path)
	if err != nil {
		t.Fatalf("ParseJSONL() error = %v", err)
	}

	if len(entries) != 3 {
		t.Errorf("ParseJSONL() returned %d entries, want 3", len(entries))
	}

	// First entry should be the user turn
	userTurns := 0
	for _, e := range entries {
		if e.IsUserTurn {
			userTurns++
		}
	}
	if userTurns != 1 {
		t.Errorf("ParseJSONL() found %d user turns, want 1", userTurns)
	}
}

func TestParseJSONLFileNotFound(t *testing.T) {
	_, err := ParseJSONL("/nonexistent/file.jsonl")
	if err == nil {
		t.Error("ParseJSONL() should return error for nonexistent file")
	}
}

func TestParseJSONLEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.jsonl")

	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	entries, err := ParseJSONL(path)
	if err != nil {
		t.Fatalf("ParseJSONL() error = %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("ParseJSONL() returned %d entries for empty file, want 0", len(entries))
	}
}
