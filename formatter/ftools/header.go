package ftools

import (
	"strings"
)

// FormatHeader tries to format a line as a Header.
// Returns the formatted string and true if successful, else original string and false.
func FormatHeader(line string) (string, bool) {
	// h1
	if strings.HasPrefix(line, "# ") {
		return "🔹 *" + strings.TrimPrefix(line, "# ") + "*", true
	}
	// h2
	if strings.HasPrefix(line, "## ") {
		return "🔹 *_" + strings.TrimPrefix(line, "## ") + "_*", true
	}
	// h3
	if strings.HasPrefix(line, "### ") {
		return "🔹 _" + strings.TrimPrefix(line, "### ") + "_", true
	}
	// h4
	if strings.HasPrefix(line, "#### ") {
		return "🔹 " + strings.TrimPrefix(line, "#### "), true
	}
	// h5
	if strings.HasPrefix(line, "##### ") {
		return "🔹 " + strings.TrimPrefix(line, "##### "), true
	}
	// h6
	if strings.HasPrefix(line, "###### ") {
		return "🔹 " + strings.TrimPrefix(line, "###### "), true
	}
	return line, false
}

func FormatHeaders(msg string) string {
	lines := strings.Split(msg, "\n")
	for i, line := range lines {
		if h, ok := FormatHeader(line); ok {
			lines[i] = h
		}
	}
	return strings.Join(lines, "\n")
}
