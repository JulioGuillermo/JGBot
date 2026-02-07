---
name: memoirs
description: Manage persistent text memoirs (list, read, create, edit, save, delete).
metadata:
  author: JulioGuillermo
  version: "1.0"
---

# Memoirs Tool

This skill allows you to store, retrieve, and manage persistent text snippets called "memoirs".

## Usage

Use this tool when you need to remember information across the conversation or manage previously stored data.

### Arguments

- `command` (required): The action to perform. One of: `list`, `read`, `create`, `edit`, `save`, `delete`.
- `name` (optional): The name of the memory. Required for all commands except `list`.
- `content` (optional): The text content to store. Required for `create`, `edit`, and `save`.

### Operations

- `list`: Returns an alphabetical list of all memoir names.
- `read`: Returns the content of the specified memoir.
- `create`: Saves new content to a memoir. Fails if the memoir already exists.
- `edit`: Updates an existing memoir. Fails if the memoir does not exist.
- `save`: Saves content regardless of whether it exists (upsert).
- `delete`: Removes the specified memoir.

## Examples

- **Agent**: "I should remember that the user's favorite color is blue."
  **Action**: Call `memoirs` with `command="save"`, `name="favorite_color"`, `content="blue"`

- **Agent**: "List all my memoirs."
  **Action**: Call `memoirs` with `command="list"`

- **Agent**: "What was the user's favorite color?"
  **Action**: Call `memoirs` with `command="read"`, `name="favorite_color"`

## Constraints

- Memoir names should be concise and use underscores instead of spaces.
- Content is stored as markdown.
