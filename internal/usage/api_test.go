package usage

import (
	"testing"
	"time"
)

func TestLimitsFiveHourRemaining(t *testing.T) {
	limits := &Limits{
		FiveHour: FiveHourLimit{
			Utilization: 50.0,
			ResetsAt:    time.Now().Add(2*time.Hour + 30*time.Minute),
		},
	}

	remaining := limits.FiveHourRemaining()

	// Should be approximately 2.5 hours
	expected := 2*time.Hour + 30*time.Minute
	tolerance := 2 * time.Second

	if remaining < expected-tolerance || remaining > expected+tolerance {
		t.Errorf("FiveHourRemaining() = %v, want ~%v", remaining, expected)
	}
}

func TestLimitsSevenDayRemaining(t *testing.T) {
	limits := &Limits{
		SevenDay: SevenDayLimit{
			Utilization: 30.0,
			ResetsAt:    time.Now().Add(5 * 24 * time.Hour),
		},
	}

	remaining := limits.SevenDayRemaining()

	// Should be approximately 5 days
	expected := 5 * 24 * time.Hour
	tolerance := 2 * time.Second

	if remaining < expected-tolerance || remaining > expected+tolerance {
		t.Errorf("SevenDayRemaining() = %v, want ~%v", remaining, expected)
	}
}

func TestLimitsExpired(t *testing.T) {
	limits := &Limits{
		FiveHour: FiveHourLimit{
			Utilization: 100.0,
			ResetsAt:    time.Now().Add(-1 * time.Hour), // Expired
		},
	}

	remaining := limits.FiveHourRemaining()

	if remaining != 0 {
		t.Errorf("FiveHourRemaining() for expired = %v, want 0", remaining)
	}
}

func TestLimitsZeroTime(t *testing.T) {
	limits := &Limits{
		FiveHour: FiveHourLimit{
			Utilization: 0,
			// ResetsAt is zero
		},
	}

	remaining := limits.FiveHourRemaining()

	if remaining != 0 {
		t.Errorf("FiveHourRemaining() for zero time = %v, want 0", remaining)
	}
}
