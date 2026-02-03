package widgets

import (
	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/cost"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// Block cost thresholds (USD) for 5-hour block.
const (
	BlockCostWarningUSD  = 2.0
	BlockCostCriticalUSD = 5.0
)

// BlockCostWidget displays the cost spent in the current 5-hour block.
// This is useful for Claude Pro users to track spending within rate limit windows.
//
// Supported Extra options:
//   - show_label: "true"/"false" - whether to show "Block$:" prefix (default: false)
//   - warn_threshold: "2.0" - USD for warning color (default: 2.0)
//   - critical_threshold: "5.0" - USD for critical/red color (default: 5.0)
type BlockCostWidget struct {
	costData *cost.CostData
}

func (w *BlockCostWidget) Name() string {
	return "block_cost"
}

// SetCostData sets the cost data for this widget.
func (w *BlockCostWidget) SetCostData(data *cost.CostData) {
	w.costData = data
}

func (w *BlockCostWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	if w.costData == nil {
		return render.Colorize("—", "gray")
	}

	value := formatCost(w.costData.FiveHourBlock)

	var text string
	if cfg.Format != "" {
		text = FormatOutput(cfg, "", value)
	} else if GetExtraBool(cfg, "show_label", false) {
		text = "Block$: " + value
	} else {
		text = value
	}

	warnThreshold := GetExtraFloat(cfg, "warn_threshold", BlockCostWarningUSD)
	criticalThreshold := GetExtraFloat(cfg, "critical_threshold", BlockCostCriticalUSD)
	color := ColorByThreshold(w.costData.FiveHourBlock, warnThreshold, criticalThreshold)
	return render.Colorize(text, color)
}

func (w *BlockCostWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return w.costData != nil && !w.costData.BlockStartTime.IsZero()
}
