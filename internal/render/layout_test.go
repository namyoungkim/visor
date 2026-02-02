package render

import (
	"strings"
	"testing"
)

const testSeparator = " | "

func TestLayout_SingleWidget(t *testing.T) {
	widgets := []string{"model"}
	result := Layout(widgets, testSeparator)

	if result != "model" {
		t.Errorf("Expected 'model', got '%s'", result)
	}
}

func TestLayout_MultipleWidgets(t *testing.T) {
	widgets := []string{"model", "context", "cost"}
	result := Layout(widgets, testSeparator)

	expected := "model | context | cost"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestLayout_EmptyWidgets(t *testing.T) {
	widgets := []string{}
	result := Layout(widgets, testSeparator)

	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestLayout_WithEmptyStrings(t *testing.T) {
	widgets := []string{"model", "", "cost", ""}
	result := Layout(widgets, testSeparator)

	expected := "model | cost"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestLayout_CustomSeparator(t *testing.T) {
	widgets := []string{"a", "b", "c"}

	tests := []struct {
		separator string
		expected  string
	}{
		{" ", "a b c"},
		{" | ", "a | b | c"},
		{" :: ", "a :: b :: c"},
		{"", "abc"},
	}

	for _, tt := range tests {
		result := Layout(widgets, tt.separator)
		if result != tt.expected {
			t.Errorf("Separator %q: expected '%s', got '%s'", tt.separator, tt.expected, result)
		}
	}
}

func TestMultiLine_SingleLine(t *testing.T) {
	lines := [][]string{
		{"model", "context"},
	}
	result := MultiLine(lines, testSeparator)

	if strings.Contains(result, "\n") {
		t.Error("Expected single line output")
	}
}

func TestMultiLine_MultipleLines(t *testing.T) {
	lines := [][]string{
		{"model", "context"},
		{"cost", "git"},
	}
	result := MultiLine(lines, testSeparator)

	parts := strings.Split(result, "\n")
	if len(parts) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(parts))
	}
}

func TestMultiLine_EmptyLines(t *testing.T) {
	lines := [][]string{
		{"model"},
		{},
		{"cost"},
	}
	result := MultiLine(lines, testSeparator)

	// Empty lines should be filtered out
	parts := strings.Split(result, "\n")
	if len(parts) != 2 {
		t.Errorf("Expected 2 non-empty lines, got %d", len(parts))
	}
}

func TestMultiLine_AllEmpty(t *testing.T) {
	lines := [][]string{
		{},
		{},
	}
	result := MultiLine(lines, testSeparator)

	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestSplitLayout_BothSides(t *testing.T) {
	left := []string{"model", "git"}
	right := []string{"cost", "cache"}
	result := SplitLayout(left, right, testSeparator)

	// Should contain both sides
	if !strings.Contains(result, "model | git") {
		t.Errorf("Expected left side 'model | git', got '%s'", result)
	}
	if !strings.Contains(result, "cost | cache") {
		t.Errorf("Expected right side 'cost | cache', got '%s'", result)
	}

	// Right side should be after left side
	leftIdx := strings.Index(result, "model")
	rightIdx := strings.Index(result, "cost")
	if leftIdx >= rightIdx {
		t.Error("Expected left side to come before right side")
	}
}

func TestSplitLayout_LeftOnly(t *testing.T) {
	left := []string{"model", "git"}
	right := []string{}
	result := SplitLayout(left, right, testSeparator)

	if !strings.Contains(result, "model | git") {
		t.Errorf("Expected 'model | git', got '%s'", result)
	}
}

func TestSplitLayout_RightOnly(t *testing.T) {
	left := []string{}
	right := []string{"cost", "cache"}
	result := SplitLayout(left, right, testSeparator)

	if !strings.Contains(result, "cost | cache") {
		t.Errorf("Expected 'cost | cache', got '%s'", result)
	}
}

func TestSplitLayout_Empty(t *testing.T) {
	left := []string{}
	right := []string{}
	result := SplitLayout(left, right, testSeparator)

	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestSplitLayout_FiltersEmptyStrings(t *testing.T) {
	left := []string{"model", "", "git"}
	right := []string{"", "cost"}
	result := SplitLayout(left, right, testSeparator)

	if !strings.Contains(result, "model | git") {
		t.Errorf("Expected 'model | git' without empty, got '%s'", result)
	}
	if !strings.Contains(result, "cost") {
		t.Errorf("Expected 'cost', got '%s'", result)
	}
}

func TestJoinLines(t *testing.T) {
	lines := []string{"line1", "line2", "line3"}
	result := JoinLines(lines)

	expected := "line1\nline2\nline3"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestJoinLines_Single(t *testing.T) {
	lines := []string{"single"}
	result := JoinLines(lines)

	if result != "single" {
		t.Errorf("Expected 'single', got '%s'", result)
	}
}

func TestJoinLines_Empty(t *testing.T) {
	lines := []string{}
	result := JoinLines(lines)

	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}
