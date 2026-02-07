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
		if maxLen > 0 && len(display) > maxLen {
			display = truncatePath(display, maxLen)
		}
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

// truncatePath shortens a path to fit within maxLen by replacing leading segments with "…/".
func truncatePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	// Minimum: "…/" + at least 1 char
	if maxLen < 3 {
		return path[:maxLen]
	}
	// Find the rightmost portion that fits within maxLen - 2 (for "…/")
	remaining := maxLen - 2 // len("…/") where … is 3 bytes but we count display chars
	// Walk from the end to find a '/' boundary
	for i := len(path) - remaining; i < len(path); i++ {
		if path[i] == '/' {
			result := "…" + path[i:]
			if len(result) <= maxLen+2 { // …is 3 bytes vs 1 display char, allow extra
				return result
			}
		}
	}
	// No good boundary found, just truncate from end
	return "…" + path[len(path)-remaining:]
}
