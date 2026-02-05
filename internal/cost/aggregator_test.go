package cost

import (
	"math"
	"testing"
	"time"
)

const floatTolerance = 0.0001

func floatEqual(a, b float64) bool {
	return math.Abs(a-b) < floatTolerance
}

func TestAggregate(t *testing.T) {
	now := time.Now()
	// Calculate today's start to ensure entries are clearly within/outside today
	todayStart := startOfDay(now)

	entries := []Entry{
		{
			// Clearly within today (after today's start + buffer)
			Timestamp: todayStart.Add(1 * time.Hour),
			CostUSD:   0.10,
		},
		{
			// Clearly within today
			Timestamp: todayStart.Add(2 * time.Hour),
			CostUSD:   0.20,
		},
		{
			// Clearly yesterday (before today's start)
			Timestamp: todayStart.Add(-10 * time.Hour),
			CostUSD:   0.50,
		},
	}

	data := Aggregate(entries, now.Add(-3*time.Hour))

	// Today should include first two entries (both are after todayStart)
	if !floatEqual(data.Today, 0.30) {
		t.Errorf("Aggregate().Today = %v, want 0.30", data.Today)
	}

	// Week should include all entries (all within same week)
	if !floatEqual(data.Week, 0.80) {
		t.Errorf("Aggregate().Week = %v, want 0.80", data.Week)
	}
}

func TestAggregateEmpty(t *testing.T) {
	data := Aggregate(nil, time.Time{})

	if data.Today != 0 {
		t.Errorf("Aggregate(nil).Today = %v, want 0", data.Today)
	}
	if data.Week != 0 {
		t.Errorf("Aggregate(nil).Week = %v, want 0", data.Week)
	}
}

func TestAggregate5HourBlock(t *testing.T) {
	now := time.Now()
	blockStart := now.Add(-2 * time.Hour)

	entries := []Entry{
		{
			Timestamp: now.Add(-1 * time.Hour), // In block
			CostUSD:   0.10,
		},
		{
			Timestamp: now.Add(-3 * time.Hour), // Before block
			CostUSD:   0.20,
		},
	}

	data := Aggregate(entries, blockStart)

	if !floatEqual(data.FiveHourBlock, 0.10) {
		t.Errorf("Aggregate().FiveHourBlock = %v, want 0.10", data.FiveHourBlock)
	}
}

func TestStartOfDay(t *testing.T) {
	// Test with a specific time
	input := time.Date(2024, 6, 15, 14, 30, 45, 0, time.Local)
	expected := time.Date(2024, 6, 15, 0, 0, 0, 0, time.Local)

	result := startOfDay(input)
	if !result.Equal(expected) {
		t.Errorf("startOfDay() = %v, want %v", result, expected)
	}
}

func TestStartOfWeek(t *testing.T) {
	tests := []struct {
		input    time.Time
		expected time.Time
	}{
		{
			input:    time.Date(2024, 6, 15, 14, 0, 0, 0, time.Local), // Saturday
			expected: time.Date(2024, 6, 10, 0, 0, 0, 0, time.Local),  // Monday
		},
		{
			input:    time.Date(2024, 6, 10, 10, 0, 0, 0, time.Local), // Monday
			expected: time.Date(2024, 6, 10, 0, 0, 0, 0, time.Local),  // Same Monday
		},
		{
			input:    time.Date(2024, 6, 16, 10, 0, 0, 0, time.Local), // Sunday
			expected: time.Date(2024, 6, 10, 0, 0, 0, 0, time.Local),  // Previous Monday
		},
	}

	for _, tt := range tests {
		result := StartOfWeek(tt.input)
		if !result.Equal(tt.expected) {
			t.Errorf("StartOfWeek(%v) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestPercentElapsedIn5HourBlock(t *testing.T) {
	// Block started 2.5 hours ago = 50% elapsed
	blockStart := time.Now().Add(-150 * time.Minute)
	pct := PercentElapsedIn5HourBlock(blockStart)

	if pct < 49 || pct > 51 {
		t.Errorf("PercentElapsedIn5HourBlock() = %v, want ~50", pct)
	}
}

func TestRemainingIn5HourBlock(t *testing.T) {
	// Block started 2 hours ago = 3 hours remaining
	blockStart := time.Now().Add(-2 * time.Hour)
	remaining := RemainingIn5HourBlock(blockStart)

	expected := 3 * time.Hour
	tolerance := 2 * time.Second

	if remaining < expected-tolerance || remaining > expected+tolerance {
		t.Errorf("RemainingIn5HourBlock() = %v, want ~%v", remaining, expected)
	}
}
