# Session Configuration Guide

This document explains the session management system in JGBot, covering the configuration of active and newly detected sessions.

## Overview

JGBot manages sessions using two distinct configuration files to separate customized settings from default ones:

1. **Configured Sessions** (`config/session.json`):
   - Contains sessions with explicitly defined settings.
2. **Unconfigured Sessions** (`config/unconfig_session.json`):
   - Contains new sessions using default parameters. This file is automatically created when a new session is detected, and new session are added to this file. This file is not used by the bot, but it is useful for tracking new sessions.

## Purpose of Each File

### Configured Sessions (`session.json`)

This file contains sessions that have been explicitly modified or approved.

**Key Features:**

- **Hot-Reloading**: Changes to this file are detected in real-time and applied without restarting the bot.
- **Persistence**: Stores custom settings for allowed status, history size, and specific LLM providers.
- **Granular Control**: Enable or disable specific tools and skills on a per-session basis.

### Unconfigured Sessions (`unconfig_session.json`)

When a new interaction is detected from a previously unknown origin (e.g., a new Telegram group or WhatsApp chat), JGBot automatically adds it to this file.

**Key Features:**

- **Automatic Discovery**: New sessions are captured without manual intervention.
- **Default Baseline**: Uses the global default settings for initial interaction.
- **Auto-Remove**: Any session added to `session.json` is automatically removed from this file.
- **Workflow**: Move an entry from this file to `session.json` to customize its behavior.

## Session Configuration Structure

### Example Session entry

```json
{
  "name": "Project Alpha Group",
  "id": "abc-123-def",
  "origin": "123456789@g.us",
  "allowed": true,
  "historySize": 50,
  "provider": "openai",
  "agentMaxIters": 10,
  "respond": {
    "always": true,
    "match": ".*"
  },
  "tools": [
     { "name": "message_reaction", "enabled": true },
     { "name": "cron", "enabled": true },
     { "name": "timer", "enabled": true },
     { "name": "javascript", "enabled": true },
     { "name": "skills", "enabled": true }
   ],
  "skills": [
    {
      "name": "weather_skill",
      "enabled": true,
      "description": "Provides real-time weather updates."
    }
  ]
}
```

### Configuration Options

#### Metadata
- **`name`**: The name of the chat. Allow user to identify the chat of this session, but is not used by the bot.
- **`id`**: A system-generated unique identifier.
- **`origin`**: A system-generated unique origin.

#### Access Control
- **`allowed`** (boolean): Set to `true` to allow the bot to process messages from this session.
  - **Default**: `false` (for unconfigured sessions).

#### Agent Behavior
- **`historySize`** (integer): The number of previous messages included in the context for the LLM.
- **`provider`** (string): The name of the LLM provider to use (must match a provider in `config.json`).
- **`agentMaxIters`** (integer): The maximum number of tool-use iterations the agent can perform per response.

#### Response Triggering
- **`respond.always`** (boolean): If `true`, the bot evaluates every message.
- **`respond.match`** (string): A regex pattern. The bot only responds if the message matches this pattern.

#### Tools & Skills
- **`tools`**: Control access to core system capabilities:
  - `message_reaction`: Allows the bot to react to messages with emojis.
  - `cron`: Enables scheduling repeat tasks and messages.
  - `timer`: Enables setting one-time alarms and timeouts.
  - `javascript`: Raw JavaScript execution for advanced logic.
  - `skills`: Access to the custom skill system.
- **`skills`**: A list of available custom skills. You can enable/disable them individually for each session.

## The Hot-Reload Workflow

JGBot continuously monitors `config/session.json`. This allows for seamless administration:

1. **Modify**: Open `config/session.json` and make your changes (e.g., toggle `allowed` to `true`).
2. **Save**: Save the file.
3. **Applied**: The bot immediately updates its internal state for that session. No restart is needed.

> [!CAUTION]
> If you introduce a JSON syntax error while editing, the hot-reload will fail, and an error will be logged. The bot will continue using the last valid configuration until the file is fixed.

## Session Lifecycle

1. **Detection**: An unknown chat sends a message. It appears in `unconfig_session.json`.
2. **Approval**: You copy the session entry into `session.json`.
3. **Configuration**: You set `allowed: true`, choose a provider, and enable specific skills.
4. **Maintenance**: You can update the history size or disable skills at any time via `session.json`.
5. **Removal**: Deleting a session from `session.json` stops the bot from responding; if that chat sends another message, it will reappear in `unconfig_session.json`.

## Best Practices

- **Security**: For public groups, keep `javascript` disabled unless absolutely necessary.
- **Resource Management**: Large `historySize` values consume more tokens and memory. Aim for a balance (typically 20–50).
- **Organization**: Use clear semantic names for sessions to simplify management as the number of chats grows.

