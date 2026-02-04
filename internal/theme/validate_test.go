package theme

import "testing"

func TestValidateColor(t *testing.T) {
	tests := []struct {
		name  string
		color string
		want  bool
	}{
		// Empty (use preset value)
		{"empty string", "", true},

		// Hex colors - 3 digits
		{"hex 3 digits lowercase", "#abc", true},
		{"hex 3 digits uppercase", "#ABC", true},
		{"hex 3 digits mixed", "#aBc", true},

		// Hex colors - 6 digits
		{"hex 6 digits lowercase", "#aabbcc", true},
		{"hex 6 digits uppercase", "#AABBCC", true},
		{"hex 6 digits mixed", "#AaBbCc", true},
		{"hex 6 digits gruvbox", "#fabd2f", true},
		{"hex 6 digits nord", "#81a1c1", true},

		// Hex colors - 8 digits (with alpha)
		{"hex 8 digits", "#aabbccdd", true},
		{"hex 8 digits uppercase", "#AABBCCDD", true},

		// Invalid hex colors
		{"hex 1 digit", "#a", false},
		{"hex 2 digits", "#ab", false},
		{"hex 4 digits", "#abcd", false},
		{"hex 5 digits", "#abcde", false},
		{"hex 7 digits", "#abcdef0", false},
		{"hex without hash", "aabbcc", false},
		{"hex invalid char", "#gggggg", false},

		// Named colors
		{"named black", "black", true},
		{"named red", "red", true},
		{"named green", "green", true},
		{"named yellow", "yellow", true},
		{"named blue", "blue", true},
		{"named magenta", "magenta", true},
		{"named cyan", "cyan", true},
		{"named white", "white", true},
		{"named gray", "gray", true},
		{"named grey", "grey", true},

		// Named colors - case insensitive
		{"named uppercase", "RED", true},
		{"named mixed case", "Blue", true},
		{"named CYAN", "CYAN", true},

		// Bright named colors
		{"bright black", "brightblack", true},
		{"bright red", "brightred", true},
		{"bright white", "brightwhite", true},

		// Invalid named colors
		{"invalid name", "orange", false},
		{"invalid name purple", "purple", false},
		{"invalid name pink", "pink", false},
		{"random string", "notacolor", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateColor(tt.color)
			if got != tt.want {
				t.Errorf("ValidateColor(%q) = %v, want %v", tt.color, got, tt.want)
			}
		})
	}
}
