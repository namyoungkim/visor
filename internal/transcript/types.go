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

// TodoStatus represents the execution state of a todo task.
type TodoStatus string

const (
	TodoPending    TodoStatus = "pending"
	TodoInProgress TodoStatus = "in_progress"
	TodoCompleted  TodoStatus = "completed"
)

// Todo represents a task created via TaskCreate/TaskUpdate tools.
type Todo struct {
	ID      string
	Subject string
	Status  TodoStatus
}

// Data holds the parsed transcript information for widgets.
type Data struct {
	Tools  []Tool
	Agents []Agent
	Todos  []Todo
}
