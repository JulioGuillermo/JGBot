package formatter

import (
	"strings"

	"JGBot/formatter/ftools"
)

func FormatMD2Telegram(msg string) string {
	// 1. Task List
	msg = strings.ReplaceAll(msg, "- [ ] ", "- ⬜ ")
	msg = strings.ReplaceAll(msg, "- [x] ", "- ✅ ")

	// 2. Tables (using ftools)
	msg = ftools.FormatTable(msg)

	// 3. Headers & Styles
	lines := strings.Split(msg, "\n")
	for i, line := range lines {
		// Try H1: # Title -> 🔹 **Title**
		if strings.HasPrefix(line, "# ") {
			lines[i] = "🔹 **" + strings.TrimPrefix(line, "# ") + "**"
			continue
		}
		// Try H2: ## Subtitle -> 🔹 __Subtitle__
		if strings.HasPrefix(line, "## ") {
			lines[i] = "🔹 __" + strings.TrimPrefix(line, "## ") + "__"
			continue
		}
		// Try H3: ### Section -> 🔹 Section
		if strings.HasPrefix(line, "### ") {
			lines[i] = "🔹 " + strings.TrimPrefix(line, "### ")
			continue
		}
		// Try H4: #### Subsection -> 🔹 Subsection
		if strings.HasPrefix(line, "#### ") {
			lines[i] = "🔹 " + strings.TrimPrefix(line, "#### ")
			continue
		}
	}
	msg = strings.Join(lines, "\n")

	return msg
}
