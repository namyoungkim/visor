package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// Compact ETA thresholds.
const (
	CompactETAWarningMin = 10.0 // Warning when <10 min to compact
	CompactETADangerMin  = 5.0  // Danger when <5 min to compact
	CompactThresholdPct  = 80.0 // Compact triggers at 80%
)

// CompactETAWidget displays estimated time until context compact (80%).
//
// Calculation: (80 - current%) / burn_rate_per_min
// where burn_rate_per_min = current_percentage / (total_duration_ms / 60000)
//
// Supported Extra options:
//   - show_when_above: context % threshold to start showing (default: "40")
//   - show_label: "true"/"false" - show "ETA:" prefix (default: false)
type CompactETAWidget struct{}

func (w *CompactETAWidget) Name() string {
	return "compact_eta"
}

func (w *CompactETAWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	pct := session.ContextWindow.UsedPercentage
	durationMs := session.Cost.TotalDurationMs

	// Already at or above compact threshold
	if pct >= CompactThresholdPct {
		return render.Colorize("compact soon", "red")
	}

	// Cannot estimate without duration
	if durationMs <= 0 {
		return render.Colorize("—", "dim")
	}

	// Calculate context burn rate (%/min)
	durationMin := float64(durationMs) / 60000.0
	if durationMin <= 0 || pct <= 0 {
		return render.Colorize("—", "dim")
	}

	burnRatePctPerMin := pct / durationMin
	if burnRatePctPerMin <= 0 {
		return render.Colorize("—", "dim")
	}

	// Estimate time to 80%
	remainingPct := CompactThresholdPct - pct
	etaMinutes := remainingPct / burnRatePctPerMin

	// Format output
	var value string
	if etaMinutes >= 60 {
		hours := int(etaMinutes / 60)
		mins := int(etaMinutes) % 60
		if mins > 0 {
			value = fmt.Sprintf("~%dh%dm", hours, mins)
		} else {
			value = fmt.Sprintf("~%dh", hours)
		}
	} else if etaMinutes >= 1 {
		value = fmt.Sprintf("~%dm", int(etaMinutes))
	} else {
		value = "<1m"
	}

	var text string
	if cfg.Format != "" {
		text = FormatOutput(cfg, "", value)
	} else if GetExtraBool(cfg, "show_label", false) {
		text = "ETA: " + value
	} else {
		text = value
	}

	// Color: closer to compact = more urgent
	var color string
	if etaMinutes <= CompactETADangerMin {
		color = "red"
	} else if etaMinutes <= CompactETAWarningMin {
		color = "yellow"
	} else {
		color = "green"
	}

	return render.Colorize(text, color)
}

func (w *CompactETAWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	// Get threshold from config (default: 40%)
	threshold := GetExtraInt(cfg, "show_when_above", 40)

	pct := session.ContextWindow.UsedPercentage
	durationMs := session.Cost.TotalDurationMs

	// Show only if above threshold and we have duration data
	return pct >= float64(threshold) && durationMs > 0
}
