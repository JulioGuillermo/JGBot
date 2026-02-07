package ftools

import "testing"

func TestFormatList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "To Do List",
			input:    "- [ ] To do",
			expected: "- ⬜ To do",
		},
		{
			name:     "To Do List",
			input:    "* [ ] To do",
			expected: "- ⬜ To do",
		},
		{
			name:     "Done List",
			input:    "- [x] Done",
			expected: "- ✅ Done",
		},
		{
			name:     "Done List",
			input:    "* [x] Done",
			expected: "- ✅ Done",
		},
		{
			name:     "Bullet List",
			input:    "- Item 1\n- Item 2",
			expected: "- Item 1\n- Item 2",
		},
		{
			name:     "Bullet List",
			input:    "* Item 1\n* Item 2",
			expected: "- Item 1\n- Item 2",
		},
		{
			name:     "Numbered List",
			input:    "1. First\n2. Second",
			expected: "1. First\n2. Second",
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
