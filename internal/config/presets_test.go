package config

import (
	"strings"
	"testing"
)

func TestGetPreset(t *testing.T) {
	tests := []struct {
		name     string
		preset   string
		wantOK   bool
		wantName string
	}{
		{"minimal exists", "minimal", true, "minimal"},
		{"default exists", "default", true, "default"},
		{"efficiency exists", "efficiency", true, "efficiency"},
		{"developer exists", "developer", true, "developer"},
		{"pro exists", "pro", true, "pro"},
		{"full exists", "full", true, "full"},
		{"unknown preset", "unknown", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, ok := GetPreset(tt.preset)
			if ok != tt.wantOK {
				t.Errorf("GetPreset(%q) ok = %v, want %v", tt.preset, ok, tt.wantOK)
			}
			if ok && p.Name != tt.wantName {
				t.Errorf("GetPreset(%q).Name = %q, want %q", tt.preset, p.Name, tt.wantName)
			}
		})
	}
}

func TestListPresets(t *testing.T) {
	output := ListPresets()

	// Check all presets are listed
	for _, name := range PresetOrder {
		if !strings.Contains(output, name) {
			t.Errorf("ListPresets() should contain %q", name)
		}
	}

	// Check usage instructions
	if !strings.Contains(output, "visor --init") {
		t.Error("ListPresets() should contain usage instructions")
	}
}

func TestGetPresetTOML(t *testing.T) {
	tests := []struct {
		name       string
		preset     string
		wantErr    bool
		wantInTOML []string
	}{
		{
			name:    "minimal preset TOML",
			preset:  "minimal",
			wantErr: false,
			wantInTOML: []string{
				"# Preset: minimal",
				`name = "model"`,
				`name = "context"`,
				`name = "cost"`,
				`name = "git"`,
			},
		},
		{
			name:    "full preset has multiple lines",
			preset:  "full",
			wantErr: false,
			wantInTOML: []string{
				"[[line]]",
				`name = "cache_hit"`,
				`name = "block_timer"`,
			},
		},
		{
			name:    "unknown preset returns error",
			preset:  "unknown",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toml, err := GetPresetTOML(tt.preset)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPresetTOML(%q) error = %v, wantErr %v", tt.preset, err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			for _, want := range tt.wantInTOML {
				if !strings.Contains(toml, want) {
					t.Errorf("GetPresetTOML(%q) should contain %q", tt.preset, want)
				}
			}
		})
	}
}

func TestPresetWidgetCounts(t *testing.T) {
	tests := []struct {
		preset    string
		wantCount int
	}{
		{"minimal", 4},
		{"default", 6},
		{"efficiency", 6},
		{"developer", 6},
		{"pro", 6},
		{"full", 18},
	}

	for _, tt := range tests {
		t.Run(tt.preset, func(t *testing.T) {
			p, ok := GetPreset(tt.preset)
			if !ok {
				t.Fatalf("Preset %q not found", tt.preset)
			}

			count := 0
			for _, line := range p.Lines {
				count += len(line)
			}

			if count != tt.wantCount {
				t.Errorf("Preset %q has %d widgets, want %d", tt.preset, count, tt.wantCount)
			}
		})
	}
}
