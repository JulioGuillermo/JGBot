package ftools

import (
	"strings"
)

// FormatHeader tries to format a line as a Header.
// Returns the formatted string and true if successful, else original string and false.
func FormatHeader(line string) (string, bool) {
	// 1️⃣2️⃣3️⃣4️⃣5️⃣6️⃣7️⃣8️⃣9️⃣🔟
	// h1
	if strings.HasPrefix(line, "# ") {
		return "1️⃣ " + strings.TrimPrefix(line, "# "), true
	}
	// h2
	if strings.HasPrefix(line, "## ") {
		return "2️⃣ " + strings.TrimPrefix(line, "## "), true
	}
	// h3
	if strings.HasPrefix(line, "### ") {
		return "3️⃣ " + strings.TrimPrefix(line, "### "), true
	}
	// h4
	if strings.HasPrefix(line, "#### ") {
		return "4️⃣ " + strings.TrimPrefix(line, "#### "), true
	}
	// h5
	if strings.HasPrefix(line, "##### ") {
		return "5️⃣ " + strings.TrimPrefix(line, "##### "), true
	}
	// h6
	if strings.HasPrefix(line, "###### ") {
		return "6️⃣ " + strings.TrimPrefix(line, "###### "), true
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

func FormatHeaderHTML(line string) (string, bool) {
	if strings.HasPrefix(line, "# ") {
		return "<h1>" + strings.TrimPrefix(line, "# ") + "</h1>", true
	}
	if strings.HasPrefix(line, "## ") {
		return "<h2>" + strings.TrimPrefix(line, "## ") + "</h2>", true
	}
	if strings.HasPrefix(line, "### ") {
		return "<h3>" + strings.TrimPrefix(line, "### ") + "</h3>", true
	}
	if strings.HasPrefix(line, "#### ") {
		return "<h4>" + strings.TrimPrefix(line, "#### ") + "</h4>", true
	}
	if strings.HasPrefix(line, "##### ") {
		return "<h5>" + strings.TrimPrefix(line, "##### ") + "</h5>", true
	}
	if strings.HasPrefix(line, "###### ") {
		return "<h6>" + strings.TrimPrefix(line, "###### ") + "</h6>", true
	}
	return line, false
}

func FormatHeadersHTML(msg string) string {
	lines := strings.Split(msg, "\n")
	for i, line := range lines {
		if h, ok := FormatHeaderHTML(line); ok {
			lines[i] = h
		}
	}
	return strings.Join(lines, "\n")
}
