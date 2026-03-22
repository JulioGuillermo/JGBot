package sessiondomain

import (
	"strings"
)

const (
	ThinkStart = "<think>"
	ThinkEnd   = "</think>"
)

// ExtractReasoning separates <think> tags from the main content.
// It returns (cleanText, thinkingContent).
func ExtractReasoning(text string) (string, string) {
	var clean, think strings.Builder
	pos := 0
	done := false

	for pos < len(text) && !done {
		pos, done = parseThink(text, pos, &clean, &think)
	}

	return strings.TrimSpace(clean.String()), think.String()
}

func parseThink(text string, pos int, clean, think *strings.Builder) (int, bool) {
	start := strings.Index(text[pos:], ThinkStart)
	if start == -1 {
		clean.WriteString(text[pos:])
		return len(text), true
	}

	// 1. Add text leading up to the tag
	clean.WriteString(text[pos : pos+start])
	pos += start + len(ThinkStart)

	// 2. Locate the end of the thinking block
	end := strings.Index(text[pos:], ThinkEnd)

	// 3. Extract content (handle unclosed tags gracefully)
	var content string
	if end == -1 {
		content = text[pos:]
		pos = len(text)
	} else {
		content = text[pos : pos+end]
		pos += end + len(ThinkEnd)
	}

	// 4. Format and append to the thinking builder
	content = strings.TrimSpace(content)
	if content != "" {
		if think.Len() > 0 {
			think.WriteString("\n\n")
		}
		think.WriteString(content)
	}

	return pos, false
}
