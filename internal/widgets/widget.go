package widgets

import (
	"strconv"
	"strings"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/history"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/transcript"
)

// Threshold constants for color coding.
const (
	// Context window thresholds
	ContextWarningPct = 60.0
	ContextDangerPct  = 80.0

	// Cost thresholds (USD)
	CostWarningUSD = 0.5
	CostDangerUSD  = 1.0

	// Cache hit rate thresholds (inverse: higher is better)
	CacheHitGoodPct    = 80.0
	CacheHitWarningPct = 50.0

	// API latency thresholds (ms)
	LatencyWarningMs = 2000
	LatencyDangerMs  = 5000

	// Block timer thresholds (percentage elapsed)
	BlockTimerWarningPct  = 80.0 // 80% elapsed = 1 hour remaining
	BlockTimerCriticalPct = 95.0 // 95% elapsed = 15 minutes remaining
)

// ColorByThreshold returns a color based on value and thresholds.
// For metrics where higher is worse (cost, latency, context usage).
func ColorByThreshold(value, warning, danger float64) string {
	if value >= danger {
		return "red"
	} else if value >= warning {
		return "yellow"
	}
	return "green"
}

// ColorByThresholdInverse returns a color based on value and thresholds.
// For metrics where higher is better (cache hit rate).
func ColorByThresholdInverse(value, good, warning float64) string {
	if value >= good {
		return "green"
	} else if value >= warning {
		return "yellow"
	}
	return "red"
}

// FormatOutput applies custom format if specified, otherwise uses default.
// Format string can use {value} placeholder.
// Example: format="Context: {value}" with value="42%" → "Context: 42%"
func FormatOutput(cfg *config.WidgetConfig, defaultFormat, value string) string {
	format := cfg.Format
	if format == "" {
		format = defaultFormat
	}

	// If no format specified, return value as-is
	if format == "" {
		return value
	}

	// Simple placeholder replacement
	result := format
	for i := 0; i <= len(result)-7; i++ {
		if result[i:i+7] == "{value}" {
			result = result[:i] + value + result[i+7:]
			break
		}
	}

	return result
}

// GetExtra returns a value from the Extra map, or defaultValue if not found.
func GetExtra(cfg *config.WidgetConfig, key, defaultValue string) string {
	if cfg.Extra == nil {
		return defaultValue
	}
	if v, ok := cfg.Extra[key]; ok {
		return v
	}
	return defaultValue
}

// GetExtraBool returns a boolean value from the Extra map.
func GetExtraBool(cfg *config.WidgetConfig, key string, defaultValue bool) bool {
	v := GetExtra(cfg, key, "")
	if v == "" {
		return defaultValue
	}
	return v == "true" || v == "1" || v == "yes"
}

// GetExtraInt returns an integer value from the Extra map.
func GetExtraInt(cfg *config.WidgetConfig, key string, defaultValue int) int {
	v := GetExtra(cfg, key, "")
	if v == "" {
		return defaultValue
	}
	if n, err := strconv.Atoi(v); err == nil {
		return n
	}
	return defaultValue
}

// GetExtraFloat returns a float64 value from the Extra map.
func GetExtraFloat(cfg *config.WidgetConfig, key string, defaultValue float64) float64 {
	v := GetExtra(cfg, key, "")
	if v == "" {
		return defaultValue
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return f
	}
	return defaultValue
}

// Progress bar characters.
const (
	BarFilled = "█"
	BarEmpty  = "░"
)

// DefaultBarWidth is the default width for progress bars.
const DefaultBarWidth = 10

// ProgressBar returns a progress bar string for the given percentage.
// The filled portion is calculated using truncation (floor), not rounding.
// Example: ProgressBar(42.0, 10) returns "████░░░░░░" (4.2 → 4 filled)
func ProgressBar(pct float64, width int) string {
	if width <= 0 {
		width = DefaultBarWidth
	}

	// Clamp percentage to 0-100
	if pct < 0 {
		pct = 0
	} else if pct > 100 {
		pct = 100
	}

	// Truncate to int (floor behavior): 42% of 10 = 4.2 → 4
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}

	return strings.Repeat(BarFilled, filled) + strings.Repeat(BarEmpty, width-filled)
}

// Widget is the interface all widgets must implement.
type Widget interface {
	Name() string
	Render(session *input.Session, cfg *config.WidgetConfig) string
	ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool
}

// Registry holds all registered widgets.
var Registry = make(map[string]Widget)

// Register adds a widget to the registry.
func Register(w Widget) {
	Registry[w.Name()] = w
}

// Get returns a widget by name.
func Get(name string) (Widget, bool) {
	w, ok := Registry[name]
	return w, ok
}

// RenderAll renders all widgets for a line configuration.
func RenderAll(session *input.Session, widgets []config.WidgetConfig) []string {
	var result []string

	for _, cfg := range widgets {
		w, ok := Get(cfg.Name)
		if !ok {
			continue
		}

		if !w.ShouldRender(session, &cfg) {
			continue
		}

		rendered := w.Render(session, &cfg)
		if rendered != "" {
			result = append(result, rendered)
		}
	}

	return result
}

// contextSparkWidget holds the singleton instance for history injection.
var contextSparkWidget = &ContextSparkWidget{}

// blockTimerWidget holds the singleton instance for history injection.
var blockTimerWidget = &BlockTimerWidget{}

// toolsWidget holds the singleton instance for transcript injection.
var toolsWidget = &ToolsWidget{}

// agentsWidget holds the singleton instance for transcript injection.
var agentsWidget = &AgentsWidget{}

// SetHistory sets the history on widgets that need it.
func SetHistory(h *history.History) {
	contextSparkWidget.SetHistory(h)
	blockTimerWidget.SetHistory(h)
}

// SetTranscript sets the transcript data on widgets that need it.
func SetTranscript(t *transcript.Data) {
	toolsWidget.SetTranscript(t)
	agentsWidget.SetTranscript(t)
}

func init() {
	// Register all built-in widgets
	Register(&ModelWidget{})
	Register(&ContextWidget{})
	Register(&GitWidget{})
	Register(&CostWidget{})
	Register(&CacheHitWidget{})
	Register(&APILatencyWidget{})
	Register(&CodeChangesWidget{})
	Register(&BurnRateWidget{})
	Register(&CompactETAWidget{})
	Register(contextSparkWidget)
	Register(blockTimerWidget)
	Register(toolsWidget)
	Register(agentsWidget)
}
