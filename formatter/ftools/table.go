package ftools

import (
	"fmt"
	"strings"
)

// FormatTable formats all tables in a message string.
func FormatTable(msg string) string {
	lines := strings.Split(msg, "\n")
	newLines := ProcessTableBlock(lines)
	return strings.Join(newLines, "\n")
}

// tableContext holds state for table processing
type tableContext struct {
	inTable  bool
	headers  []string
	rowCount int
}

func (tc *tableContext) reset() {
	tc.inTable = false
	tc.headers = nil
	tc.rowCount = 0
}

// ProcessTableBlock processes a block of lines to find and transform tables.
func ProcessTableBlock(lines []string) []string {
	var newLines []string
	ctx := &tableContext{}

	for i, line := range lines {
		if !IsTableLine(line) {
			ctx.reset()
			newLines = append(newLines, line)
			continue
		}

		// It is a table line
		cols := ExtractTableCells(line)

		// Case 1: Start of a new table (Header)
		if !ctx.inTable {
			if isHeader(lines, i) {
				ctx.inTable = true
				ctx.headers = cols
				ctx.rowCount = 0
				// Skip printing header line, we use it for rows
				// We DO NOT skip the separator line here explicitly,
				// the loop will hit it next and handle it in Case 2.
				continue
			}
			// Not a valid table header (no separator next), treat as text
			newLines = append(newLines, line)
			continue
		}

		// Case 2: Separator line
		if strings.Contains(line, "---") {
			continue // Skip
		}

		// Case 3: Data row
		ctx.rowCount++
		newLines = append(newLines, formatTableRow(cols, ctx.headers, ctx.rowCount)...)
	}

	return newLines
}

func isHeader(lines []string, currentIndex int) bool {
	if currentIndex+1 >= len(lines) {
		return false
	}
	nextLine := lines[currentIndex+1]
	return IsTableLine(nextLine) && strings.Contains(nextLine, "---")
}

func formatTableRow(cols []string, headers []string, rowNum int) []string {
	var rows []string
	rows = append(rows, fmt.Sprintf("%d. •••", rowNum))
	for j, val := range cols {
		header := ""
		if j < len(headers) {
			header = headers[j]
		}
		rows = append(rows, fmt.Sprintf("\t- **%s**: %s", header, val))
	}
	return rows
}

// Helper functions (kept from before, ensured simple)

func IsTableLine(line string) bool {
	trim := strings.TrimSpace(line)
	return strings.HasPrefix(trim, "|") && strings.HasSuffix(trim, "|")
}

func ExtractTableCells(line string) []string {
	trim := strings.TrimSpace(line)
	cols := strings.Split(trim, "|")
	if len(cols) <= 2 {
		return nil
	}
	contentCols := cols[1 : len(cols)-1]
	for k := range contentCols {
		contentCols[k] = strings.TrimSpace(contentCols[k])
	}
	return contentCols
}
