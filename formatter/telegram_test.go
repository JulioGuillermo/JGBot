package formatter

import (
	"testing"
)

func TestMD2Telegram(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic styles are supported by telegram from md
		// Basic list are supported by telegram from md
		// Blockquote are supported by telegram from md
		// URL are supported by telegram from md

		// 1. Task List
		{name: "Task List Todo", input: "- [ ] To do", expected: "- ⬜ To do"},
		{name: "Task List Done", input: "- [x] Done", expected: "- ✅ Done"},

		// 2. Tables (Custom mapping to numbered list)
		{
			name:     "Markdown Table",
			input:    "| Name | Age |\n|---|---|\n| Bob | 30 |",
			expected: "1. •••\n\t- **Name**: Bob\n\t- **Age**: 30",
		},

		// 3. Headings (Custom mapping for Telegram)
		{name: "H1 Header", input: "# Title", expected: "1️⃣ Title"},
		{name: "H2 Header", input: "## Subtitle", expected: "2️⃣ Subtitle"},
		{name: "H3 Header", input: "### Section", expected: "3️⃣ Section"},
		{name: "H4 Header", input: "#### Subsection", expected: "4️⃣ Subsection"},
		{name: "H5 Header", input: "##### Subsection", expected: "5️⃣ Subsection"},
		{name: "H6 Header", input: "###### Subsection", expected: "6️⃣ Subsection"},

		// 4. Mixed Content
		{
			name:     "Complex Message",
			input:    "# Header\n- Item 1\n- Item 2\nCheck [link](https://test.com)!",
			expected: "1️⃣ Header\n- Item 1\n- Item 2\nCheck [link](https://test.com)!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatMD2Telegram(tt.input)
			if got != tt.expected {
				t.Errorf("FormatMD2Telegram(%q) == %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
