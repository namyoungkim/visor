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

# Run tests (CGO_ENABLED=0 required on macOS to avoid dyld LC_UUID errors)
CGO_ENABLED=0 go test ./...

# Manual testing
echo '{"session_id":"test","model":{"display_name":"Opus"},"context_window":{"used_percentage":42.5},"cost":{"total_cost_usd":0.48,"total_duration_ms":45000}}' | ./visor

# CLI flags
./visor --version         # Version info
./visor --init            # Generate config with 'default' preset
./visor --init minimal    # Generate config with specific preset
./visor --init help       # Show available presets
./visor --setup           # Configure Claude Code statusline
./visor --check           # Validate config file
./visor --debug           # Debug output to stderr
./visor --tui             # Interactive TUI config editor

# Install globally
go install github.com/namyoungkim/visor@latest
```

## Release

Releases are automated via GitHub Actions. When you push a version tag, GoReleaser builds binaries for Linux/macOS (amd64/arm64) and creates a GitHub Release.

```bash
# Create annotated tag
git tag -a v0.8.0 -m "v0.8.0: Brief description

Features:
- Feature 1
- Feature 2"

# Push tag to remote (triggers GitHub Actions release)
git push origin v0.8.0

# List all tags
git tag -l

# Local build with version
go build -ldflags "-X main.version=0.8.0" -o visor ./cmd/visor
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
internal/config/            # TOML config loading + saving
internal/widgets/           # Widget interface + implementations
internal/render/            # Layout, ANSI colors, truncation
internal/git/               # git CLI wrapper
internal/history/           # Session history buffer
internal/transcript/        # JSONL transcript parsing (v0.3)
internal/tui/               # Interactive TUI config editor (v0.5)
internal/theme/             # Theme presets and management (v0.6)
internal/cost/              # JSONL parsing and cost aggregation (v0.6)
internal/auth/              # OAuth credential providers (v0.6)
internal/usage/             # Usage limit API client (v0.6)
internal/claudeconfig/      # Claude config parsing (v0.10)
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
- **Global block state**: JSON at `~/.cache/visor/block_state.json` (v0.11.5, shared across sessions)
- **Git info**: External `git` CLI calls with 200ms timeout (zero dependencies)
- **TUI**: Charm ecosystem (bubbletea, bubbles, lipgloss) for interactive config editor
- **Dependencies**: `BurntSushi/toml`, `charmbracelet/bubbletea`, `charmbracelet/bubbles`, `charmbracelet/lipgloss`

## Widgets (v0.11.5)

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

### Transcript Widgets (v0.3, v0.7, v0.10)
| Widget | Identifier | Output Example | Unique? |
|--------|------------|----------------|---------|
| Tool status | `tools` | `✓Bash ×7 \| ✓Edit ×4 \| ✓Read ×6` | **Yes** |
| Agent status | `agents` | `✓Explore: Analyze... (42s)` or `◐Plan: Impl... (5s...)` | **Yes** |
| Task progress | `todos` | `⊙ Task name (3/5)` or `✓ All done (5/5)` | **Yes** |

### Rate Limit Widget (v0.4)
| Widget | Identifier | Output Example | Unique? |
|--------|------------|----------------|---------|
| Block timer | `block_timer` | `Block: 4h23m` | **Yes** |

### Cost Tracking Widgets (v0.6)
| Widget | Identifier | Output Example | Unique? |
|--------|------------|----------------|---------|
| Daily cost | `daily_cost` | `$2.34 today` | **Yes** |
| Weekly cost | `weekly_cost` | `$15.67 week` | **Yes** |
| Block cost | `block_cost` | `$0.45 block` | **Yes** |
| 5-hour limit | `block_limit` | `5h: 42%` | **Yes** |
| 7-day limit | `week_limit` | `7d: 69%` | **Yes** |

### Session Info Widgets (v0.10, 6 widgets)
| Widget | Identifier | Output Example | Unique? |
|--------|------------|----------------|---------|
| Session ID | `session_id` | `abc123de` | No |
| Session duration | `duration` | `⏱️ 5m` or `1h23m` | No |
| Token speed | `token_speed` | `42.1 tok/s` | **Yes** |
| Plan type | `plan` | `Pro` or `API` or `Bedrock` | No |
| Task progress | `todos` | `⊙ Task (3/5)` or `✓ All done (5/5)` | **Yes** |
| Config counts | `config_counts` | `2📄 3🔒 2🔌 1🪝` | **Yes** |
| Working directory | `cwd` | `~/project/visor` | No |

### Widget Formulas
- Cache hit rate: `cache_read_input_tokens / (cache_read_input_tokens + input_tokens) × 100`
- API latency (per-call): `total_api_duration_ms / total_api_calls`
- Burn rate: `total_cost_usd / (total_duration_ms / 60000)`
- Compact ETA: `(80 - current%) / context_burn_rate_per_min`
- Block timer: Remaining time in 5-hour Claude Pro rate limit block
- Token speed: `total_output_tokens / (total_api_duration_ms / 1000)`

