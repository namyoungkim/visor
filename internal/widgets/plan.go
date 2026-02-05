package widgets

import (
	"os"
	"strings"

	"github.com/namyoungkim/visor/internal/auth"
	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// PlanWidget displays the detected Claude plan type.
//
// Detection logic:
//  1. Model ID contains "bedrock" → "Bedrock"
//  2. Model ID contains "vertex" → "Vertex"
//  3. ANTHROPIC_API_KEY is set → "API"
//  4. OAuth credentials available (DefaultProvider) → "Pro"
//  5. No API key and no bedrock/vertex → "Pro" (subscription assumed)
//
// Supported Extra options:
//   - show_label: "true"/"false" - whether to show "Plan:" prefix (default: false)
//
// Output format: "Pro" or "API" or "Bedrock" or "Vertex" (or "Plan: Pro" with show_label)
type PlanWidget struct{}

func (w *PlanWidget) Name() string {
	return "plan"
}

func (w *PlanWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	modelID := strings.ToLower(session.Model.ID)
	showLabel := GetExtraBool(cfg, "show_label", false)

	var text, color string
	switch {
	case strings.Contains(modelID, "bedrock"):
		text, color = "Bedrock", "magenta"
	case strings.Contains(modelID, "vertex"):
		text, color = "Vertex", "blue"
	case isAPIKeyUser():
		text, color = "API", "yellow"
	default:
		// No API key + no cloud provider = subscription user
		text, color = detectSubscriptionType()
	}

	if showLabel {
		text = "Plan: " + text
	}

	return render.Colorize(text, color)
}

func (w *PlanWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return true
}

// isAPIKeyUser checks if the user is using an API key.
func isAPIKeyUser() bool {
	return os.Getenv("ANTHROPIC_API_KEY") != ""
}

// detectSubscriptionType returns the subscription display name and color.
// Checks OAuth credentials for subscriptionType field.
func detectSubscriptionType() (string, string) {
	provider := auth.DefaultProvider()
	creds, err := provider.Get()
	if err != nil {
		return "Pro", "cyan" // Default assumption for non-API users
	}

	switch strings.ToLower(creds.SubscriptionType) {
	case "max":
		return "Max", "cyan"
	case "team":
		return "Team", "cyan"
	case "pro":
		return "Pro", "cyan"
	default:
		return "Pro", "cyan"
	}
}
