// Package claudeconfig provides utilities for parsing Claude Code configuration files.
package claudeconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Counts holds aggregated counts from Claude configuration.
type Counts struct {
	ClaudeMDCount int // CLAUDE.md files from cwd up to root
	RulesCount    int // Permission rules (permissions.allow entries)
	MCPCount      int // Enabled MCP plugins
	HooksCount    int // Configured hooks
}

// settingsJSON represents the structure of ~/.claude/settings.json
type settingsJSON struct {
	Permissions struct {
		Allow []interface{} `json:"allow"`
		Deny  []interface{} `json:"deny"`
	} `json:"permissions"`
	EnabledPlugins []string               `json:"enabledPlugins"`
	Hooks          map[string]interface{} `json:"hooks"`
}

// LoadCounts loads configuration counts from Claude config files.
// cwd is the current working directory to start searching for CLAUDE.md files.
func LoadCounts(cwd string) *Counts {
	counts := &Counts{}

	// Count CLAUDE.md files from cwd up to root
	counts.ClaudeMDCount = countClaudeMDFiles(cwd)

	// Parse ~/.claude/settings.json for rules, MCPs, hooks
	home, err := os.UserHomeDir()
	if err == nil {
		settingsPath := filepath.Join(home, ".claude", "settings.json")
		parseSettings(settingsPath, counts)
	}

	return counts
}

// countClaudeMDFiles counts CLAUDE.md files from dir up to filesystem root.
func countClaudeMDFiles(dir string) int {
	if dir == "" {
		return 0
	}

	count := 0
	current := dir

	for {
		// Check for CLAUDE.md in current directory
		claudeMD := filepath.Join(current, "CLAUDE.md")
		if _, err := os.Stat(claudeMD); err == nil {
			count++
		}

		// Also check for .claude.md (hidden variant)
		claudeMDHidden := filepath.Join(current, ".claude.md")
		if _, err := os.Stat(claudeMDHidden); err == nil {
			count++
		}

		// Move to parent directory
		parent := filepath.Dir(current)
		if parent == current {
			// Reached root
			break
		}
		current = parent
	}

	return count
}

// parseSettings parses settings.json and updates counts.
func parseSettings(path string, counts *Counts) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var settings settingsJSON
	if err := json.Unmarshal(data, &settings); err != nil {
		return
	}

	counts.RulesCount = len(settings.Permissions.Allow)
	counts.MCPCount = len(settings.EnabledPlugins)
	counts.HooksCount = countHooks(settings.Hooks)
}

// countHooks counts the total number of hook entries.
func countHooks(hooks map[string]interface{}) int {
	count := 0
	for _, v := range hooks {
		// Each hook type can have multiple entries (array)
		if arr, ok := v.([]interface{}); ok {
			count += len(arr)
		} else {
			count++
		}
	}
	return count
}
