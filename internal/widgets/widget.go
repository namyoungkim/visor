package widgets

import (
	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
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

func init() {
	// Register all built-in widgets
	Register(&ModelWidget{})
	Register(&ContextWidget{})
	Register(&GitWidget{})
	Register(&CostWidget{})
	Register(&CacheHitWidget{})
	Register(&APILatencyWidget{})
	Register(&CodeChangesWidget{})
}
