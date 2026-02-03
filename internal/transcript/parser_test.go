package transcript

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse_EmptyPath(t *testing.T) {
	data := Parse("")
	if data == nil {
		t.Fatal("expected non-nil Data")
	}
	if len(data.Tools) != 0 {
		t.Errorf("expected empty Tools, got %d", len(data.Tools))
	}
	if len(data.Agents) != 0 {
		t.Errorf("expected empty Agents, got %d", len(data.Agents))
	}
}

func TestParse_NonExistentFile(t *testing.T) {
	data := Parse("/nonexistent/path/file.jsonl")
	if data == nil {
		t.Fatal("expected non-nil Data")
	}
	if len(data.Tools) != 0 {
		t.Errorf("expected empty Tools, got %d", len(data.Tools))
	}
}

func TestParse_ToolUseAndResult(t *testing.T) {
	content := `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_001","name":"Read"}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_001"}]}}
{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_002","name":"Write"}]}}
`
	path := writeTempFile(t, content)
	defer os.Remove(path)

	data := Parse(path)
	if len(data.Tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(data.Tools))
	}

	// Check that we have both Read (completed) and Write (running)
	toolStatus := make(map[string]ToolStatus)
	for _, tool := range data.Tools {
		toolStatus[tool.Name] = tool.Status
	}

	if toolStatus["Read"] != ToolCompleted {
		t.Errorf("expected Read to be completed, got %s", toolStatus["Read"])
	}
	if toolStatus["Write"] != ToolRunning {
		t.Errorf("expected Write to be running, got %s", toolStatus["Write"])
	}
}

func TestParse_ToolError(t *testing.T) {
	content := `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_001","name":"Bash"}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_001","is_error":true}]}}
`
	path := writeTempFile(t, content)
	defer os.Remove(path)

	data := Parse(path)
	if len(data.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(data.Tools))
	}
	if data.Tools[0].Status != ToolError {
		t.Errorf("expected error status, got %s", data.Tools[0].Status)
	}
}

func TestParse_TaskAgent(t *testing.T) {
	content := `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_001","name":"Task","input":{"subagent_type":"Explore"}}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_001"}]}}
`
	path := writeTempFile(t, content)
	defer os.Remove(path)

	data := Parse(path)
	if len(data.Agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(data.Agents))
	}
	if data.Agents[0].Type != "Explore" {
		t.Errorf("expected Explore agent type, got %s", data.Agents[0].Type)
	}
	if data.Agents[0].Status != "completed" {
		t.Errorf("expected completed status, got %s", data.Agents[0].Status)
	}
}

func TestParse_MalformedJSON(t *testing.T) {
	content := `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_001","name":"Read"}]}}
not valid json
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_001"}]}}
`
	path := writeTempFile(t, content)
	defer os.Remove(path)

	data := Parse(path)
	// Should still parse valid lines
	if len(data.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(data.Tools))
	}
	if data.Tools[0].Status != ToolCompleted {
		t.Errorf("expected completed status, got %s", data.Tools[0].Status)
	}
}

func TestParse_ToolOrder(t *testing.T) {
	content := `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_001","name":"Read"}]}}
{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_002","name":"Write"}]}}
{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_003","name":"Bash"}]}}
`
	path := writeTempFile(t, content)
	defer os.Remove(path)

	data := Parse(path)
	if len(data.Tools) != 3 {
		t.Fatalf("expected 3 tools, got %d", len(data.Tools))
	}

	// Verify insertion order is preserved
	expectedOrder := []string{"Read", "Write", "Bash"}
	for i, expected := range expectedOrder {
		if data.Tools[i].Name != expected {
			t.Errorf("tool %d: expected %s, got %s", i, expected, data.Tools[i].Name)
		}
	}
}

func TestTailLines(t *testing.T) {
	content := "line1\nline2\nline3\nline4\nline5\n"
	path := writeTempFile(t, content)
	defer os.Remove(path)

	lines, err := tailLines(path, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "line3" {
		t.Errorf("expected 'line3', got '%s'", lines[0])
	}
	if lines[2] != "line5" {
		t.Errorf("expected 'line5', got '%s'", lines[2])
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.jsonl")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}
