package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// Token speed thresholds (tokens per second - lower is worse)
const (
	TokenSpeedWarningTPS  = 20.0
	TokenSpeedCriticalTPS = 10.0
)

// TokenSpeedWidget displays the output token generation speed.
//
// Supported Extra options:
//   - show_label: "true"/"false" - show "out:" prefix (default: false)
//   - warn_threshold: tokens/sec below which to show warning (default: 20)
//   - critical_threshold: tokens/sec below which to show critical (default: 10)
//
// Output format: "42.1 tok/s"
type TokenSpeedWidget struct{}

func (w *TokenSpeedWidget) Name() string {
	return "token_speed"
}

func (w *TokenSpeedWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	tokens := session.GetTotalOutputTokens()
	durationMs := session.Cost.TotalAPIDurationMs

	if durationMs <= 0 || tokens <= 0 {
		return "—"
	}

	durationSec := float64(durationMs) / 1000.0
	speed := float64(tokens) / durationSec

	var value string
	if speed >= 100 {
		value = fmt.Sprintf("%.0f tok/s", speed)
	} else if speed >= 10 {
		value = fmt.Sprintf("%.1f tok/s", speed)
	} else {
		value = fmt.Sprintf("%.2f tok/s", speed)
	}

	var text string
	if GetExtraBool(cfg, "show_label", false) {
		text = "out: " + value
	} else {
		text = value
	}

	// Lower speed is worse (inverse threshold)
	warnThreshold := GetExtraFloat(cfg, "warn_threshold", TokenSpeedWarningTPS)
	criticalThreshold := GetExtraFloat(cfg, "critical_threshold", TokenSpeedCriticalTPS)
	color := colorByThresholdLowerIsWorse(speed, warnThreshold, criticalThreshold)

	return render.Colorize(text, color)
}

func (w *TokenSpeedWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return session.Cost.TotalAPIDurationMs > 0 && session.GetTotalOutputTokens() > 0
}

// colorByThresholdLowerIsWorse returns color where lower values are worse.
// Above warn = green, between warn and critical = yellow, below critical = red
func colorByThresholdLowerIsWorse(value, warn, critical float64) string {
	if value <= critical {
		return "red"
	} else if value <= warn {
		return "yellow"
	}
	return "green"
}
