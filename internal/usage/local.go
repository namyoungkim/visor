package usage

import (
	"strings"
	"time"

	"github.com/namyoungkim/visor/internal/cost"
)

// DefaultLimitForTier returns default message limits based on subscription tier.
// Returns (fiveHourLimit, sevenDayLimit).
//
// These values are estimates based on observed Claude Pro/Max behavior (as of 2026-02).
// Pro: ~45 messages/5h, Max 5x: ~225/5h, Max 20x: ~900/5h.
// 7-day = 5-hour limit * 15 (conservative: ~3 blocks/day * 5 days).
func DefaultLimitForTier(tier string) (int, int) {
	t := strings.ToLower(tier)
	switch {
	case strings.Contains(t, "20x"):
		return 900, 13500
	case strings.Contains(t, "5x"):
		return 225, 3375
	default:
		// Pro default
		return 45, 675
	}
}

// EstimateLimits creates Limits from local JSONL message counts.
// tier is the subscription rate limit tier (e.g. "default_claude_max_5x").
// If fiveHourLimit or sevenDayLimit is 0, they are auto-detected from the tier.
// If tier is also empty, Pro defaults are used.
func EstimateLimits(costData *cost.CostData, blockStart time.Time, tier string, fiveHourLimit, sevenDayLimit int) *Limits {
	if costData == nil {
		return nil
	}

	// Auto-detect limits from tier if not specified
	if fiveHourLimit == 0 || sevenDayLimit == 0 {
		autoFive, autoSeven := DefaultLimitForTier(tier)
		if fiveHourLimit == 0 {
			fiveHourLimit = autoFive
		}
		if sevenDayLimit == 0 {
			sevenDayLimit = autoSeven
		}
	}

	now := time.Now()

	// 5-hour block utilization
	fiveHourUtil := float64(0)
	if fiveHourLimit > 0 {
		fiveHourUtil = float64(costData.FiveHourBlockMessages) / float64(fiveHourLimit) * 100
		if fiveHourUtil > 100 {
			fiveHourUtil = 100
		}
	}

	// 7-day utilization
	sevenDayUtil := float64(0)
	if sevenDayLimit > 0 {
		sevenDayUtil = float64(costData.WeekMessages) / float64(sevenDayLimit) * 100
		if sevenDayUtil > 100 {
			sevenDayUtil = 100
		}
	}

	// Calculate reset times
	var fiveHourReset time.Time
	if !blockStart.IsZero() {
		fiveHourReset = blockStart.Add(cost.BlockDuration)
	}

	weekStart := cost.StartOfWeek(now)
	sevenDayReset := weekStart.Add(7 * 24 * time.Hour)

	fiveHourRemaining := fiveHourLimit - costData.FiveHourBlockMessages
	if fiveHourRemaining < 0 {
		fiveHourRemaining = 0
	}
	sevenDayRemaining := sevenDayLimit - costData.WeekMessages
	if sevenDayRemaining < 0 {
		sevenDayRemaining = 0
	}

	return &Limits{
		FiveHour: FiveHourLimit{
			Utilization: fiveHourUtil,
			ResetsAt:    fiveHourReset,
			Remaining:   fiveHourRemaining,
			Total:       fiveHourLimit,
		},
		SevenDay: SevenDayLimit{
			Utilization: sevenDayUtil,
			ResetsAt:    sevenDayReset,
			Remaining:   sevenDayRemaining,
			Total:       sevenDayLimit,
		},
	}
}
