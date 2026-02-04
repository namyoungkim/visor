package widgets

import (
	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// SessionIDWidget displays the current session ID.
//
// Supported Extra options:
//   - show_label: "true"/"false" - show "Session:" prefix (default: false)
//   - max_length: maximum length to display (default: 8, 0 = full)
//
// Output format: "a1b2c3d4" or "Session: a1b2c3d4"
type SessionIDWidget struct{}

func (w *SessionIDWidget) Name() string {
	return "session_id"
}

func (w *SessionIDWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	id := session.SessionID
	if id == "" {
		return ""
	}

	maxLen := GetExtraInt(cfg, "max_length", 8)
	if maxLen > 0 && len(id) > maxLen {
		id = id[:maxLen]
	}

	var text string
	if GetExtraBool(cfg, "show_label", false) {
		text = "Session: " + id
	} else {
		text = id
	}

	return render.Colorize(text, "gray")
}

func (w *SessionIDWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return session.SessionID != ""
}
