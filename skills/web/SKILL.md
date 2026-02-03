---
name: get_link
description: Fetches the content of a URL and converts the HTML into Markdown format.
metadata:
  author: JulioGuillermo
  version: "1.0"
---

# Get Link

This skill allows the bot to access external web content. It takes a URL, fetches the HTML content of the page, and converts it into a clean Markdown format for easier processing and reading.

## Usage

Use this tool whenever a user provides a URL and asks to read, summarize, analyze, or extract information from that webpage.

### Arguments

- `url`: The full web address (including http:// or https://) of the page you want to fetch.

## Examples

- **User**: "Summarize the information on this page: https://example.com/article"
  **Action**: Call `get_link` with `url="https://example.com/article"`
- **User**: "What are the main points of https://en.wikipedia.org/wiki/JavaScript?"
  **Action**: Call `get_link` with `url="https://en.wikipedia.org/wiki/JavaScript"`

## Notes

The tool will return a string containing the Markdown representation of the page. If the request fails (e.g., 404 error or timeout), it will return an error message describing the status code.
