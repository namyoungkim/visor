package transcript

// ToolStatus represents the execution state of a tool.
type ToolStatus string

const (
	ToolRunning   ToolStatus = "running"
	ToolCompleted ToolStatus = "completed"
	ToolError     ToolStatus = "error"
)

// Tool represents a single tool invocation from the transcript.
type Tool struct {
	ID     string
	Name   string
	Status ToolStatus
}

// Agent represents a sub-agent spawned via the Task tool.
type Agent struct {
	ID     string
	Type   string // e.g., "Explore", "Plan"
	Status string // "running" | "completed"
}

// Data holds the parsed transcript information for widgets.
type Data struct {
	Tools  []Tool
	Agents []Agent
}
