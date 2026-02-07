package formatter

import (
	"testing"
)

func TestMD2WhatsApp(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// 1. Basic Text Styles
		{name: "Bold Markdown **", input: "**bold text**", expected: "*bold text*"},
		{name: "Bold Markdown __", input: "__bold text__", expected: "*bold text*"},
		{name: "Italic Markdown *", input: "*italic text*", expected: "_italic text_"},
		{name: "Italic Markdown _", input: "_italic text_", expected: "_italic text_"},
		{name: "Strikethrough Markdown ~~", input: "~~strikethrough~~", expected: "~strikethrough~"},
		{name: "Inline Code `", input: "`code`", expected: "`code`"},
		{name: "Block Code ```", input: "```\nblock code\n```", expected: "```\nblock code\n```"},

		// 2. Lists (2024 Features)
		{name: "Bullet List -", input: "- Item 1\n- Item 2", expected: "- Item 1\n- Item 2"},
		{name: "Bullet List *", input: "* Item 1\n* Item 2", expected: "- Item 1\n- Item 2"},
		{name: "Numbered List", input: "1. First\n2. Second", expected: "1. First\n2. Second"},

		// 3. Blockquotes (2024 Feature)
		{name: "Blockquote", input: "> This is a quote", expected: "> This is a quote"},

		// 4. Nested Styles
		{name: "Bold and Italic 1", input: "**bold _italic_** ", expected: "*bold _italic_* "},
		{name: "Bold and Italic 2", input: "__bold *italic*__ ", expected: "*bold _italic_* "},
		{name: "Strikethrough in Bold", input: "**bold ~~strike~~**", expected: "*bold ~strike~*"},
		{name: "Bold in Italic Block", input: "_Italic with **bold** content_", expected: "_Italic with *bold* content_"},

		// 5. Headings (Custom mapping to Emojis)
		{name: "H1 Header", input: "# Title", expected: "1️⃣ Title"},
		{name: "H2 Header", input: "## Subtitle", expected: "2️⃣ Subtitle"},
		{name: "H3 Header", input: "### Section", expected: "3️⃣ Section"},
		{name: "H4 Header", input: "#### Subsection", expected: "4️⃣ Subsection"},
		{name: "H5 Header", input: "##### Subsection", expected: "5️⃣ Subsection"},
		{name: "H6 Header", input: "###### Subsection", expected: "6️⃣ Subsection"},

		// 6. Task Lists (Custom mapping)
		{name: "Task List Todo", input: "- [ ] To do", expected: "- ⬜ To do"},
		{name: "Task List Done", input: "- [x] Done", expected: "- ✅ Done"},

		// 7. Links (Custom mapping to (title: url))
		{name: "Markdown Link", input: "[Google](https://google.com)", expected: "(Google: https://google.com)"},
		{name: "Raw URL", input: "Check https://google.com", expected: "Check https://google.com"},

		// 8. Tables (Custom mapping to numbered list)
		{
			name:     "Markdown Table",
			input:    "| Name | Age |\n|---|---|\n| Bob | 30 |",
			expected: "1. •••\n\t- *Name*: Bob\n\t- *Age*: 30",
		},

		// 9. Edge Cases
		// {name: "Escaped Character", input: "\\*not bold\\*", expected: "*not bold*"}, // MD escapes usually just render the char
		{name: "Empty String", input: "", expected: ""},
		{name: "Multiple Newlines", input: "Line 1\n\n\nLine 2", expected: "Line 1\n\n\nLine 2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatMD2WhatsApp(tt.input)
			if got != tt.expected {
				t.Errorf("FormatMD2WhatsApp(%q) == %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
