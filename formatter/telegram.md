# Telegram Formatting Rules

Telegram supports several formatting options for messages, primarily through its **MarkdownV2** style and **HTML**. This document focuses on the Markdown-style syntax which is commonly used in bots and manual messages.

## Text Styles (MarkdownV2)

| Style | MarkdownV2 Syntax | HTML Syntax | Example |
| :--- | :--- | :--- | :--- |
| **Bold** | `*text*` or `**text**` | `<b>text</b>` | **bold** |
| *Italic* | `_text_` or `__text__` | `<i>text</i>` | *italic* |
| __Underline__ | `__text__` | `<u>text</u>` | <u>underline</u> |
| ~~Strikethrough~~ | `~text~` | `<s>text</s>` | ~~strikethrough~~ |
| || Spoiler || | `||text||` | `<tg-spoiler>text</tg-spoiler>` | [Spoiler] |
| `Inline Code` | `` `text` `` | `<code>text</code>` | `code` |
| `Block Code` | ` ```text``` ` | `<pre>text</pre>` | ```block``` |

## Links and Mentions

| Element | MarkdownV2 Syntax | HTML Syntax |
| :--- | :--- | :--- |
| **Inline Link** | `[text](url)` | `<a href="url">text</a>` |
| **User Mention** | `[text](tg://user?id=123)` | `<a href="tg://user?id=123">text</a>` |

## Mapping in this Package (Planned)

The `formatter` package aims to provide conversion to Telegram-compatible Markdown. 

> [!IMPORTANT]
> **MarkdownV2 Escaping:**
> Telegram's MarkdownV2 is very sensitive. Any of the following characters MUST be escaped with a backslash `\` if they are not part of a formatting tag:
> `_`, `*`, `[`, `]`, `(`, `)`, `~`, `` ` ``, `>`, `#`, `+`, `-`, `=`, `|`, `{`, `}`, `.`, `!`

### Potential Transformations:
- **Headings:** Telegram doesn't have native headers.
    - `# Title` → `*TITLE*` (Bold & Caps)
    - `## Subtitle` → `*Subtitle*` (Bold)
- **Lists:** Standard Markdown lists `-` or `1.` are generally kept as is, but require proper escaping of the dot `.` or hyphen `-`.
- **Tables:** Not natively supported. Likely converted to pre-formatted text blocks (```) or list-based representations.
