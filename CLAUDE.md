# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**visor** (Claude Code Efficiency Dashboard) is a Go-based statusline for Claude Code focused on real-time efficiency metrics rather than just status display. The key differentiator is exposing cache hit rate, API latency, burn rate, and context prediction—data already in stdin JSON that no other project uses.

## Environment Setup

Go 1.22+ is required. If using mise (recommended):
```bash
# Activate mise to add go to PATH
eval "$(mise activate bash)"  # or zsh

# Or run go directly via mise
mise exec go -- build ./cmd/visor

# Check go version
go version  # should be 1.22+
```

## Build & Development Commands

```bash
# Build
go build -o visor ./cmd/visor

# Run tests
go test ./...

# Manual testing
echo '{"session_id":"test","model":{"display_name":"Opus"},"context_window":{"used_percentage":42.5},"cost":{"total_cost_usd":0.48,"total_duration_ms":45000}}' | ./visor

# CLI flags
./visor --version   # Version info
./visor --init      # Generate ~/.config/visor/config.toml
./visor --setup     # Configure Claude Code statusline
./visor --check     # Validate config file
./visor --debug     # Debug output to stderr

# Install globally
go install github.com/namyoungkim/visor@latest
```

## Release

Releases are automated via GitHub Actions. When you push a version tag, GoReleaser builds binaries for Linux/macOS (amd64/arm64) and creates a GitHub Release.

```bash
# Create annotated tag
git tag -a v0.4.0 -m "v0.4.0: Brief description

Features:
- Feature 1
- Feature 2"

# Push tag to remote (triggers GitHub Actions release)
git push origin v0.4.0

# List all tags
git tag -l

# Local build with version
go build -ldflags "-X main.version=0.4.0" -o visor ./cmd/visor
```

## Architecture

### Data Flow
```
stdin JSON → input.Parse() → Session struct
                                   │
config.Load() → Config             │
                  │                │
history.Load() → History           │
                  │                │
transcript.Parse() → Data          │  (v0.3)
                  │                │
                  ▼                ▼
            widgets.RenderAll(session, config)
                      │
                      ▼
            render.Layout() or render.SplitLayout()
                      │
                      ▼
                  stdout ANSI
```

### Project Structure (cmd/ + internal/)
```
cmd/visor/main.go           # CLI entry point only
internal/input/             # stdin JSON parsing → Session struct
internal/config/            # TOML config loading
internal/widgets/           # Widget interface + implementations
internal/render/            # Layout, ANSI colors, truncation
internal/git/               # git CLI wrapper
internal/history/           # Session history buffer
internal/transcript/        # JSONL transcript parsing (v0.3)
```

### Widget Interface Pattern
All widgets implement this interface. Add new widgets by creating a file in `internal/widgets/` and registering in the Registry:

```go
type Widget interface {
    Name() string
    Render(session *Session, cfg *WidgetConfig) string
    ShouldRender(session *Session, cfg *WidgetConfig) bool
}
```

## Key Technical Decisions

- **Language**: Go (1-2ms startup, fills empty niche in ecosystem)
- **Config**: TOML at `~/.config/visor/config.toml` (uses `[[line]]` for multiline layout)
- **History**: JSON at `~/.cache/visor/history_<session_id>.json`
- **Git info**: External `git` CLI calls with 200ms timeout (zero dependencies)
- **Dependencies**: Only `BurntSushi/toml` for config parsing

## Widgets (v0.4.0)

### Core Widgets (v0.1)
| Widget | Identifier | Unique? |
|--------|------------|---------|
| Model name | `model` | No |
| Context % + progress bar | `context` | No |
| Git status | `git` | No |
| Cost | `cost` | No |
| Cache hit rate | `cache_hit` | **Yes** |
| API latency | `api_latency` | **Yes** |
| Code changes | `code_changes` | **Yes** |

### Efficiency Widgets (v0.2)
| Widget | Identifier | Output Example | Unique? |
|--------|------------|----------------|---------|
| Burn rate | `burn_rate` | `64.0¢/min` | **Yes** |
| Compact ETA | `compact_eta` | `~18m` | **Yes** |
| Context sparkline | `context_spark` | `▂▃▄▅▆` | **Yes** |

### Transcript Widgets (v0.3)
| Widget | Identifier | Output Example | Unique? |
|--------|------------|----------------|---------|
| Tool status | `tools` | `✓Read ✓Write ◐Bash` | **Yes** |
| Agent status | `agents` | `✓Plan ◐Explore` | **Yes** |

### Rate Limit Widget (v0.4)
| Widget | Identifier | Output Example | Unique? |
|--------|------------|----------------|---------|
| Block timer | `block_timer` | `Block: 4h23m` | **Yes** |

### Widget Formulas
- Cache hit rate: `cache_read_tokens / (cache_read + input_tokens) × 100`
- Burn rate: `total_cost_usd / (total_duration_ms / 60000)`
- Compact ETA: `(80 - current%) / context_burn_rate_per_min`
- Block timer: Remaining time in 5-hour Claude Pro rate limit block

## Config Options (v0.4.0)

### General
- `[general].separator` - Widget separator (default: `" | "`)

### Layout Types
```toml
# Single-line layout
[[line]]
  [[line.widget]]
  name = "model"

# Split layout (left/right aligned)
[[line]]
  [[line.left]]
  name = "model"
  [[line.right]]
  name = "cost"
```

### Widget Extras
- `context`: `show_bar`, `bar_width`, `warn_threshold` (60), `critical_threshold` (80)
- `compact_eta`: `show_when_above` (default: 40)
- `context_spark`: `width` (default: 8)
- `burn_rate`: `show_label`, `warn_threshold` (10), `critical_threshold` (25) (cents/min)
- `cost`: `show_label`, `warn_threshold` (0.5), `critical_threshold` (1.0) (USD)
- `cache_hit`: `show_label`, `good_threshold` (80), `warn_threshold` (50)
- `api_latency`: `warn_threshold` (2000), `critical_threshold` (5000) (ms)
- `block_timer`: `show_label`, `warn_threshold` (80), `critical_threshold` (95) (% elapsed)
- `tools`: `max_display` (default: 3), `show_label`
- `agents`: `max_display` (default: 3), `show_label`

## Performance Requirements

- Cold startup: < 5ms
- No panics on malformed/missing JSON fields
- Graceful fallback (empty output on parse failure)

## Related Documentation

- `docs/00_PRD.md` — Full product spec, widgets, phases
- `docs/01_IMPACT_MAPPING.md` — Goals and deliverables
- `docs/02_USER_STORY_MAPPING.md` — User journey and validation scenarios
- `docs/03_C4_MODEL.md` — System architecture diagrams
- `docs/04_ADR.md` — Architecture decisions with rationale
- `docs/05_IMPLEMENTATION.md` — Code structure, APIs, extension guide
- `docs/06_PROGRESS.md` — PRD progress tracking
