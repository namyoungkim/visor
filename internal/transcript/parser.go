package transcript

import (
	"bufio"
	"encoding/json"
	"os"
)

// maxLines limits the number of transcript lines to parse.
// 100 lines is sufficient for typical sessions:
// - Average tool call produces ~2 lines (tool_use + tool_result)
// - 100 lines ≈ 50 tool invocations worth of history
// - Keeps memory usage bounded for long-running sessions
// See: https://github.com/namyoungkim/visor/issues/16 for optimization plans
const maxLines = 100

// transcriptEntry represents a single line in the JSONL transcript.
type transcriptEntry struct {
	Type    string `json:"type"`
	Message struct {
		Content []contentBlock `json:"content"`
	} `json:"message"`
	Data struct {
		AgentID string `json:"agentId"`
	} `json:"data"`
	ToolUseID string `json:"toolUseID"`
}

// contentBlock represents a content block in the message.
type contentBlock struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	ToolUseID string `json:"tool_use_id"`
	IsError   *bool  `json:"is_error"`
	Input     struct {
		SubagentType string `json:"subagent_type"`
	} `json:"input"`
}

// Parse reads a JSONL transcript file and extracts tool/agent data.
// Returns empty Data on any error (graceful fallback).
func Parse(path string) *Data {
	if path == "" {
		return &Data{}
	}

	lines, err := tailLines(path, maxLines)
	if err != nil {
		return &Data{}
	}

	return parseLines(lines)
}

// tailLines reads the last n lines from a file.
func tailLines(path string, n int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	// Increase buffer size for large JSON lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > n {
			lines = lines[1:]
		}
	}

	return lines, scanner.Err()
}

// parseLines processes JSONL lines and extracts tools/agents.
func parseLines(lines []string) *Data {
	data := &Data{
		Tools:  make([]Tool, 0),
		Agents: make([]Agent, 0),
	}

	toolMap := make(map[string]*Tool)
	agentMap := make(map[string]*Agent)
	var toolOrder []string  // Maintain insertion order
	var agentOrder []string // Maintain insertion order

	for _, line := range lines {
		var entry transcriptEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		switch entry.Type {
		case "assistant":
			processAssistant(&entry, toolMap, agentMap, &toolOrder, &agentOrder)
		case "user":
			processToolResult(&entry, toolMap, agentMap)
		}
	}

	// Convert maps to slices in insertion order
	for _, id := range toolOrder {
		if tool, ok := toolMap[id]; ok {
			data.Tools = append(data.Tools, *tool)
		}
	}
	for _, id := range agentOrder {
		if agent, ok := agentMap[id]; ok {
			data.Agents = append(data.Agents, *agent)
		}
	}

	return data
}

// processAssistant handles assistant messages containing tool_use.
func processAssistant(entry *transcriptEntry, toolMap map[string]*Tool, agentMap map[string]*Agent, toolOrder, agentOrder *[]string) {
	for _, block := range entry.Message.Content {
		if block.Type != "tool_use" {
			continue
		}

		// Track tool (only if not already seen)
		if _, exists := toolMap[block.ID]; !exists {
			tool := &Tool{
				ID:     block.ID,
				Name:   block.Name,
				Status: ToolRunning,
			}
			toolMap[block.ID] = tool
			*toolOrder = append(*toolOrder, block.ID)
		}

		// Check if this is a Task tool (spawns agent)
		if block.Name == "Task" && block.Input.SubagentType != "" {
			if _, exists := agentMap[block.ID]; !exists {
				agent := &Agent{
					ID:     block.ID,
					Type:   block.Input.SubagentType,
					Status: "running",
				}
				agentMap[block.ID] = agent
				*agentOrder = append(*agentOrder, block.ID)
			}
		}
	}
}

// processToolResult handles tool_result messages to update tool and agent status.
func processToolResult(entry *transcriptEntry, toolMap map[string]*Tool, agentMap map[string]*Agent) {
	for _, block := range entry.Message.Content {
		if block.Type != "tool_result" {
			continue
		}

		// Update tool status
		if tool, ok := toolMap[block.ToolUseID]; ok {
			if block.IsError != nil && *block.IsError {
				tool.Status = ToolError
			} else {
				tool.Status = ToolCompleted
			}
		}

		// Update agent status (Task tool completion)
		if agent, ok := agentMap[block.ToolUseID]; ok {
			agent.Status = "completed"
		}
	}
}
