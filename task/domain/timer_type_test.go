package taskdomain

import (
	"testing"
)

func TestTimerTypeConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		timerType TimerType
		expected  string
	}{
		{
			name:      "ALARM type",
			timerType: ALARM,
			expected:  "ALARM",
		},
		{
			name:      "TIMEOUT type",
			timerType: TIMEOUT,
			expected:  "TIMEOUT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if string(tt.timerType) != tt.expected {
				t.Errorf("TimerType = %v, want %v", tt.timerType, tt.expected)
			}
		})
	}
}

func TestTimerTypeInequality(t *testing.T) {
	t.Parallel()

	if ALARM == TIMEOUT {
		t.Error("ALARM should not equal TIMEOUT")
	}
}

func TestTimerTypeFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected TimerType
	}{
		{
			name:     "ALARM from string",
			input:    "ALARM",
			expected: ALARM,
		},
		{
			name:     "TIMEOUT from string",
			input:    "TIMEOUT",
			expected: TIMEOUT,
		},
		{
			name:     "unknown type",
			input:    "UNKNOWN",
			expected: "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := TimerType(tt.input)
			if result != tt.expected {
				t.Errorf("TimerType(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
