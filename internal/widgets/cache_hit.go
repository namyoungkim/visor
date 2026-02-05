package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// CacheHitWidget displays the cache hit rate.
// This is a unique metric that no other statusline exposes.
// Formula: cache_read_input_tokens / (cache_read_input_tokens + input_tokens) * 100
//
// Supported Extra options:
//   - show_label: "true"/"false" - whether to show "Cache:" prefix (default: true)
//   - good_threshold: "80" - percentage for good/green color (default: 80)
//   - warn_threshold: "50" - percentage for warning color (default: 50)
type CacheHitWidget struct{}

func (w *CacheHitWidget) Name() string {
	return "cache_hit"
}

func (w *CacheHitWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	// Check if current_usage is available
	cu := session.GetCurrentUsage()
	if cu == nil {
		label := "Cache: —"
		if !GetExtraBool(cfg, "show_label", true) {
			label = "—"
		}
		return render.Colorize(label, "gray")
	}

	cacheRead := cu.GetCacheReadTokens()
	inputTokens := cu.InputTokens

	total := cacheRead + inputTokens
	if total == 0 {
		label := "Cache: —"
		if !GetExtraBool(cfg, "show_label", true) {
			label = "—"
		}
		return render.Colorize(label, "gray")
	}

	rate := float64(cacheRead) / float64(total) * 100
	goodThreshold := GetExtraFloat(cfg, "good_threshold", CacheHitGoodPct)
	warnThreshold := GetExtraFloat(cfg, "warn_threshold", CacheHitWarningPct)
	color := ColorByThresholdInverse(rate, goodThreshold, warnThreshold)

	value := fmt.Sprintf("%.0f%%", rate)

	var text string
	if cfg.Format != "" {
		text = FormatOutput(cfg, "", value)
	} else if GetExtraBool(cfg, "show_label", true) {
		text = "Cache: " + value
	} else {
		text = value
	}

	return render.Colorize(text, color)
}

func (w *CacheHitWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return true
}
