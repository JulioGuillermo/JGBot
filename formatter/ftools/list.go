package ftools

import (
	"regexp"
	"strings"
)

func FormatList(msg string) string {
	reTodo := regexp.MustCompile(`^(\*|-) \[ \] `)
	reDone := regexp.MustCompile(`(?i)^(\*|-) \[x\] `)
	reBullet := regexp.MustCompile(`^(\*|-) `)
	reNumbered := regexp.MustCompile(`^(\d+)\. `)

	lines := strings.Split(msg, "\n")
	for i, line := range lines {
		// To Do List
		if match := reTodo.FindStringSubmatch(line); len(match) > 1 {
			lines[i] = "- ⬜ " + line[len(match[0]):]
			continue
		}
		// Done List
		if match := reDone.FindStringSubmatch(line); len(match) > 1 {
			lines[i] = "- ✅ " + line[len(match[0]):]
			continue
		}
		// Bullet list
		if match := reBullet.FindStringSubmatch(line); len(match) > 1 {
			lines[i] = "- " + line[len(match[0]):]
			continue
		}
		// Numbered list
		if match := reNumbered.FindStringSubmatch(line); len(match) > 1 {
			lines[i] = match[1] + ". " + line[len(match[0]):]
			continue
		}
	}
	return strings.Join(lines, "\n")
}
