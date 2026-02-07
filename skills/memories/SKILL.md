---
name: memories
description: Manage persistent text memories (list, read, create, edit, save, delete).
metadata:
  author: JulioGuillermo
  version: "1.0"
---

# Memories Tool

This skill allows you to store, retrieve, and manage persistent text snippets called "memories".

## Usage

Use this tool when you need to remember information across the conversation or manage previously stored data.

### Arguments

- `command` (required): The action to perform. One of: `list`, `read`, `create`, `edit`, `save`, `delete`.
- `name` (optional): The name of the memory. Required for all commands except `list`.
- `content` (optional): The text content to store. Required for `create`, `edit`, and `save`.

### Operations

- `list`: Returns an alphabetical list of all memory names.
- `read`: Returns the content of the specified memory.
- `create`: Saves new content to a memory. Fails if the memory already exists.
- `edit`: Updates an existing memory. Fails if the memory does not exist.
- `save`: Saves content regardless of whether it exists (upsert).
- `delete`: Removes the specified memory.

## Examples

- **Agent**: "I should remember that the user's favorite color is blue."
  **Action**: Call `memories` with `command="save"`, `name="favorite_color"`, `content="blue"`

- **Agent**: "List all my memories."
  **Action**: Call `memories` with `command="list"`

- **Agent**: "What was the user's favorite color?"
  **Action**: Call `memories` with `command="read"`, `name="favorite_color"`

## Constraints

- Memory names should be concise and use underscores instead of spaces.
- Content is stored as markdown.
