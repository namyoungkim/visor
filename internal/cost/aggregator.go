package cost

import (
	"sort"
	"time"
)

// CostData holds aggregated cost information.
type CostData struct {
	Today         float64 // Cost in current calendar day
	Week          float64 // Cost in current week (Monday-Sunday)
	Month         float64 // Cost in current month
	FiveHourBlock float64 // Cost in current 5-hour block

	BlockStartTime time.Time // Start time of current 5-hour block
	Provider       Provider
	LastUpdated    time.Time

	// Per-session cost (for current session)
	SessionCost float64

	// Message counts for local usage estimation
	TodayMessages         int // Messages in current calendar day
	WeekMessages          int // Messages in current week (Monday-Sunday)
	FiveHourBlockMessages int // Messages in current 5-hour block
}

// Aggregate computes aggregated costs from entries.
func Aggregate(entries []Entry, blockStart time.Time) *CostData {
	now := time.Now()
	data := &CostData{
		Provider:       DetectProvider(),
		LastUpdated:    now,
		BlockStartTime: blockStart,
	}

	if len(entries) == 0 {
		return data
	}

	// Sort entries by timestamp
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})

	// Get time boundaries
	todayStart := startOfDay(now)
	weekStart := startOfWeek(now)
	monthStart := startOfMonth(now)
	blockEnd := blockStart.Add(5 * time.Hour)

	for _, e := range entries {
		// Today's cost
		if !e.Timestamp.Before(todayStart) {
			data.Today += e.CostUSD
			data.TodayMessages++
		}

		// Week's cost
		if !e.Timestamp.Before(weekStart) {
			data.Week += e.CostUSD
			data.WeekMessages++
		}

		// Month's cost
		if !e.Timestamp.Before(monthStart) {
			data.Month += e.CostUSD
		}

		// 5-hour block cost
		if !blockStart.IsZero() && !e.Timestamp.Before(blockStart) && e.Timestamp.Before(blockEnd) {
			data.FiveHourBlock += e.CostUSD
			data.FiveHourBlockMessages++
		}
	}

	return data
}

// AggregateSession computes cost for a specific session.
func AggregateSession(entries []Entry, sessionID string) float64 {
	var total float64
	for _, e := range entries {
		if e.SessionID == sessionID {
			total += e.CostUSD
		}
	}
	return total
}

// startOfDay returns the start of the current day in local time.
func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// startOfWeek returns the start of the current week (Monday) in local time.
func startOfWeek(t time.Time) time.Time {
	day := startOfDay(t)
	weekday := int(day.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	return day.AddDate(0, 0, -(weekday - 1))
}

// startOfMonth returns the start of the current month in local time.
func startOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

// RemainingIn5HourBlock returns the remaining time in the 5-hour block.
func RemainingIn5HourBlock(blockStart time.Time) time.Duration {
	if blockStart.IsZero() {
		return 0
	}

	blockEnd := blockStart.Add(5 * time.Hour)
	remaining := time.Until(blockEnd)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// PercentElapsedIn5HourBlock returns the percentage elapsed in the 5-hour block.
func PercentElapsedIn5HourBlock(blockStart time.Time) float64 {
	if blockStart.IsZero() {
		return 0
	}

	elapsed := time.Since(blockStart)
	total := 5 * time.Hour

	pct := float64(elapsed) / float64(total) * 100
	if pct < 0 {
		return 0
	}
	if pct > 100 {
		return 100
	}
	return pct
}
