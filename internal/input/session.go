package input

// Session represents the parsed stdin JSON from Claude Code.
type Session struct {
	SessionID      string        `json:"session_id"`
	Model          Model         `json:"model"`
	Cost           Cost          `json:"cost"`
	ContextWindow  ContextWindow `json:"context_window"`
	Workspace      Workspace     `json:"workspace"`
	CurrentUsage   *CurrentUsage `json:"current_usage"`
	TranscriptPath string        `json:"transcript_path"`
	CWD            string        `json:"cwd"`
}

// Model contains model information.
type Model struct {
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
}

// Cost contains billing and API performance data.
type Cost struct {
	TotalCostUSD          float64 `json:"total_cost_usd"`
	TotalDurationMs       int64   `json:"total_duration_ms"`
	TotalAPIDurationMs    int64   `json:"total_api_duration_ms"`
	TotalAPICalls         int     `json:"total_api_calls"`
	TotalInputTokens      int     `json:"total_input_tokens"`
	TotalOutputTokens     int     `json:"total_output_tokens"`
	TotalCacheReadTokens  int     `json:"total_cache_read_tokens"`
	TotalCacheWriteTokens int     `json:"total_cache_write_tokens"`
}

// ContextWindow contains context usage information.
type ContextWindow struct {
	UsedPercentage    float64       `json:"used_percentage"`
	UsedTokens        int           `json:"used_tokens"`
	MaxTokens         int           `json:"max_tokens"`
	TotalInputTokens  int           `json:"total_input_tokens"`
	TotalOutputTokens int           `json:"total_output_tokens"`
	CurrentUsage      *CurrentUsage `json:"current_usage"`
}

// Workspace contains information about code changes.
type Workspace struct {
	LinesAdded   int `json:"lines_added"`
	LinesRemoved int `json:"lines_removed"`
	FilesChanged int `json:"files_changed"`
}

// CurrentUsage contains token usage for the current request.
type CurrentUsage struct {
	InputTokens          int `json:"input_tokens"`
	CacheReadTokens      int `json:"cache_read_tokens"`
	CacheReadInputTokens int `json:"cache_read_input_tokens"`
}

// GetCacheReadTokens returns the cache read token count,
// preferring cache_read_input_tokens over cache_read_tokens for compatibility.
func (cu *CurrentUsage) GetCacheReadTokens() int {
	if cu.CacheReadInputTokens > 0 {
		return cu.CacheReadInputTokens
	}
	return cu.CacheReadTokens
}

// GetCurrentUsage returns the CurrentUsage, checking ContextWindow first,
// then falling back to the top-level field for backward compatibility.
func (s *Session) GetCurrentUsage() *CurrentUsage {
	if s.ContextWindow.CurrentUsage != nil {
		return s.ContextWindow.CurrentUsage
	}
	return s.CurrentUsage
}

// GetTotalOutputTokens returns total output tokens, checking ContextWindow first,
// then falling back to Cost for backward compatibility.
func (s *Session) GetTotalOutputTokens() int {
	if s.ContextWindow.TotalOutputTokens > 0 {
		return s.ContextWindow.TotalOutputTokens
	}
	return s.Cost.TotalOutputTokens
}
