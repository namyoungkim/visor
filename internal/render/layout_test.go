package render

import (
	"strings"
	"testing"
)

func TestLayout_SingleWidget(t *testing.T) {
	widgets := []string{"model"}
	result := Layout(widgets)

	if result != "model" {
		t.Errorf("Expected 'model', got '%s'", result)
	}
}

func TestLayout_MultipleWidgets(t *testing.T) {
	widgets := []string{"model", "context", "cost"}
	result := Layout(widgets)

	if !strings.Contains(result, "model") {
		t.Error("Expected result to contain 'model'")
	}
	if !strings.Contains(result, "context") {
		t.Error("Expected result to contain 'context'")
	}
	if !strings.Contains(result, "cost") {
		t.Error("Expected result to contain 'cost'")
	}
}

func TestLayout_EmptyWidgets(t *testing.T) {
	widgets := []string{}
	result := Layout(widgets)

	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestLayout_WithEmptyStrings(t *testing.T) {
	widgets := []string{"model", "", "cost", ""}
	result := Layout(widgets)

	// Empty strings should be filtered out
	if strings.Contains(result, "  ") {
		t.Error("Expected empty widgets to be filtered out")
	}
}

func TestMultiLine_SingleLine(t *testing.T) {
	lines := [][]string{
		{"model", "context"},
	}
	result := MultiLine(lines)

	if strings.Contains(result, "\n") {
		t.Error("Expected single line output")
	}
}

func TestMultiLine_MultipleLines(t *testing.T) {
	lines := [][]string{
		{"model", "context"},
		{"cost", "git"},
	}
	result := MultiLine(lines)

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
	result := MultiLine(lines)

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
	result := MultiLine(lines)

	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}
