# Available Custom Skills

This document lists the custom skills and built-in tools available in JGBot. These skills can be enabled or disabled per session in the `config/session.json` file.

## Native Tools (Built-in)

These are core system capabilities that the AI agent can use:

### [message_reaction](reactiontool/reaction.go)

Allows the agent to react to specific messages with emojis.

- **Tool Name**: `message_reaction`
- **Arguments**: `message_id`, `reaction` (emoji string)
- **Session Control**: Can be enabled/disabled per session

### [javascript](skills/javascript/SKILL.md)

This skill provides the system instructions and technical reference for the JavaScript runtime environment. It ensures the AI agent understands how to use the built-in `javascript` tool effectively.

- **Tool Name**: `javascript` (native tool)
- **Primary Use**: Advanced logic, file operations, and manual HTTP requests.
- **Session Control**: Can be enabled/disabled per session (recommended: disable for public groups)

### [skills](skills/)

Access to the custom skill system. Allows the agent to use loaded custom skills.

- **Tool Name**: `skills` (native tool)
- **Session Control**: Can be enabled/disabled per session

### [cron](SCHEDULED_TASKS.md)

Enables scheduling recurring tasks using cron expressions.

- **Tool Name**: `cron`
- **Actions**: `list`, `read`, `add`, `remove`
- **Session Control**: Can be enabled/disabled per session

### [timer](SCHEDULED_TASKS.md)

Enables setting one-time alarms and timeouts.

- **Tool Name**: `timer`
- **Actions**: `list`, `read`, `add`, `remove`
- **Session Control**: Can be enabled/disabled per session

## Custom Skills

### [memories](skills/memories/SKILL.md)

Allows the agent to manage persistent text snippets. This is useful for remembering user preferences, facts, or any data that needs to persist across different conversations or after a restart.

- **Tool Name**: `memories`
- **Actions**: `list`, `read`, `create`, `edit`, `save`, `delete`
- **Arguments**: `command`, `name`, `content`

### [web](skills/web/SKILL.md)

Enables the agent to access the internet. It can perform web searches or fetch the content of specific URLs, converting the results into clean Markdown.

- **Tool Name**: `web`
- **Arguments**: `query` (for search) or `url` (for direct fetching)

## Utility & Testing

### [test_skill](skills/test/SKILL.md)

A simple skill used to verify that the agent can correctly pass arguments to local functions. It greets the user by name.

- **Tool Name**: `test_skill`
- **Arguments**: `name`

## Administrative Tools

These tools require administrator permissions:

### send_message

Allows sending messages to other sessions from within the agent.

- **Tool Name**: `send_message`
- **Arguments**: `session` (origin), `message`
- **Permission**: Admin only

### list_sessions

Lists all configured sessions.

- **Tool Name**: `list_sessions`
- **Permission**: Admin only

---

## How to Enable Skills

To use these skills, they must be enabled for the specific session in `config/session.json`.

Example configuration:

```json
{
  "name": "My Chat",
  "allowed": true,
  "tools": [
    { "name": "message_reaction", "enabled": true },
    { "name": "cron", "enabled": true },
    { "name": "timer", "enabled": true },
    { "name": "javascript", "enabled": false },
    { "name": "skills", "enabled": true }
  ],
  "skills": [
    { "name": "memories", "enabled": true },
    { "name": "javascript", "enabled": false },
    { "name": "web", "enabled": true }
  ]
}
```

> [!TIP]
> You can find detailed development instructions for each skill in its respective `SKILL.md` file within the `skills/` directory.
