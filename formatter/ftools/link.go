package ftools

import (
	"fmt"
	"regexp"
)

const LinkPlaceholder = "«LINKPH»"

// MapLinks finds [text](url) and raw URLs, replacing them with placeholders to protect them.
func MapLinks(msg string) (string, []string) {
	var links []string

	// 1. Markdown Links [text](url)
	reMD := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	msg = reMD.ReplaceAllStringFunc(msg, func(match string) string {
		links = append(links, match)
		return fmt.Sprintf("%s%d", LinkPlaceholder, len(links)-1)
	})

	// 2. Raw URLs
	// Simple regex for any protocol
	reRaw := regexp.MustCompile(`[a-zA-Z]+://[^\s]+`)
	msg = reRaw.ReplaceAllStringFunc(msg, func(match string) string {
		// Avoid double matching if it's already inside a placeholder?
		// Placeholders are LINK_PH_... no http there.
		links = append(links, match)
		return fmt.Sprintf("%s%d", LinkPlaceholder, len(links)-1)
	})

	return msg, links
}

// RestoreLinks restores links and formats them:
// [text](url) -> (text: url)
// Raw URL -> returns as is (protected from other formatting)
func RestoreLinks(msg string, links []string, supportMD bool) string {
	re := regexp.MustCompile(LinkPlaceholder + `(\d+)`)
	return re.ReplaceAllStringFunc(msg, func(match string) string {
		var idx int
		n, err := fmt.Sscanf(match, LinkPlaceholder+"%d", &idx)
		if err != nil || n != 1 {
			return match
		}

		if idx < 0 || idx >= len(links) {
			return match
		}

		content := links[idx]
		if content == "" {
			return match
		}

		if supportMD {
			return content
		}

		// Check if it is MD link or Raw
		if content[0] != '[' {
			return content
		}

		reMD := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
		sub := reMD.FindStringSubmatch(content)
		if len(sub) != 3 {
			return content
		}

		title := sub[1]
		url := sub[2]
		// Requested format: (title: url)
		return fmt.Sprintf("(%s: %s)", title, url)
	})
}

func RestoreLinksHTML(msg string, links []string) string {
	re := regexp.MustCompile(LinkPlaceholder + `(\d+)`)
	return re.ReplaceAllStringFunc(msg, func(match string) string {
		var idx int
		n, err := fmt.Sscanf(match, LinkPlaceholder+"%d", &idx)
		if err != nil || n != 1 {
			return match
		}

		if idx < 0 || idx >= len(links) {
			return match
		}

		content := links[idx]
		if content == "" {
			return match
		}

		// Check if it is MD link or Raw
		if content[0] != '[' {
			// Raw URL
			return fmt.Sprintf(`<a href="%s">%s</a>`, content, content)
		}

		reMD := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
		sub := reMD.FindStringSubmatch(content)
		if len(sub) != 3 {
			return content
		}

		title := sub[1]
		url := sub[2]
		return fmt.Sprintf(`<a href="%s">%s</a>`, url, title)
	})
}
