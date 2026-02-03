package tui

// OptionType defines the type of a widget option
type OptionType int

const (
	OptionTypeBool OptionType = iota
	OptionTypeInt
	OptionTypeFloat
	OptionTypeString
)

// OptionDef defines metadata for a widget option
type OptionDef struct {
	Key          string
	Type         OptionType
	DefaultValue string
	Description  string
}

// WidgetMeta contains metadata for a widget
type WidgetMeta struct {
	Name        string
	Description string
	Options     []OptionDef
}

// AllWidgets returns metadata for all available widgets
func AllWidgets() []WidgetMeta {
	return []WidgetMeta{
		{
			Name:        "model",
			Description: "Display model name (e.g., Opus)",
			Options:     nil,
		},
		{
			Name:        "context",
			Description: "Context window usage with progress bar",
			Options: []OptionDef{
				{Key: "show_label", Type: OptionTypeBool, DefaultValue: "true", Description: "Show 'Ctx:' prefix"},
				{Key: "show_bar", Type: OptionTypeBool, DefaultValue: "true", Description: "Show progress bar"},
				{Key: "bar_width", Type: OptionTypeInt, DefaultValue: "10", Description: "Progress bar width"},
				{Key: "warn_threshold", Type: OptionTypeInt, DefaultValue: "60", Description: "Warning threshold %"},
				{Key: "critical_threshold", Type: OptionTypeInt, DefaultValue: "80", Description: "Critical threshold %"},
			},
		},
		{
			Name:        "context_spark",
			Description: "Context history sparkline",
			Options: []OptionDef{
				{Key: "width", Type: OptionTypeInt, DefaultValue: "8", Description: "Sparkline width"},
				{Key: "show_label", Type: OptionTypeBool, DefaultValue: "false", Description: "Show 'Ctx:' prefix"},
			},
		},
		{
			Name:        "compact_eta",
			Description: "Estimated time until context full",
			Options: []OptionDef{
				{Key: "show_when_above", Type: OptionTypeInt, DefaultValue: "40", Description: "Show only above this %"},
				{Key: "show_label", Type: OptionTypeBool, DefaultValue: "false", Description: "Show 'ETA:' prefix"},
			},
		},
		{
			Name:        "block_timer",
			Description: "Claude Pro rate limit block timer",
			Options: []OptionDef{
				{Key: "show_label", Type: OptionTypeBool, DefaultValue: "true", Description: "Show 'Block:' prefix"},
				{Key: "warn_threshold", Type: OptionTypeInt, DefaultValue: "80", Description: "Warning at % elapsed"},
				{Key: "critical_threshold", Type: OptionTypeInt, DefaultValue: "95", Description: "Critical at % elapsed"},
			},
		},
		{
			Name:        "cache_hit",
			Description: "API cache hit rate",
			Options: []OptionDef{
				{Key: "show_label", Type: OptionTypeBool, DefaultValue: "false", Description: "Show 'Cache:' prefix"},
				{Key: "good_threshold", Type: OptionTypeInt, DefaultValue: "80", Description: "Good/green threshold %"},
				{Key: "warn_threshold", Type: OptionTypeInt, DefaultValue: "50", Description: "Warning threshold %"},
			},
		},
		{
			Name:        "api_latency",
			Description: "Average API response time",
			Options: []OptionDef{
				{Key: "show_label", Type: OptionTypeBool, DefaultValue: "false", Description: "Show 'Latency:' prefix"},
				{Key: "warn_threshold", Type: OptionTypeInt, DefaultValue: "2000", Description: "Warning threshold ms"},
				{Key: "critical_threshold", Type: OptionTypeInt, DefaultValue: "5000", Description: "Critical threshold ms"},
			},
		},
		{
			Name:        "cost",
			Description: "Session cost in USD",
			Options: []OptionDef{
				{Key: "show_label", Type: OptionTypeBool, DefaultValue: "false", Description: "Show 'Cost:' prefix"},
				{Key: "warn_threshold", Type: OptionTypeFloat, DefaultValue: "0.5", Description: "Warning threshold USD"},
				{Key: "critical_threshold", Type: OptionTypeFloat, DefaultValue: "1.0", Description: "Critical threshold USD"},
			},
		},
		{
			Name:        "burn_rate",
			Description: "Cost per minute",
			Options: []OptionDef{
				{Key: "show_label", Type: OptionTypeBool, DefaultValue: "false", Description: "Show 'Burn:' prefix"},
				{Key: "warn_threshold", Type: OptionTypeInt, DefaultValue: "10", Description: "Warning threshold cents/min"},
				{Key: "critical_threshold", Type: OptionTypeInt, DefaultValue: "25", Description: "Critical threshold cents/min"},
			},
		},
		{
			Name:        "code_changes",
			Description: "Lines added/removed in session",
			Options:     nil,
		},
		{
			Name:        "git",
			Description: "Git branch and status",
			Options:     nil,
		},
		{
			Name:        "tools",
			Description: "Active tool calls status",
			Options: []OptionDef{
				{Key: "max_display", Type: OptionTypeInt, DefaultValue: "3", Description: "Max tools to display"},
				{Key: "show_label", Type: OptionTypeBool, DefaultValue: "false", Description: "Show 'Tools:' prefix"},
			},
		},
		{
			Name:        "agents",
			Description: "Active agent status",
			Options: []OptionDef{
				{Key: "max_display", Type: OptionTypeInt, DefaultValue: "3", Description: "Max agents to display"},
				{Key: "show_label", Type: OptionTypeBool, DefaultValue: "false", Description: "Show 'Agents:' prefix"},
			},
		},
		// v0.6 Cost tracking widgets
		{
			Name:        "daily_cost",
			Description: "Today's aggregated cost",
			Options: []OptionDef{
				{Key: "show_label", Type: OptionTypeBool, DefaultValue: "false", Description: "Show 'Today:' prefix"},
				{Key: "warn_threshold", Type: OptionTypeFloat, DefaultValue: "5.0", Description: "Warning threshold USD"},
				{Key: "critical_threshold", Type: OptionTypeFloat, DefaultValue: "10.0", Description: "Critical threshold USD"},
			},
		},
		{
			Name:        "weekly_cost",
			Description: "This week's aggregated cost",
			Options: []OptionDef{
				{Key: "show_label", Type: OptionTypeBool, DefaultValue: "false", Description: "Show 'Week:' prefix"},
				{Key: "warn_threshold", Type: OptionTypeFloat, DefaultValue: "25.0", Description: "Warning threshold USD"},
				{Key: "critical_threshold", Type: OptionTypeFloat, DefaultValue: "50.0", Description: "Critical threshold USD"},
			},
		},
		{
			Name:        "block_cost",
			Description: "Cost in current 5-hour block",
			Options: []OptionDef{
				{Key: "show_label", Type: OptionTypeBool, DefaultValue: "false", Description: "Show 'Block$:' prefix"},
				{Key: "warn_threshold", Type: OptionTypeFloat, DefaultValue: "2.0", Description: "Warning threshold USD"},
				{Key: "critical_threshold", Type: OptionTypeFloat, DefaultValue: "5.0", Description: "Critical threshold USD"},
			},
		},
		// v0.6 Usage limit widgets
		{
			Name:        "block_limit",
			Description: "5-hour rate limit utilization",
			Options: []OptionDef{
				{Key: "show_label", Type: OptionTypeBool, DefaultValue: "true", Description: "Show '5h:' prefix"},
				{Key: "show_remaining", Type: OptionTypeBool, DefaultValue: "true", Description: "Show remaining time"},
				{Key: "warn_threshold", Type: OptionTypeInt, DefaultValue: "70", Description: "Warning threshold %"},
				{Key: "critical_threshold", Type: OptionTypeInt, DefaultValue: "90", Description: "Critical threshold %"},
			},
		},
		{
			Name:        "week_limit",
			Description: "7-day rate limit utilization",
			Options: []OptionDef{
				{Key: "show_label", Type: OptionTypeBool, DefaultValue: "true", Description: "Show '7d:' prefix"},
				{Key: "show_remaining", Type: OptionTypeBool, DefaultValue: "false", Description: "Show remaining time"},
				{Key: "warn_threshold", Type: OptionTypeInt, DefaultValue: "70", Description: "Warning threshold %"},
				{Key: "critical_threshold", Type: OptionTypeInt, DefaultValue: "90", Description: "Critical threshold %"},
			},
		},
	}
}

// GetWidgetMeta returns metadata for a specific widget
func GetWidgetMeta(name string) *WidgetMeta {
	for _, w := range AllWidgets() {
		if w.Name == name {
			return &w
		}
	}
	return nil
}

// AvailableWidgetNames returns a list of all widget names
func AvailableWidgetNames() []string {
	widgets := AllWidgets()
	names := make([]string, len(widgets))
	for i, w := range widgets {
		names[i] = w.Name
	}
	return names
}
