package formatter

import (
	"testing"
)

func TestFormatMD2HTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Bold and Italic",
			input:    "**bold** *italic*",
			expected: "<b>bold</b> <i>italic</i>",
		},
		{
			name:     "Bold and Italic2",
			input:    "**bold _italic_**",
			expected: "<b>bold <i>italic</i></b>",
		},
		{
			name:     "Bold and Italic2",
			input:    "*bold __italic__*",
			expected: "<i>bold <b>italic</b></i>",
		},
		{
			name:     "Headers",
			input:    "# H1\n## H2",
			expected: "<h1>H1</h1>\n<h2>H2</h2>",
		},
		{
			name:     "HTML Escaping",
			input:    "1 < 2 & 3 > 2",
			expected: "1 &lt; 2 &amp; 3 &gt; 2",
		},
		{
			name:     "Markdown Links",
			input:    "[text](https://example.com)",
			expected: `<a href="https://example.com">text</a>`,
		},
		{
			name:     "Code Blocks",
			input:    "```\ncode\n```",
			expected: "<pre><code>code</code></pre>",
		},
		{
			name:     "Inline Code",
			input:    "`code`",
			expected: "<code>code</code>",
		},
		{
			name:     "Strike",
			input:    "~~strike~~",
			expected: "<s>strike</s>",
		},
		{
			name:     "Mixed HTML and MD", // MD should be escaped if not intentional
			input:    "**bold** <script>alert(1)</script>",
			expected: "<b>bold</b> &lt;script&gt;alert(1)&lt;/script&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatMD2HTML(tt.input)
			if got != tt.expected {
				t.Errorf("FormatMD2HTML() = %q, want %q", got, tt.expected)
			}
		})
	}
}
