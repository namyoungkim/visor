package transcript

// ToolStatus represents the execution state of a tool.
type ToolStatus string

const (
	ToolRunning   ToolStatus = "running"
	ToolCompleted ToolStatus = "completed"
	ToolError     ToolStatus = "error"
)

// Tool represents a grouped tool usage from the transcript.
// Tools with the same Name are grouped together, with Count tracking invocations.
type Tool struct {
	ID     string
	Name   string
	Status ToolStatus
	Count  int // Number of times this tool was invoked
}

// Agent represents a sub-agent spawned via the Task tool.
type Agent struct {
	ID          string
	Type        string // e.g., "Explore", "Plan"
	Status      string // "running" | "completed"
	Description string // Task description from input
	StartTime   int64  // Timestamp when tool_use was issued (ms)
	EndTime     int64  // Timestamp when tool_result was received (ms)
}

// Data holds the parsed transcript information for widgets.
type Data struct {
	Tools  []Tool
	Agents []Agent
}
