package transcript

import (
	"bufio"
	"encoding/json"
	"os"
)

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

	for _, line := range lines {
		var entry transcriptEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		switch entry.Type {
		case "assistant":
			processAssistant(&entry, toolMap, agentMap)
		case "user":
			processToolResult(&entry, toolMap)
		case "progress":
			processProgress(&entry, agentMap)
		}
	}

	// Convert maps to slices (maintaining insertion order via map iteration is fine for display)
	for _, tool := range toolMap {
		data.Tools = append(data.Tools, *tool)
	}
	for _, agent := range agentMap {
		data.Agents = append(data.Agents, *agent)
	}

	return data
}

// processAssistant handles assistant messages containing tool_use.
func processAssistant(entry *transcriptEntry, toolMap map[string]*Tool, agentMap map[string]*Agent) {
	for _, block := range entry.Message.Content {
		if block.Type != "tool_use" {
			continue
		}

		// Track tool
		tool := &Tool{
			ID:     block.ID,
			Name:   block.Name,
			Status: ToolRunning,
		}
		toolMap[block.ID] = tool

		// Check if this is a Task tool (spawns agent)
		if block.Name == "Task" && block.Input.SubagentType != "" {
			agent := &Agent{
				ID:     block.ID,
				Type:   block.Input.SubagentType,
				Status: "running",
			}
			agentMap[block.ID] = agent
		}
	}
}

// processToolResult handles tool_result messages to update tool status.
func processToolResult(entry *transcriptEntry, toolMap map[string]*Tool) {
	for _, block := range entry.Message.Content {
		if block.Type != "tool_result" {
			continue
		}

		tool, ok := toolMap[block.ToolUseID]
		if !ok {
			continue
		}

		if block.IsError != nil && *block.IsError {
			tool.Status = ToolError
		} else {
			tool.Status = ToolCompleted
		}
	}
}

// processProgress handles progress messages to track agent status.
func processProgress(entry *transcriptEntry, agentMap map[string]*Agent) {
	// Progress entries with agentId indicate running agents
	if entry.Data.AgentID == "" {
		return
	}

	// Try to find and update the agent by matching toolUseID
	for _, agent := range agentMap {
		if agent.Status == "running" {
			// Agent is still running, nothing to update
			return
		}
	}
}
