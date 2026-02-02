package widgets

import (
	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/history"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// Sparkline characters (8 levels: 0-12.5%, 12.5-25%, ..., 87.5-100%).
var sparkChars = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// ContextSparkWidget displays a sparkline of recent context usage.
//
// Supported Extra options:
//   - width: number of characters in sparkline (default: "8")
//   - show_label: "true"/"false" - show prefix (default: false)
type ContextSparkWidget struct {
	history *history.History
}

func (w *ContextSparkWidget) Name() string {
	return "context_spark"
}

// SetHistory sets the history for this widget.
func (w *ContextSparkWidget) SetHistory(h *history.History) {
	w.history = h
}

func (w *ContextSparkWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	if w.history == nil {
		return render.Colorize("—", "dim")
	}

	width := GetExtraInt(cfg, "width", 8)
	values := w.history.GetContextHistory(width)

	if len(values) < 2 {
		// Need at least 2 data points for a meaningful sparkline
		return render.Colorize("—", "dim")
	}

	// Build sparkline
	spark := sparkline(values)

	var text string
	if cfg.Format != "" {
		text = FormatOutput(cfg, "", spark)
	} else if GetExtraBool(cfg, "show_label", false) {
		text = "Ctx: " + spark
	} else {
		text = spark
	}

	// Color based on trend (last value compared to average)
	color := sparkColor(values)
	return render.Colorize(text, color)
}

func (w *ContextSparkWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	// Only render if we have history data
	return w.history != nil && w.history.Count() >= 2
}

// sparkline converts values (0-100) to sparkline characters.
func sparkline(values []float64) string {
	if len(values) == 0 {
		return ""
	}

	result := make([]rune, len(values))
	for i, v := range values {
		// Clamp to 0-100
		if v < 0 {
			v = 0
		} else if v > 100 {
			v = 100
		}

		// Map to spark character index (0-7)
		idx := int(v / 100.0 * float64(len(sparkChars)-1))
		if idx >= len(sparkChars) {
			idx = len(sparkChars) - 1
		}
		result[i] = sparkChars[idx]
	}

	return string(result)
}

// sparkColor determines color based on trend.
func sparkColor(values []float64) string {
	if len(values) < 2 {
		return "white"
	}

	// Compare last value to previous average
	last := values[len(values)-1]
	sum := 0.0
	for i := 0; i < len(values)-1; i++ {
		sum += values[i]
	}
	avg := sum / float64(len(values)-1)

	// Rising trend (context filling up) = more urgent
	if last > avg+5 {
		return "red"
	} else if last < avg-5 {
		return "green"
	}
	return "yellow"
}
