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
			name:    "user message ignored",
			line:    `{"type":"user","content":"hello"}`,
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

func TestParseJSONL(t *testing.T) {
	// Create a temporary JSONL file
	dir := t.TempDir()
	path := filepath.Join(dir, "test.jsonl")

	content := `{"type":"user","content":"hello"}
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

	if len(entries) != 2 {
		t.Errorf("ParseJSONL() returned %d entries, want 2", len(entries))
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
