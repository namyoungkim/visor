package input

// Session represents the parsed stdin JSON from Claude Code.
type Session struct {
	Model         Model         `json:"model"`
	Cost          Cost          `json:"cost"`
	ContextWindow ContextWindow `json:"context_window"`
	Workspace     Workspace     `json:"workspace"`
	CurrentUsage  *CurrentUsage `json:"current_usage"`
}

// Model contains model information.
type Model struct {
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
}

// Cost contains billing and API performance data.
type Cost struct {
	TotalCostUSD          float64 `json:"total_cost_usd"`
	TotalAPIDurationMs    int64   `json:"total_api_duration_ms"`
	TotalAPICalls         int     `json:"total_api_calls"`
	TotalInputTokens      int     `json:"total_input_tokens"`
	TotalOutputTokens     int     `json:"total_output_tokens"`
	TotalCacheReadTokens  int     `json:"total_cache_read_tokens"`
	TotalCacheWriteTokens int     `json:"total_cache_write_tokens"`
}

// ContextWindow contains context usage information.
type ContextWindow struct {
	UsedPercentage float64 `json:"used_percentage"`
	UsedTokens     int     `json:"used_tokens"`
	MaxTokens      int     `json:"max_tokens"`
}

// Workspace contains information about code changes.
type Workspace struct {
	LinesAdded   int `json:"lines_added"`
	LinesRemoved int `json:"lines_removed"`
	FilesChanged int `json:"files_changed"`
}

// CurrentUsage contains token usage for the current request.
type CurrentUsage struct {
	InputTokens     int `json:"input_tokens"`
	CacheReadTokens int `json:"cache_read_tokens"`
}
