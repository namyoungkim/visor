package widgets

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// PlanWidget displays the detected Claude plan type.
//
// Detection logic:
//  1. Model ID contains "bedrock" → "Bedrock"
//  2. Model ID contains "vertex" → "Vertex"
//  3. OAuth credentials exist at ~/.claude/auth.json → "Pro" (cannot distinguish Max/Team)
//  4. No OAuth → "API"
//
// Output format: "Pro" or "API" or "Bedrock" or "Vertex"
type PlanWidget struct{}

func (w *PlanWidget) Name() string {
	return "plan"
}

func (w *PlanWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	modelID := strings.ToLower(session.Model.ID)

	// Check for cloud providers first
	if strings.Contains(modelID, "bedrock") {
		return render.Colorize("Bedrock", "magenta")
	}
	if strings.Contains(modelID, "vertex") {
		return render.Colorize("Vertex", "blue")
	}

	// Check for OAuth credentials (Pro/Max/Team users)
	if hasOAuthCredentials() {
		return render.Colorize("Pro", "cyan")
	}

	// Default to API
	return render.Colorize("API", "yellow")
}

func (w *PlanWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return true
}

// hasOAuthCredentials checks if OAuth credentials exist at ~/.claude/auth.json
func hasOAuthCredentials() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	authPath := filepath.Join(home, ".claude", "auth.json")
	info, err := os.Stat(authPath)
	if err != nil {
		return false
	}

	// File exists and has content
	return info.Size() > 2 // More than just "{}" or "[]"
}