## TUI Config Editor (v0.5)

Interactive terminal UI for configuration editing:

```bash
./visor --tui
```

### Keybindings
| Key | Action |
|-----|--------|
| `j/k` | Move cursor |
| `a` | Add widget |
| `d` | Delete widget |
| `e` | Edit widget options |
| `J/K` | Reorder widgets |
| `L` | Change layout (single/split) |
| `n` | Add new line |
| `s` | Save |
| `t` | Change theme |
| `q` | Quit |

### TUI Package Structure
```
internal/tui/
├── tui.go              # Run() entry point
├── model.go            # Model struct (state)
├── update.go           # Update() message handling
├── view.go             # View() rendering
├── styles.go           # lipgloss styles
├── keys.go             # Keybinding definitions
├── widget_options.go   # Widget option metadata
└── preview.go          # Sample session & preview
```

## Themes (v0.6, v0.8)

visor supports multiple theme presets:

| Theme | Description |
|-------|-------------|
| `default` | Standard ASCII separators |
| `powerline` | Powerline glyphs (, ) |
| `gruvbox` | Gruvbox color palette |
| `nord` | Nord color palette |
| `gruvbox-powerline` | Gruvbox + Powerline |
| `nord-powerline` | Nord + Powerline |

Theme picker in TUI: press `t` key.

### Custom Theme Overrides (v0.8)

Override preset colors and separators:

```toml
[theme]
name = "gruvbox"       # Base preset
powerline = true       # Enable Powerline style

[theme.colors]
warning = "#ff00ff"    # Hex color
critical = "red"       # Named color
backgrounds = ["#111111", "#222222"]

[theme.separators]
left = " :: "
right = " :: "
```

**Supported color formats:**
- Hex: `#RGB`, `#RRGGBB`, `#RRGGBBAA`
- Named: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`, `gray`
- Bright: `brightred`, `brightgreen`, etc.

**Color fields:** `normal`, `warning`, `critical`, `good`, `primary`, `secondary`, `muted`, `backgrounds`

**Separator fields:** `left`, `right`, `left_soft`, `right_soft`, `left_hard`, `right_hard`

## Config Options (v0.8.0)

### General
- `[general].separator` - Widget separator (default: `" | "`)
- `[general].debug` - Enable debug output to stderr (default: `false`). Can also use `--debug` CLI flag.

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
- `api_latency`: `warn_threshold` (2000), `critical_threshold` (5000) (ms, per-call average)
- `block_timer`: `show_label`, `warn_threshold` (80), `critical_threshold` (95) (% elapsed)
- `block_limit`: `show_label`, `show_remaining`, `show_bar` (false), `bar_width` (10), `warn_threshold` (70), `critical_threshold` (90)
- `tools`: `max_display` (default: 0, unlimited), `show_label`, `show_count` (default: true)
- `agents`: `max_display` (default: 2), `show_label`, `show_description` (default: true), `show_duration` (default: true), `max_description_len` (default: 15)
- `duration`: `show_icon` (default: true) - show ⏱️ prefix
- `token_speed`: `show_label`, `warn_threshold` (20), `critical_threshold` (10) (tokens/sec, lower is worse)
- `todos`: `show_label`, `max_subject_len` (default: 30)
- `config_counts`: `show_claude_md` (default: true), `show_rules` (default: true), `show_mcps` (default: true), `show_hooks` (default: true)
- `plan`: `show_label` (default: false) - show "Plan:" prefix
- `session_id`: `show_label`, `max_length` (default: 0, 0 = full)
- `cwd`: `show_label`, `show_basename` (default: false, show only directory name), `max_length` (default: 0, 0 = auto ~width/3)

## Configuration Presets (v0.11)

Initialize config with presets for different use cases:

```bash
./visor --init            # 'default' preset
./visor --init minimal    # Essential 4 widgets
./visor --init efficiency # Cost optimization focus
./visor --init developer  # Tool/agent monitoring
./visor --init pro        # Claude Pro/Max rate limits
./visor --init full       # All 22 widgets, multi-line
./visor --init help       # List available presets
```

| Preset | Widgets | Description |
|--------|---------|-------------|
| `minimal` | 4 | model, context, cost, git |
| `default` | 6 | model, context, cache_hit, api_latency, cost, git |
| `efficiency` | 6 | model, context, burn_rate, cache_hit, compact_eta, cost |
| `developer` | 7 | model, context, tools, agents, todos, code_changes, git |
| `pro` | 6 | model, context, block_limit, week_limit, daily_cost, cost |
| `full` | 23 | All widgets in 7 lines by category |

Rate limit widgets (`block_limit`, `week_limit`) count user turns (`type="user"` + `isMeta=false`) for message limits. Tier auto-detection from OAuth credentials (Pro: 45/5h, Max 5x: 225/5h, Max 20x: 900/5h).

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
