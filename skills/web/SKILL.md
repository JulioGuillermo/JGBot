---
name: web
description: Performs a web search or fetches the content of a URL, converting the result into Markdown format.
metadata:
  author: JulioGuillermo
  version: "1.1"
---

# Web Tool

This skill allows you to access external web content either by searching or by direct URL access. It fetches the content and converts it into a clean Markdown format.

## Usage

Use this tool when you need to:
1.  **Search the web**: Provide a `query` to find information.
2.  **Read a webpage**: Provide a `url` to read specific content.

### Arguments

- `query` (optional): The search terms to find information on the web.
- `url` (optional): The full web address (including http:// or https://) of the page to fetch.

**Note**: You must provide either `query` or `url`, but not both (url takes precedence).

## Examples

- **User**: "Who won the super bowl 2024?"
  **Action**: Call `web` with `query="super bowl 2024 winner"`
- **User**: "Summarize this article: https://example.com/article"
  **Action**: Call `web` with `url="https://example.com/article"`

## Notes

The tool returns a Markdown string. If using `query`, it performs a search (e.g., via Mojeek) and returns the results. If using `url`, it fetches and converts the page content.
