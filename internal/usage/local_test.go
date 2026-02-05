package usage

import (
	"testing"
	"time"

	"github.com/namyoungkim/visor/internal/cost"
)

func TestDefaultLimitForTier(t *testing.T) {
	tests := []struct {
		tier      string
		wantFive  int
		wantSeven int
	}{
		{"", 45, 675},
		{"pro", 45, 675},
		{"default_claude_max_5x", 225, 3375},
		{"DEFAULT_CLAUDE_MAX_5X", 225, 3375},
		{"default_claude_max_20x", 900, 13500},
		{"some_unknown_tier", 45, 675},
	}

	for _, tt := range tests {
		t.Run(tt.tier, func(t *testing.T) {
			gotFive, gotSeven := DefaultLimitForTier(tt.tier)
			if gotFive != tt.wantFive {
				t.Errorf("DefaultLimitForTier(%q) fiveHour = %d, want %d", tt.tier, gotFive, tt.wantFive)
			}
			if gotSeven != tt.wantSeven {
				t.Errorf("DefaultLimitForTier(%q) sevenDay = %d, want %d", tt.tier, gotSeven, tt.wantSeven)
			}
		})
	}
}

func TestEstimateLimits_Basic(t *testing.T) {
	now := time.Now()
	blockStart := now.Add(-1 * time.Hour)

	costData := &cost.CostData{
		FiveHourBlockMessages: 10,
		WeekMessages:          50,
	}

	limits := EstimateLimits(costData, blockStart, "", 45, 675)
	if limits == nil {
		t.Fatal("expected non-nil limits")
	}

	// 10/45 = ~22.2%
	expectedFiveHour := float64(10) / float64(45) * 100
	if limits.FiveHour.Utilization != expectedFiveHour {
		t.Errorf("FiveHour.Utilization = %f, want %f", limits.FiveHour.Utilization, expectedFiveHour)
	}

	// 50/675 = ~7.4%
	expectedSevenDay := float64(50) / float64(675) * 100
	if limits.SevenDay.Utilization != expectedSevenDay {
		t.Errorf("SevenDay.Utilization = %f, want %f", limits.SevenDay.Utilization, expectedSevenDay)
	}

	// Remaining
	if limits.FiveHour.Remaining != 35 {
		t.Errorf("FiveHour.Remaining = %d, want 35", limits.FiveHour.Remaining)
	}
	if limits.SevenDay.Remaining != 625 {
		t.Errorf("SevenDay.Remaining = %d, want 625", limits.SevenDay.Remaining)
	}

	// Total
	if limits.FiveHour.Total != 45 {
		t.Errorf("FiveHour.Total = %d, want 45", limits.FiveHour.Total)
	}
	if limits.SevenDay.Total != 675 {
		t.Errorf("SevenDay.Total = %d, want 675", limits.SevenDay.Total)
	}

	// ResetsAt for 5-hour block
	expectedReset := blockStart.Add(cost.BlockDuration)
	if !limits.FiveHour.ResetsAt.Equal(expectedReset) {
		t.Errorf("FiveHour.ResetsAt = %v, want %v", limits.FiveHour.ResetsAt, expectedReset)
	}
}

func TestEstimateLimits_CustomLimits(t *testing.T) {
	now := time.Now()
	blockStart := now.Add(-30 * time.Minute)

	costData := &cost.CostData{
		FiveHourBlockMessages: 100,
		WeekMessages:          500,
	}

	limits := EstimateLimits(costData, blockStart, "", 225, 3375)
	if limits == nil {
		t.Fatal("expected non-nil limits")
	}

	// 100/225 = ~44.4%
	expectedFiveHour := float64(100) / float64(225) * 100
	if limits.FiveHour.Utilization != expectedFiveHour {
		t.Errorf("FiveHour.Utilization = %f, want %f", limits.FiveHour.Utilization, expectedFiveHour)
	}

	// 500/3375 = ~14.8%
	expectedSevenDay := float64(500) / float64(3375) * 100
	if limits.SevenDay.Utilization != expectedSevenDay {
		t.Errorf("SevenDay.Utilization = %f, want %f", limits.SevenDay.Utilization, expectedSevenDay)
	}
}

