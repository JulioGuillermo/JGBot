# WhatsApp Formatting Rules

This document outlines the formatting rules supported by WhatsApp, which are used as the target for the Markdown conversion logic in this package.

## Text Styles

| Style | Markdown Syntax | WhatsApp Syntax | Example |
| :--- | :--- | :--- | :--- |
| **Bold** | `**text**` or `__text__` | `*text*` | *bold* |
| *Italic* | `*text*` or `_text_` | `_text_` | _italic_ |
| ~~Strikethrough~~ | `~~text~~` | `~text~` | ~strikethrough~ |
| `Monospace` | ` ```text``` ` | ` ```text``` ` | ```monospace``` |
| `Inline Code` | `` `text` `` | `` `text` `` | `code` |

## Blocks and Lists (2024+ Features)

| Element | Markdown Syntax | WhatsApp Syntax |
| :--- | :--- | :--- |
| **Bullet List** | `- Item` or `* Item` | `- Item` |
| **Numbered List** | `1. Item` | `1. Item` |
| **Blockquote** | `> Text` | `> Text` |

## Mapping in this Package

The `formatter` package performs the following mappings specifically for WhatsApp:

- **Headings:**
    - `# Title` → `🔹 *Title*`
    - `## Subtitle` → `🔹 *_Subtitle_*`
    - `### Section` → `🔹 _Section_`
- **Task Lists:**
    - `[ ]` → `白` (White square)
    - `[x]` → `✅` (Check mark)
- **Tables:** Transferred to a numbered list format where headers are bolded.
- **Links:** `[title](url)` → `(title: url)`.

> [!IMPORTANT]
> WhatsApp formatting requires the symbols to be placed immediately adjacent to the text without spaces (e.g., `*bold*` works, `* bold *` may not).
