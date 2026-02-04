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

func TestParse_ToolCount(t *testing.T) {
	// Multiple Read calls should be grouped together
	content := `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_001","name":"Read"}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_001"}]}}
{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_002","name":"Read"}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_002"}]}}
{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_003","name":"Read"}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_003"}]}}
{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_004","name":"Edit"}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_004"}]}}
`
	path := writeTempFile(t, content)
	defer os.Remove(path)

	data := Parse(path)

	// Should have 2 tools: Read and Edit (grouped by name)
	if len(data.Tools) != 2 {
		t.Fatalf("expected 2 tools (grouped by name), got %d", len(data.Tools))
	}

	// Find Read tool
	var readTool *Tool
	for i := range data.Tools {
		if data.Tools[i].Name == "Read" {
			readTool = &data.Tools[i]
			break
		}
	}
	if readTool == nil {
		t.Fatal("expected Read tool")
	}
	if readTool.Count != 3 {
		t.Errorf("expected Read count=3, got %d", readTool.Count)
	}
	if readTool.Status != ToolCompleted {
		t.Errorf("expected Read status=completed, got %s", readTool.Status)
	}
}

func TestParse_ToolRunningStatus(t *testing.T) {
	// Last invocation is still running
	content := `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_001","name":"Bash"}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_001"}]}}
{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_002","name":"Bash"}]}}
`
	path := writeTempFile(t, content)
	defer os.Remove(path)

	data := Parse(path)

	if len(data.Tools) != 1 {
		t.Fatalf("expected 1 tool (grouped by name), got %d", len(data.Tools))
	}
	if data.Tools[0].Count != 2 {
		t.Errorf("expected count=2, got %d", data.Tools[0].Count)
	}
	// Should be running because the latest invocation hasn't completed
	if data.Tools[0].Status != ToolRunning {
		t.Errorf("expected running status, got %s", data.Tools[0].Status)
	}
}

func TestParse_AgentDescription(t *testing.T) {
	content := `{"type":"assistant","timestamp":1000,"message":{"content":[{"type":"tool_use","id":"toolu_001","name":"Task","input":{"subagent_type":"Explore","description":"Analyze widget structure"}}]}}
{"type":"user","timestamp":43000,"message":{"content":[{"type":"tool_result","tool_use_id":"toolu_001"}]}}
`
	path := writeTempFile(t, content)
	defer os.Remove(path)

	data := Parse(path)

	if len(data.Agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(data.Agents))
	}
	if data.Agents[0].Description != "Analyze widget structure" {
		t.Errorf("expected description 'Analyze widget structure', got '%s'", data.Agents[0].Description)
	}
	if data.Agents[0].StartTime != 1000 {
		t.Errorf("expected StartTime=1000, got %d", data.Agents[0].StartTime)
	}
	if data.Agents[0].EndTime != 43000 {
		t.Errorf("expected EndTime=43000, got %d", data.Agents[0].EndTime)
	}
}

func TestParse_AgentRunning(t *testing.T) {
	content := `{"type":"assistant","timestamp":1000,"message":{"content":[{"type":"tool_use","id":"toolu_001","name":"Task","input":{"subagent_type":"Plan","description":"Plan implementation"}}]}}
`
	path := writeTempFile(t, content)
	defer os.Remove(path)

	data := Parse(path)

	if len(data.Agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(data.Agents))
	}
	if data.Agents[0].Status != "running" {
		t.Errorf("expected running status, got '%s'", data.Agents[0].Status)
	}
	if data.Agents[0].EndTime != 0 {
		t.Errorf("expected EndTime=0 for running agent, got %d", data.Agents[0].EndTime)
	}
}

func TestParseTimestamp_Int64(t *testing.T) {
	raw := []byte(`1738750823478`)
	ts := parseTimestamp(raw)
	if ts != 1738750823478 {
		t.Errorf("expected 1738750823478, got %d", ts)
	}
}

func TestParseTimestamp_ISO8601(t *testing.T) {
	raw := []byte(`"2026-02-05T13:40:23.478Z"`)
	ts := parseTimestamp(raw)
	// Should be non-zero
	if ts == 0 {
		t.Error("expected non-zero timestamp for ISO 8601")
	}
	// Should be in reasonable range (year 2026)
	if ts < 1700000000000 || ts > 2000000000000 {
		t.Errorf("timestamp %d seems out of range", ts)
	}
}

func TestParseTimestamp_Empty(t *testing.T) {
	ts := parseTimestamp(nil)
	if ts != 0 {
		t.Errorf("expected 0 for nil, got %d", ts)
	}
	ts = parseTimestamp([]byte{})
	if ts != 0 {
		t.Errorf("expected 0 for empty, got %d", ts)
	}
}

func TestGetMaxLines_Default(t *testing.T) {
	os.Unsetenv("VISOR_TRANSCRIPT_MAX_LINES")
	n := getMaxLines()
	if n != defaultMaxLines {
		t.Errorf("expected %d, got %d", defaultMaxLines, n)
	}
}

func TestGetMaxLines_EnvOverride(t *testing.T) {
	os.Setenv("VISOR_TRANSCRIPT_MAX_LINES", "1000")
	defer os.Unsetenv("VISOR_TRANSCRIPT_MAX_LINES")

	n := getMaxLines()
	if n != 1000 {
		t.Errorf("expected 1000, got %d", n)
	}
}

func TestGetMaxLines_InvalidEnv(t *testing.T) {
	os.Setenv("VISOR_TRANSCRIPT_MAX_LINES", "invalid")
	defer os.Unsetenv("VISOR_TRANSCRIPT_MAX_LINES")

	n := getMaxLines()
	if n != defaultMaxLines {
		t.Errorf("expected default %d for invalid env, got %d", defaultMaxLines, n)
	}
}

func TestParse_StringContent(t *testing.T) {
	// Content can be a string (e.g., text responses) - should not crash
	content := `{"type":"assistant","message":{"content":"Just a text response"}}
{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_001","name":"Read"}]}}
`
	path := writeTempFile(t, content)
	defer os.Remove(path)

	data := Parse(path)
	if len(data.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(data.Tools))
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
