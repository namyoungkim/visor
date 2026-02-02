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
type CacheHitWidget struct{}

func (w *CacheHitWidget) Name() string {
	return "cache_hit"
}

func (w *CacheHitWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	// Check if current_usage is available
	if session.CurrentUsage == nil {
		return render.Colorize("Cache: —", "gray")
	}

	cacheRead := session.CurrentUsage.CacheReadTokens
	inputTokens := session.CurrentUsage.InputTokens

	total := cacheRead + inputTokens
	if total == 0 {
		return render.Colorize("Cache: —", "gray")
	}

	rate := float64(cacheRead) / float64(total) * 100

	// Color based on cache hit rate
	color := "red"
	if rate >= 80 {
		color = "green"
	} else if rate >= 50 {
		color = "yellow"
	}

	text := fmt.Sprintf("Cache: %.0f%%", rate)
	return render.Colorize(text, color)
}

func (w *CacheHitWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return true
}
