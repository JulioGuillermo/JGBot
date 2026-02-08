package formatter

import (
	"testing"
)

func TestFormatMD2TelegramHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Combined Bold and Italic",
			input:    "**bold** *italic*",
			expected: "<b>bold</b> <i>italic</i>",
		},
		{
			name:     "Table and HTML",
			input:    "| H1 | H2 |\n|---|---|\n| V1 | V2 |",
			expected: "1. •••\n\t- <b>H1</b>: V1\n\t- <b>H2</b>: V2",
		},
		{
			name:     "Headers as H tags",
			input:    "# Header 1\n## Header 2",
			expected: "1️⃣ Header 1\n2️⃣ Header 2",
		},
		{
			name:     "Escaping special characters",
			input:    "Check this: 1 < 2 & 3 > 2",
			expected: "Check this: 1 &lt; 2 &amp; 3 &gt; 2",
		},
		{
			name:     "Mixed Formatting",
			input:    "# Title\n- [ ] Todo\n**Important** `code`",
			expected: "1️⃣ Title\n- ⬜ Todo\n<b>Important</b> <code>code</code>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatMD2TelegramHTML(tt.input)
			if got != tt.expected {
				t.Errorf("FormatMD2TelegramHTML() = %q, want %q", got, tt.expected)
			}
		})
	}
}
