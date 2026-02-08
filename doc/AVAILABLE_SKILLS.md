# Available Custom Skills

This document lists the custom skills currently available in JGBot. These skills can be enabled or disabled per session in the `config/session.json` file.

## Core Skills

### [javascript](file:///home/jg/Documents/mispro/jgbot/skills/javascript/SKILL.md)
This skill provides the system instructions and technical reference for the the JavaScript runtime environment. It ensures the AI agent understands how to use the built-in `javascript` tool effectively, including the naming conventions (Go methods start with Uppercase) and available global objects like `VirtualFiles` and `HttpRequest`.

- **Tool Name**: `javascript` (native tool)
- **Primary Use**: Advanced logic, file operations, and manual HTTP requests.

### [memories](file:///home/jg/Documents/mispro/jgbot/skills/memories/SKILL.md)
Allows the agent to manage persistent text snippets. This is useful for remembering user preferences, facts, or any data that needs to persist across different conversations or after a restart.

- **Tool Name**: `memories`
- **Actions**: `list`, `read`, `create`, `edit`, `save`, `delete`.
- **Arguments**: `command`, `name`, `content`.

### [web](file:///home/jg/Documents/mispro/jgbot/skills/web/SKILL.md)
Enables the agent to access the internet. It can perform web searches or fetch the content of specific URLs, converting the results into clean Markdown.

- **Tool Name**: `web`
- **Arguments**: `query` (for search) or `url` (for direct fetching).

## Utility & Testing

### [test_skill](file:///home/jg/Documents/mispro/jgbot/skills/test/SKILL.md)
A simple skill used to verify that the agent can correctly pass arguments to local functions. It greets the user by name.

- **Tool Name**: `test_skill` (defined in the `test` directory)
- **Arguments**: `name`.

---

## How to Enable Skills

To use these skills, they must be enabled for the specific session in `config/session.json`.

Example configuration:

```json
{
  "name": "My Chat",
  "allowed": true,
  "skills": [
    { "name": "memories", "enabled": true },
    { "name": "web", "enabled": true },
    { "name": "javascript", "enabled": true }
  ]
}
```

> [!TIP]
> You can find detailed development instructions for each skill in its respective `SKILL.md` file within the `skills/` directory.
