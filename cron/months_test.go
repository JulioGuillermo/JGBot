package cron

import "testing"

func TestReplaceMonth(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Simple tests
		{"jan", "1"},
		{"feb", "2"},
		{"mar", "3"},
		{"apr", "4"},
		{"may", "5"},
		{"jun", "6"},
		{"jul", "7"},
		{"aug", "8"},
		{"sep", "9"},
		{"oct", "10"},
		{"nov", "11"},
		{"dec", "12"},
		{"Jan", "1"},
		{"Feb", "2"},
		{"Mar", "3"},
		{"Apr", "4"},
		{"May", "5"},
		{"Jun", "6"},
		{"Jul", "7"},
		{"Aug", "8"},
		{"Sep", "9"},
		{"Oct", "10"},
		{"Nov", "11"},
		{"Dec", "12"},
		{"January", "1"},
		{"February", "2"},
		{"March", "3"},
		{"April", "4"},
		{"May", "5"},
		{"June", "6"},
		{"July", "7"},
		{"August", "8"},
		{"September", "9"},
		{"October", "10"},
		{"November", "11"},
		{"December", "12"},
		{"Jan", "1"},
		{"Feb", "2"},
		{"Mar", "3"},
		{"Apr", "4"},
		{"May", "5"},
		{"Jun", "6"},
		{"Jul", "7"},
		{"Aug", "8"},
		{"Sep", "9"},
		{"Oct", "10"},
		{"Nov", "11"},
		{"Dec", "12"},

		// List tests
		{"jan,feb,mar", "1,2,3"},
		{"jan,mar", "1,3"},
		{"jan,feb,mar,apr,may,jun, jul,aug,sep,oct,nov,dec", "1,2,3,4,5,6, 7,8,9,10,11,12"},
		{"jan,dec", "1,12"},
		{"jan,feb,mar,apr,may,jun, jul,aug,sep,oct,nov,dec,jan,feb,mar,apr,may,jun,jul,aug,sep,oct,nov,dec", "1,2,3,4,5,6, 7,8,9,10,11,12,1,2,3,4,5,6,7,8,9,10,11,12"},
		{"[jan,feb,mar]", "[1,2,3]"},
		{"[jan,mar]", "[1,3]"},
		{"[jan,feb,mar,apr,may, jun, jul,aug,sep,oct,nov,dec]", "[1,2,3,4,5, 6, 7,8,9,10,11,12]"},
		{"[jan,dec]", "[1,12]"},
		{"[jan,feb,mar,apr,may,jun,jul,aug,sep,oct,nov,dec,jan,feb,mar,apr,may,jun,jul,aug,sep,oct,nov,dec]", "[1,2,3,4,5,6,7,8,9,10,11,12,1,2,3,4,5,6,7,8,9,10,11,12]"},

		// Every tests
		{"every jan", "every 1"},
		{"every (jan)", "every (1)"},
	}

	for _, test := range tests {
		if result := ReplaceMonth(test.input); result != test.expected {
			t.Errorf("ReplaceMonth(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}
