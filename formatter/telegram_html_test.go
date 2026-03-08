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
		{
			name:     "No italic for underscores in identifiers",
			input:    "4. **ak_arch** & **ak_arch_frontend** - stack",
			expected: "4. <b>ak_arch</b> &amp; <b>ak_arch_frontend</b> - stack",
		},
		{
			name:     "Uppercase task list done",
			input:    "- [X] Done",
			expected: "- ✅ Done",
		},
		{
			name:     "Raw url becomes anchor after telegram preprocessing",
			input:    "See https://example.com",
			expected: `See <a href="https://example.com">https://example.com</a>`,
		},
		{
			name:     "Inline code protects nested markdown",
			input:    "`**bold** _italic_` and _real_",
			expected: "<code>**bold** _italic_</code> and <i>real</i>",
		},
		{
			name:     "Fenced code block survives telegram html conversion",
			input:    "```\n# title\n**bold**\n```",
			expected: "<pre><code># title\n**bold**</code></pre>",
		},
		{
			name:     "Two row table renders both rows",
			input:    "| Name | Age |\n|---|---|\n| Bob | 30 |\n| Ana | 25 |",
			expected: "1. •••\n\t- <b>Name</b>: Bob\n\t- <b>Age</b>: 30\n2. •••\n\t- <b>Name</b>: Ana\n\t- <b>Age</b>: 25",
		},
		{
			name:     "Header list and emphasis together",
			input:    "# Title\n- **bold**\n- _italic_",
			expected: "1️⃣ Title\n- <b>bold</b>\n- <i>italic</i>",
		},
		{
			name:     "Markdown link preserved then converted to html anchor",
			input:    "[docs](https://example.com/path?a=1&b=2)",
			expected: `<a href="https://example.com/path?a=1&b=2">docs</a>`,
		},
		{
			name:     "Unmatched delimiters stay literal",
			input:    "**bold _italic",
			expected: "**bold _italic",
		},
		{
			name:     "Header content can include emphasis",
			input:    "## **Build** _status_",
			expected: "2️⃣ <b>Build</b> <i>status</i>",
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
