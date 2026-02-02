# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**visor** (Claude Code Efficiency Dashboard) is a Go-based statusline for Claude Code focused on real-time efficiency metrics rather than just status display. The key differentiator is exposing cache hit rate, API latency, and code changes—data already in stdin JSON that no other project uses.

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
echo '{"model":{"display_name":"Opus"},...}' | ./visor

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

```bash
# Create annotated tag
git tag -a v0.1.x -m "v0.1.x: Brief description

Features:
- Feature 1
- Feature 2"

# Push tag to remote
git push origin v0.1.x

# List all tags
git tag -l
```

## Architecture

### Data Flow
```
stdin JSON → input.Parse() → Session struct
                                   │
config.Load() → Config             │
                  │                │
                  ▼                ▼
            widgets.Registry.RenderAll(session, config)
                                   │
                                   ▼
                          render.Layout() → stdout ANSI
```

### Project Structure (cmd/ + internal/)
```
cmd/visor/main.go           # CLI entry point only
internal/input/             # stdin JSON parsing → Session struct
internal/config/            # TOML config loading
internal/widgets/           # Widget interface + implementations
internal/render/            # Layout, ANSI colors, truncation
internal/git/               # git CLI wrapper
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
- **Git info**: External `git` CLI calls with 200ms timeout (zero dependencies)
- **Dependencies**: Only `BurntSushi/toml` for config parsing

## MVP Widgets (v0.1)

| Widget | Identifier | Unique? |
|--------|------------|---------|
| Model name | `model` | No |
| Context % + progress bar | `context` | No |
| Git status | `git` | No |
| Cost | `cost` | No |
| Cache hit rate | `cache_hit` | **Yes** |
| API latency | `api_latency` | **Yes** |
| Code changes | `code_changes` | **Yes** |

Cache hit rate formula: `cache_read_input_tokens / (cache_read + input_tokens) × 100`

## Config Options (v0.1.2)

- `[general].separator` - Widget separator (default: `" | "`)
- `context` widget extras: `show_bar`, `bar_width` for progress bar customization

## Performance Requirements

- Cold startup: < 5ms
- No panics on malformed/missing JSON fields
- Graceful fallback (empty output on parse failure)

## Related Documentation

- `docs/00_PRD.md` — Full product spec, widgets, phases
- `docs/01_IMPACT_MAPPING.md` — Goals and deliverables
- `docs/02_USER_STORY_MAPPING.md` — User journey and validation scenarios
- `docs/03_C4_MODEL.md` — System architecture diagrams
- `docs/04_ADR.md` — 7 architecture decisions with rationale
- `docs/05_IMPLEMENTATION.md` — Code structure, APIs, extension guide
- `docs/06_PROGRESS.md` — PRD progress tracking
