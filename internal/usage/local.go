package usage

import (
	"strings"
	"time"

	"github.com/namyoungkim/visor/internal/auth"
	"github.com/namyoungkim/visor/internal/cost"
)

// DefaultLimitForTier returns default message limits based on subscription tier.
// Returns (fiveHourLimit, sevenDayLimit).
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
// If fiveHourLimit or sevenDayLimit is 0, it attempts to auto-detect from the
// subscription tier via auth credentials.
func EstimateLimits(costData *cost.CostData, blockStart time.Time, fiveHourLimit, sevenDayLimit int) *Limits {
	if costData == nil {
		return nil
	}

	// Auto-detect limits from tier if not specified
	if fiveHourLimit == 0 || sevenDayLimit == 0 {
		tier := detectTier()
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
		fiveHourReset = blockStart.Add(5 * time.Hour)
	}

	weekStart := startOfWeek(now)
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

// detectTier attempts to read the rate limit tier from auth credentials.
func detectTier() string {
	provider := auth.DefaultProvider()
	creds, err := provider.Get()
	if err != nil || creds == nil {
		return ""
	}
	return creds.RateLimitTier
}

// startOfWeek returns the start of the current week (Monday) in local time.
func startOfWeek(t time.Time) time.Time {
	day := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	weekday := int(day.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	return day.AddDate(0, 0, -(weekday - 1))
}
