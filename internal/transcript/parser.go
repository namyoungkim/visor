package transcript

import (
	"encoding/json"
	"os"
)

// maxLines limits the number of transcript lines to parse.
// 100 lines is sufficient for typical sessions:
// - Average tool call produces ~2 lines (tool_use + tool_result)
// - 100 lines ≈ 50 tool invocations worth of history
// - Keeps memory usage bounded for long-running sessions
const maxLines = 100

// transcriptEntry represents a single line in the JSONL transcript.
type transcriptEntry struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"` // Unix timestamp in milliseconds
	Message   struct {
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
		Description  string `json:"description"` // Task description
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

// tailLines reads the last n lines from a file efficiently.
// It seeks from EOF and reads backwards to avoid loading the entire file.
func tailLines(path string, n int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Get file size
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := stat.Size()

	if fileSize == 0 {
		return nil, nil
	}

	// Start with an estimated chunk size (average JSONL line ~2KB, read extra for safety)
	chunkSize := int64(n * 4 * 1024) // 4KB per line estimate
	if chunkSize > fileSize {
		chunkSize = fileSize
	}

	var lines []string
	offset := fileSize

	for len(lines) < n && offset > 0 {
		// Calculate read position
		offset -= chunkSize
		if offset < 0 {
			chunkSize += offset // Adjust chunk size for final read
			offset = 0
		}

		// Seek and read chunk
		if _, err := f.Seek(offset, 0); err != nil {
			return nil, err
		}

		chunk := make([]byte, chunkSize)
		bytesRead, err := f.Read(chunk)
		if err != nil {
			return nil, err
		}
		chunk = chunk[:bytesRead]

		// Parse lines from chunk
		chunkLines := splitLines(chunk)

		// If not at start of file, first line might be partial - discard it
		if offset > 0 && len(chunkLines) > 0 {
			chunkLines = chunkLines[1:]
		}

		// Prepend to existing lines
		lines = append(chunkLines, lines...)

		// Double chunk size for next iteration if needed
		chunkSize *= 2
		if chunkSize > 1024*1024 { // Cap at 1MB chunks
			chunkSize = 1024 * 1024
		}
	}

	// Return only the last n lines
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}

	return lines, nil
}

// splitLines splits a byte slice into lines, handling \n and \r\n.
func splitLines(data []byte) []string {
	var lines []string
	start := 0

	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			end := i
			if end > start && data[end-1] == '\r' {
				end--
			}
			if end > start { // Skip empty lines
				lines = append(lines, string(data[start:end]))
			}
			start = i + 1
		}
	}

	// Handle last line without newline
	if start < len(data) {
		lines = append(lines, string(data[start:]))
	}

	return lines
}

// parseLines processes JSONL lines and extracts tools/agents.
func parseLines(lines []string) *Data {
	data := &Data{
		Tools:  make([]Tool, 0),
		Agents: make([]Agent, 0),
	}

	// toolMap groups tools by Name (not ID) to count invocations
	toolMap := make(map[string]*Tool)    // key: tool Name
	toolIDMap := make(map[string]string) // key: tool ID -> tool Name (for result lookup)
	agentMap := make(map[string]*Agent)  // key: tool ID
	var toolOrder []string               // Maintain insertion order by Name
	var agentOrder []string              // Maintain insertion order by ID

	for _, line := range lines {
		var entry transcriptEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		switch entry.Type {
		case "assistant":
			processAssistant(&entry, toolMap, toolIDMap, agentMap, &toolOrder, &agentOrder)
		case "user":
			processToolResult(&entry, toolMap, toolIDMap, agentMap)
		}
	}

	// Convert maps to slices in insertion order
	for _, name := range toolOrder {
		if tool, ok := toolMap[name]; ok {
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
func processAssistant(entry *transcriptEntry, toolMap map[string]*Tool, toolIDMap map[string]string, agentMap map[string]*Agent, toolOrder, agentOrder *[]string) {
	for _, block := range entry.Message.Content {
		if block.Type != "tool_use" {
			continue
		}

		// Track tool by Name (group same tools together)
		toolIDMap[block.ID] = block.Name // Map ID -> Name for result lookup

		if existing, exists := toolMap[block.Name]; exists {
			// Tool already seen: increment count, update status to running
			existing.Count++
			existing.Status = ToolRunning
			existing.ID = block.ID // Update to latest ID
		} else {
			// New tool: create entry
			tool := &Tool{
				ID:     block.ID,
				Name:   block.Name,
				Status: ToolRunning,
				Count:  1,
			}
			toolMap[block.Name] = tool
			*toolOrder = append(*toolOrder, block.Name)
		}

		// Check if this is a Task tool (spawns agent)
		if block.Name == "Task" && block.Input.SubagentType != "" {
			if _, exists := agentMap[block.ID]; !exists {
				agent := &Agent{
					ID:          block.ID,
					Type:        block.Input.SubagentType,
					Status:      "running",
					Description: block.Input.Description,
					StartTime:   entry.Timestamp,
				}
				agentMap[block.ID] = agent
				*agentOrder = append(*agentOrder, block.ID)
			}
		}
	}
}

// processToolResult handles tool_result messages to update tool and agent status.
func processToolResult(entry *transcriptEntry, toolMap map[string]*Tool, toolIDMap map[string]string, agentMap map[string]*Agent) {
	for _, block := range entry.Message.Content {
		if block.Type != "tool_result" {
			continue
		}

		// Update tool status (lookup by ID -> Name)
		if toolName, ok := toolIDMap[block.ToolUseID]; ok {
			if tool, ok := toolMap[toolName]; ok {
				// Only update status if this is the latest invocation
				if tool.ID == block.ToolUseID {
					if block.IsError != nil && *block.IsError {
						tool.Status = ToolError
					} else {
						tool.Status = ToolCompleted
					}
				}
			}
		}

		// Update agent status (Task tool completion)
		if agent, ok := agentMap[block.ToolUseID]; ok {
			agent.Status = "completed"
			agent.EndTime = entry.Timestamp
		}
	}
}
