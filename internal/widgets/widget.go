package widgets

import (
	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
)

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
