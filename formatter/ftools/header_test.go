package ftools

import "testing"

func TestFormatHeader(t *testing.T) {
	tests := []struct {
		input string
		want  string
		ok    bool
	}{
		{"# Title", "1️⃣ Title", true},
		{"## Title", "2️⃣ Title", true},
		{"### Title", "3️⃣ Title", true},
		{"#### Title", "4️⃣ Title", true},
		{"##### Title", "5️⃣ Title", true},
		{"###### Title", "6️⃣ Title", true},
		{"Not Header", "Not Header", false},
	}

	for _, tt := range tests {
		got, ok := FormatHeader(tt.input)
		if ok != tt.ok || got != tt.want {
			t.Errorf("FormatHeader(%q) = %q, %v; want %q, %v", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}
