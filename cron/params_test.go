package cron

import "testing"

func TestGetNumberParam(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1", "1"},
		{"1,2,3", "1,2,3"},
		{"1,2,3,4,5,6, 7,8,9,10,11,12", "1,2,3,4,5,6,7,8,9,10,11,12"},
		{"1,2,3,4,5,6, 7,8,9,10,11,12,1,2,3,4,5,6,7,8,9,10,11,12", "1,2,3,4,5,6,7,8,9,10,11,12,1,2,3,4,5,6,7,8,9,10,11,12"},
		{"[1,2, 3]", "1,2,3"},
		{"[1,2,3,4,5,6,7,8,9,10,11,12]", "1,2,3,4,5,6,7,8,9,10,11,12"},
		{"[1,2,3,4,5,6,7,8,9,10,11,12,1,2,3,4,5,6,7,8,9,10,11,12]", "1,2,3,4,5,6,7,8,9,10,11,12,1,2,3,4,5,6,7,8,9,10,11,12"},
		{"every 1", "*/1"},
		{"every (1)", "*/1"},
		{"every(1)", "*/1"},
	}

	for _, test := range tests {
		if result := GetNumberParam(test.input); result != test.expected {
			t.Errorf("GetNumberParam(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}
