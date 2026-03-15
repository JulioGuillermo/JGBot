package taskdomain

import (
	"testing"
)

func TestFileConstants(t *testing.T) {
	t.Parallel()

	if TIMER_FILE == "" {
		t.Error("TIMER_FILE constant should not be empty")
	}

	if CRON_FILE == "" {
		t.Error("CRON_FILE constant should not be empty")
	}

	if TIMER_FILE == CRON_FILE {
		t.Error("TIMER_FILE and CRON_FILE should be different")
	}
}

func TestTimeConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		constant int
		expected int
	}{
		{
			name:     "Minute",
			constant: Minute,
			expected: 60,
		},
		{
			name:     "Hour",
			constant: Hour,
			expected: 3600,
		},
		{
			name:     "Day",
			constant: Day,
			expected: 86400,
		},
		{
			name:     "Month",
			constant: Month,
			expected: 2592000,
		},
		{
			name:     "Year",
			constant: Year,
			expected: 31104000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.constant != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

func TestTimeConstantRelationships(t *testing.T) {
	t.Parallel()

	// Verify relationships between constants
	if Hour != 60*Minute {
		t.Errorf("Hour (%d) should equal 60 * Minute (%d)", Hour, 60*Minute)
	}

	if Day != 24*Hour {
		t.Errorf("Day (%d) should equal 24 * Hour (%d)", Day, 24*Hour)
	}

	if Month != 30*Day {
		t.Errorf("Month (%d) should equal 30 * Day (%d)", Month, 30*Day)
	}

	if Year != 12*Month {
		t.Errorf("Year (%d) should equal 12 * Month (%d)", Year, 12*Month)
	}
}

func TestTimeConstantsOrdering(t *testing.T) {
	t.Parallel()

	// Verify that constants are in ascending order
	if !(Minute < Hour) {
		t.Error("Minute should be less than Hour")
	}
	if !(Hour < Day) {
		t.Error("Hour should be less than Day")
	}
	if !(Day < Month) {
		t.Error("Day should be less than Month")
	}
	if !(Month < Year) {
		t.Error("Month should be less than Year")
	}
}
