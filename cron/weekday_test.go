package cron

import "testing"

func TestReplaceWeekday(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"sun", "0"},
		{"mon", "1"},
		{"tue", "2"},
		{"wed", "3"},
		{"thu", "4"},
		{"fri", "5"},
		{"sat", "6"},
		{"Sun", "0"},
		{"Mon", "1"},
		{"Tue", "2"},
		{"Wed", "3"},
		{"Thu", "4"},
		{"Fri", "5"},
		{"Sat", "6"},
		{"Sunday", "0"},
		{"Monday", "1"},
		{"Tuesday", "2"},
		{"Wednesday", "3"},
		{"Thursday", "4"},
		{"Friday", "5"},
		{"Saturday", "6"},
		{"Sun", "0"},
		{"Mon", "1"},
		{"Tue", "2"},
		{"Wed", "3"},
		{"Thu", "4"},
		{"Fri", "5"},
		{"Sat", "6"},
	}

	for _, test := range tests {
		if result := ReplaceWeekday(test.input); result != test.expected {
			t.Errorf("ReplaceWeekday(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}
