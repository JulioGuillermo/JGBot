package ftools

import "testing"

func TestFormatHeader(t *testing.T) {
	tests := []struct {
		input string
		want  string
		ok    bool
	}{
		{"# Title", "🔹 *Title*", true},
		{"## Title", "🔹 *_Title_*", true},
		{"### Title", "🔹 _Title_", true},
		{"Not Header", "Not Header", false},
	}

	for _, tt := range tests {
		got, ok := FormatHeader(tt.input)
		if ok != tt.ok || got != tt.want {
			t.Errorf("FormatHeader(%q) = %q, %v; want %q, %v", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}
