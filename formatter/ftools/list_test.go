package ftools

import "testing"

func TestFormatList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Todo list with dash",
			input:    "- [ ] To do",
			expected: "- ⬜ To do",
		},
		{
			name:     "Todo list with star",
			input:    "* [ ] To do",
			expected: "- ⬜ To do",
		},
		{
			name:     "Done list lowercase",
			input:    "- [x] Done",
			expected: "- ✅ Done",
		},
		{
			name:     "Done list uppercase",
			input:    "* [X] Done",
			expected: "- ✅ Done",
		},
		{
			name:     "Bullet list",
			input:    "* Item 1\n* Item 2",
			expected: "- Item 1\n- Item 2",
		},
		{
			name:     "Numbered list",
			input:    "1. First\n2. Second",
			expected: "1. First\n2. Second",
		},
		{
			name:     "Indented lines stay untouched",
			input:    "  - nested",
			expected: "  - nested",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatList(tt.input)
			if got != tt.expected {
				t.Errorf("FormatList(%q) == %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
