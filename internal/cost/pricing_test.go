package cost

import (
	"math"
	"testing"
)

func TestGetPricing(t *testing.T) {
	tests := []struct {
		name      string
		modelID   string
		wantInput float64
	}{
		{
			name:      "opus 4.5",
			modelID:   "claude-opus-4-5-20251101",
			wantInput: 15.0,
		},
		{
			name:      "sonnet 4",
			modelID:   "claude-sonnet-4-20250514",
			wantInput: 3.0,
		},
		{
			name:      "3.5 sonnet",
			modelID:   "claude-3-5-sonnet-20241022",
			wantInput: 3.0,
		},
		{
			name:      "3.5 haiku",
			modelID:   "claude-3-5-haiku-20241022",
			wantInput: 1.0,
		},
		{
			name:      "unknown model uses default",
			modelID:   "claude-unknown-99",
			wantInput: 3.0, // default pricing
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := GetPricing(tt.modelID)
			if p.InputPer1M != tt.wantInput {
				t.Errorf("GetPricing(%q).InputPer1M = %v, want %v", tt.modelID, p.InputPer1M, tt.wantInput)
			}
		})
	}
}

func TestCalculateCost(t *testing.T) {
	// Test with known values
	// 1000 input tokens at $3/1M = $0.003
	// 500 output tokens at $15/1M = $0.0075
	// Total = $0.0105
	cost := CalculateCost("claude-3-5-sonnet-20241022", 1000, 500, 0, 0)

	expected := 0.0105
	if math.Abs(cost-expected) > 0.0001 {
		t.Errorf("CalculateCost() = %v, want %v", cost, expected)
	}
}

func TestCalculateCostWithCache(t *testing.T) {
	// Test with cache tokens
	// 1000 input at $3/1M = $0.003
	// 500 output at $15/1M = $0.0075
	// 2000 cache read at $0.30/1M = $0.0006
	// Total = $0.0111
	cost := CalculateCost("claude-3-5-sonnet-20241022", 1000, 500, 2000, 0)

	expected := 0.0111
	if math.Abs(cost-expected) > 0.0001 {
		t.Errorf("CalculateCost() with cache = %v, want %v", cost, expected)
	}
}
