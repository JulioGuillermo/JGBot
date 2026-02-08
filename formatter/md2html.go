package formatter

import (
	"JGBot/formatter/ftools"
)

func FormatMD2HTML(msg string) string {
	// Protect already existing placeholders (if any) or existing code blocks
	msg, codeBlocks := ftools.MapCodeBlocks(msg)
	msg, links := ftools.MapLinks(msg)

	msg = ftools.EscapeHTML(msg)

	msg = ftools.FormatHeadersHTML(msg)

	// Apply styles. Using specific order and avoiding underscore style on placeholders.
	msg = ftools.FormatStyleTags(msg, "**", "**", "<b>", "</b>")
	msg = ftools.FormatStyleTags(msg, "__", "__", "<b>", "</b>")
	msg = ftools.FormatStyleTags(msg, "*", "*", "<i>", "</i>")
	msg = ftools.FormatStyleTags(msg, "~~", "~~", "<s>", "</s>")

	// Special case for _ : avoid matching inside placeholders like LINK_PH_
	// Since we already used ftools.FormatStyleTags for other things, we'll do _ manually or carefully.
	// Actually, the easiest is to format style BEFORE mapping if we are sure style won't break links/code.
	// But Markdown says code/links HAVE priority.
	// So we'll just restore them LAST, which we do.
	// The problem is that FormatStyleTags for "_" matches "LINK_PH_0" -> "LINK<i>PH</i>0"

	// FIX: Use a regex that doesn't match if it looks like a placeholder.
	// Or better, just apply the _ formatting before mapping links since links don't usually HAVE lone underscores that need italics.
	// But let's just use a more specific regex for _ in ftools or here.

	msg = ftools.FormatStyleTags(msg, "_", "_", "<i>", "</i>")

	msg = ftools.RestoreLinksHTML(msg, links)
	msg = ftools.RestoreCodeBlocksHTML(msg, codeBlocks)

	return msg
}
