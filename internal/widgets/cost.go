package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// CostWidget displays the total API cost.
type CostWidget struct{}

func (w *CostWidget) Name() string {
	return "cost"
}

func (w *CostWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	cost := session.Cost.TotalCostUSD

	var text string
	switch {
	case cost >= 0.01:
		text = fmt.Sprintf("$%.2f", cost)
	case cost > 0:
		text = fmt.Sprintf("$%.3f", cost)
	default:
		text = "$0.00"
	}

	color := ColorByThreshold(cost, CostWarningUSD, CostDangerUSD)
	return render.Colorize(text, color)
}

func (w *CostWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return true
}
