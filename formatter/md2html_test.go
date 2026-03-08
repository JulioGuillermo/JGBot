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
		{
			name:     "Underscore inside words should not become italic",
			input:    "**ak_arch** and **ak_arch_frontend**",
			expected: "<b>ak_arch</b> and <b>ak_arch_frontend</b>",
		},
		{
			name:     "Unmatched delimiters stay literal",
			input:    "**bold *italic _under",
			expected: "**bold *italic _under",
		},
		{
			name:     "Inline code protects markdown markers",
			input:    "`**bold** _italic_` outside **bold**",
			expected: "<code>**bold** _italic_</code> outside <b>bold</b>",
		},
		{
			name:     "Fenced code block protects markdown and html",
			input:    "```\n**bold** <tag>\n```",
			expected: "<pre><code>**bold** &lt;tag&gt;</code></pre>",
		},
		{
			name:     "Raw URL becomes anchor",
			input:    "Visit https://example.com/docs?q=1&lang=en",
			expected: `Visit <a href="https://example.com/docs?q=1&lang=en">https://example.com/docs?q=1&lang=en</a>`,
		},
		{
			name:     "Markdown link title is protected from style parsing",
			input:    "[**Docs**](https://example.com)",
			expected: `<a href="https://example.com">**Docs**</a>`,
		},
		{
			name:     "Header content can still be styled",
			input:    "# **Important** _now_",
			expected: "<h1><b>Important</b> <i>now</i></h1>",
		},
		{
			name:     "Italic underscore keeps punctuation boundaries",
			input:    "Start _italic_, end.",
			expected: "Start <i>italic</i>, end.",
		},
		{
			name:     "Multiple raw urls are restored independently",
			input:    "A https://a.test B https://b.test",
			expected: `A <a href="https://a.test">https://a.test</a> B <a href="https://b.test">https://b.test</a>`,
		},
		{
			name:     "Escaped html around markdown and code",
			input:    "<b>unsafe</b> and `x < y`",
			expected: "&lt;b&gt;unsafe&lt;/b&gt; and <code>x &lt; y</code>",
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
