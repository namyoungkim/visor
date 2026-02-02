package render

import "testing"

func TestTruncate_Short(t *testing.T) {
	result := Truncate("hello", 10)
	if result != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result)
	}
}

func TestTruncate_Exact(t *testing.T) {
	result := Truncate("hello", 5)
	if result != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result)
	}
}

func TestTruncate_Long(t *testing.T) {
	result := Truncate("hello world", 8)
	if len(result) > 8 {
		t.Errorf("Result too long: '%s' (len=%d)", result, len(result))
	}
}

func TestTruncate_WithANSI(t *testing.T) {
	input := "\033[31mhello\033[0m"
	result := Truncate(input, 5)

	// ANSI codes should not count toward width
	visible := VisibleLength(result)
	if visible > 5 {
		t.Errorf("Visible length %d exceeds max 5", visible)
	}
}

func TestVisibleLength_Plain(t *testing.T) {
	length := VisibleLength("hello")
	if length != 5 {
		t.Errorf("Expected 5, got %d", length)
	}
}

func TestVisibleLength_WithANSI(t *testing.T) {
	length := VisibleLength("\033[31mhello\033[0m")
	if length != 5 {
		t.Errorf("Expected 5, got %d", length)
	}
}

func TestVisibleLength_CJK(t *testing.T) {
	length := VisibleLength("한글")
	if length != 4 {
		t.Errorf("Expected 4 (2 chars * 2 width), got %d", length)
	}
}