func TestEstimateLimits_ZeroMessages(t *testing.T) {
	now := time.Now()
	blockStart := now.Add(-10 * time.Minute)

	costData := &cost.CostData{
		FiveHourBlockMessages: 0,
		WeekMessages:          0,
	}

	limits := EstimateLimits(costData, blockStart, "", 45, 675)
	if limits == nil {
		t.Fatal("expected non-nil limits")
	}

	if limits.FiveHour.Utilization != 0 {
		t.Errorf("FiveHour.Utilization = %f, want 0", limits.FiveHour.Utilization)
	}
	if limits.SevenDay.Utilization != 0 {
		t.Errorf("SevenDay.Utilization = %f, want 0", limits.SevenDay.Utilization)
	}
	if limits.FiveHour.Remaining != 45 {
		t.Errorf("FiveHour.Remaining = %d, want 45", limits.FiveHour.Remaining)
	}
	if limits.SevenDay.Remaining != 675 {
		t.Errorf("SevenDay.Remaining = %d, want 675", limits.SevenDay.Remaining)
	}
}

func TestEstimateLimits_OverLimit(t *testing.T) {
	now := time.Now()
	blockStart := now.Add(-4 * time.Hour)

	costData := &cost.CostData{
		FiveHourBlockMessages: 60,
		WeekMessages:          800,
	}

	limits := EstimateLimits(costData, blockStart, "", 45, 675)
	if limits == nil {
		t.Fatal("expected non-nil limits")
	}

	// Should cap at 100%
	if limits.FiveHour.Utilization != 100 {
		t.Errorf("FiveHour.Utilization = %f, want 100", limits.FiveHour.Utilization)
	}
	if limits.SevenDay.Utilization != 100 {
		t.Errorf("SevenDay.Utilization = %f, want 100", limits.SevenDay.Utilization)
	}

	// Remaining should be 0
	if limits.FiveHour.Remaining != 0 {
		t.Errorf("FiveHour.Remaining = %d, want 0", limits.FiveHour.Remaining)
	}
	if limits.SevenDay.Remaining != 0 {
		t.Errorf("SevenDay.Remaining = %d, want 0", limits.SevenDay.Remaining)
	}
}

func TestEstimateLimits_NilCostData(t *testing.T) {
	limits := EstimateLimits(nil, time.Now(), "", 45, 675)
	if limits != nil {
		t.Error("expected nil limits for nil costData")
	}
}

func TestEstimateLimits_TierAutoDetect(t *testing.T) {
	now := time.Now()
	blockStart := now.Add(-1 * time.Hour)

	costData := &cost.CostData{
		FiveHourBlockMessages: 100,
		WeekMessages:          500,
	}

	// Pass tier but zero limits → should auto-detect Max 5x limits (225, 3375)
	limits := EstimateLimits(costData, blockStart, "default_claude_max_5x", 0, 0)
	if limits == nil {
		t.Fatal("expected non-nil limits")
	}

	if limits.FiveHour.Total != 225 {
		t.Errorf("FiveHour.Total = %d, want 225 (Max 5x)", limits.FiveHour.Total)
	}
	if limits.SevenDay.Total != 3375 {
		t.Errorf("SevenDay.Total = %d, want 3375 (Max 5x)", limits.SevenDay.Total)
	}

	// 100/225 = ~44.4%
	expectedUtil := float64(100) / float64(225) * 100
	if limits.FiveHour.Utilization != expectedUtil {
		t.Errorf("FiveHour.Utilization = %f, want %f", limits.FiveHour.Utilization, expectedUtil)
	}
}

func TestEstimateLimits_ZeroBlockStart(t *testing.T) {
	costData := &cost.CostData{
		FiveHourBlockMessages: 5,
		WeekMessages:          20,
	}

	limits := EstimateLimits(costData, time.Time{}, "", 45, 675)
	if limits == nil {
		t.Fatal("expected non-nil limits")
	}

	// ResetsAt should be zero for 5-hour block when blockStart is zero
	if !limits.FiveHour.ResetsAt.IsZero() {
		t.Errorf("FiveHour.ResetsAt should be zero, got %v", limits.FiveHour.ResetsAt)
	}
}
