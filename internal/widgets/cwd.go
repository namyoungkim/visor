package widgets

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
)

// CWDWidget displays the current working directory path.
//
// Supported Extra options:
//   - show_label: "true"/"false" - show "CWD:" prefix (default: false)
//   - show_basename: "true"/"false" - show only directory name (default: false)
//   - max_length: maximum path length, 0 = full (default: 0). Truncates with "…/" prefix.
//
// Output format: "~/project/visor" or "CWD: visor"
type CWDWidget struct{}

func (w *CWDWidget) Name() string {
	return "cwd"
}

func (w *CWDWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	cwd := session.CWD
	if cwd == "" {
		return ""
	}

	showBasename := GetExtraBool(cfg, "show_basename", false)
	maxLen := GetExtraInt(cfg, "max_length", 0)

	var display string
	if showBasename {
		display = filepath.Base(cwd)
	} else {
		display = abbreviateHome(cwd)
	}

	if maxLen > 0 && len([]rune(display)) > maxLen {
		display = truncatePath(display, maxLen)
	}

	if GetExtraBool(cfg, "show_label", false) {
		display = "CWD: " + display
	}

	return render.Colorize(display, "cyan")
}

func (w *CWDWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return session.CWD != ""
}

// abbreviateHome replaces the home directory prefix with ~.
func abbreviateHome(path string) string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return path
	}
	if path == home {
		return "~"
	}
	if strings.HasPrefix(path, home+"/") {
		return "~" + path[len(home):]
	}
	return path
}

// truncatePath shortens a path to fit within maxLen (in runes) by replacing
// leading segments with "…/". Uses rune counting for proper Unicode handling.
func truncatePath(path string, maxLen int) string {
	runes := []rune(path)
	if len(runes) <= maxLen {
		return path
	}
	if maxLen < 2 {
		return string(runes[:maxLen])
	}
	// "…" takes 1 rune, so we have maxLen-1 runes for the tail
	tail := runes[len(runes)-(maxLen-1):]
	// Try to find a '/' boundary for cleaner truncation
	for i, r := range tail {
		if r == '/' {
			return "…" + string(tail[i:])
		}
	}
	return "…" + string(tail)
}
