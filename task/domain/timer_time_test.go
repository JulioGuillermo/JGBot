package taskdomain

import (
	"testing"
	"time"
)

func TestTimerTimeZeroValues(t *testing.T) {
	t.Parallel()

	timerTime := TimerTime{}

	if timerTime.Day != 0 {
		t.Errorf("TimerTime.Day = %d, want 0", timerTime.Day)
	}
	if timerTime.Month != 0 {
		t.Errorf("TimerTime.Month = %d, want 0", timerTime.Month)
	}
	if timerTime.Year != 0 {
		t.Errorf("TimerTime.Year = %d, want 0", timerTime.Year)
	}
	if timerTime.Hour != 0 {
		t.Errorf("TimerTime.Hour = %d, want 0", timerTime.Hour)
	}
	if timerTime.Minute != 0 {
		t.Errorf("TimerTime.Minute = %d, want 0", timerTime.Minute)
	}
	if timerTime.Second != 0 {
		t.Errorf("TimerTime.Second = %d, want 0", timerTime.Second)
	}
}

func TestTimerTimeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		timerTime TimerTime
		expected  string
	}{
		{
			name: "time only (zero date values)",
			timerTime: TimerTime{
				Day:    0,
				Month:  0,
				Year:   0,
				Hour:   10,
				Minute: 30,
				Second: 45,
			},
			expected: "10:30:45",
		},
		{
			name: "full date and time",
			timerTime: TimerTime{
				Day:    15,
				Month:  6,
				Year:   2024,
				Hour:   14,
				Minute: 30,
				Second: 0,
			},
			expected: "14:30:00 06/15/2024",
		},
		{
			name: "midnight",
			timerTime: TimerTime{
				Day:    1,
				Month:  1,
				Year:   2024,
				Hour:   0,
				Minute: 0,
				Second: 0,
			},
			expected: "00:00:00 01/01/2024",
		},
		{
			name: "end of day",
			timerTime: TimerTime{
				Day:    31,
				Month:  12,
				Year:   2024,
				Hour:   23,
				Minute: 59,
				Second: 59,
			},
			expected: "23:59:59 12/31/2024",
		},
		{
			name: "single digit values",
			timerTime: TimerTime{
				Day:    1,
				Month:  1,
				Year:   2024,
				Hour:   1,
				Minute: 1,
				Second: 1,
			},
			expected: "01:01:01 01/01/2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.timerTime.String()
			if result != tt.expected {
				t.Errorf("TimerTime.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTimerTimeToDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		timerTime TimerTime
		expected  time.Duration
	}{
		{
			name: "one second",
			timerTime: TimerTime{
				Second: 1,
			},
			expected: 1 * time.Second,
		},
		{
			name: "one minute",
			timerTime: TimerTime{
				Minute: 1,
			},
			expected: 1 * time.Minute,
		},
		{
			name: "one hour",
			timerTime: TimerTime{
				Hour: 1,
			},
			expected: 1 * time.Hour,
		},
		{
			name: "one day",
			timerTime: TimerTime{
				Day: 1,
			},
			expected: 24 * time.Hour,
		},
		{
			name: "one month (30 days)",
			timerTime: TimerTime{
				Month: 1,
			},
			expected: 30 * 24 * time.Hour,
		},
		{
			name: "one year (12 months)",
			timerTime: TimerTime{
				Year: 1,
			},
			expected: 12 * 30 * 24 * time.Hour,
		},
		{
			name: "complex duration",
			timerTime: TimerTime{
				Year:   1,
				Month:  2,
				Day:    3,
				Hour:   4,
				Minute: 5,
				Second: 6,
			},
			expected: (1*Year + 2*Month + 3*Day + 4*Hour + 5*Minute + 6) * time.Second,
		},
		{
			name:      "zero duration",
			timerTime: TimerTime{},
			expected:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.timerTime.ToDuration()
			if result != tt.expected {
				t.Errorf("TimerTime.ToDuration() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTimerTimeToTime(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name      string
		timerTime TimerTime
		validate  func(time.Time) bool
	}{
		{
			name: "full date and time",
			timerTime: TimerTime{
				Day:    15,
				Month:  6,
				Year:   2024,
				Hour:   14,
				Minute: 30,
				Second: 45,
			},
			validate: func(result time.Time) bool {
				return result.Year() == 2024 &&
					int(result.Month()) == 6 &&
					result.Day() == 15 &&
					result.Hour() == 14 &&
					result.Minute() == 30 &&
					result.Second() == 45
			},
		},
		{
			name: "time only (should use current date)",
			timerTime: TimerTime{
				Day:    0,
				Month:  0,
				Year:   0,
				Hour:   10,
				Minute: 30,
				Second: 0,
			},
			validate: func(result time.Time) bool {
				return result.Year() == now.Year() &&
					int(result.Month()) == int(now.Month()) &&
					result.Day() == now.Day() &&
					result.Hour() == 10 &&
					result.Minute() == 30 &&
					result.Second() == 0
			},
		},
		{
			name: "partial date (only month and year)",
			timerTime: TimerTime{
				Day:    0,
				Month:  12,
				Year:   2025,
				Hour:   0,
				Minute: 0,
				Second: 0,
			},
			validate: func(result time.Time) bool {
				return result.Year() == 2025 &&
					int(result.Month()) == 12 &&
					result.Hour() == 0 &&
					result.Minute() == 0 &&
					result.Second() == 0
			},
		},
		{
			name: "partial date (only year)",
			timerTime: TimerTime{
				Day:    0,
				Month:  0,
				Year:   2025,
				Hour:   12,
				Minute: 0,
				Second: 0,
			},
			validate: func(result time.Time) bool {
				return result.Year() == 2025 &&
					result.Hour() == 12 &&
					result.Minute() == 0 &&
					result.Second() == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.timerTime.ToTime()
			if !tt.validate(result) {
				t.Errorf("TimerTime.ToTime() = %v, does not match expected criteria", result)
			}
		})
	}
}

func TestTimerTimeToDurationCalculation(t *testing.T) {
	t.Parallel()

	// Test that the duration calculation uses the constants correctly
	timerTime := TimerTime{
		Year:   1,
		Month:  1,
		Day:    1,
		Hour:   1,
		Minute: 1,
		Second: 1,
	}

	expected := time.Duration(
		1*Year+
			1*Month+
			1*Day+
			1*Hour+
			1*Minute+
			1,
	) * time.Second

	result := timerTime.ToDuration()
	if result != expected {
		t.Errorf("TimerTime.ToDuration() = %v, want %v", result, expected)
	}
}

func TestTimerTimeEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("maximum values", func(t *testing.T) {
		timerTime := TimerTime{
			Day:    31,
			Month:  12,
			Year:   9999,
			Hour:   23,
			Minute: 59,
			Second: 59,
		}

		// Should not panic
		_ = timerTime.String()
		_ = timerTime.ToDuration()
		_ = timerTime.ToTime()
	})

	t.Run("negative values", func(t *testing.T) {
		timerTime := TimerTime{
			Day:    -1,
			Month:  -1,
			Year:   -1,
			Hour:   -1,
			Minute: -1,
			Second: -1,
		}

		// Should not panic, but may produce unexpected results
		_ = timerTime.String()
		_ = timerTime.ToDuration()
		_ = timerTime.ToTime()
	})
}

func TestTimerTimeConstants(t *testing.T) {
	t.Parallel()

	// Verify the constants are defined correctly
	if Minute != 60 {
		t.Errorf("Minute = %d, want 60", Minute)
	}
	if Hour != 60*Minute {
		t.Errorf("Hour = %d, want %d", Hour, 60*Minute)
	}
	if Day != 24*Hour {
		t.Errorf("Day = %d, want %d", Day, 24*Hour)
	}
	if Month != 30*Day {
		t.Errorf("Month = %d, want %d", Month, 30*Day)
	}
	if Year != 12*Month {
		t.Errorf("Year = %d, want %d", Year, 12*Month)
	}
}
