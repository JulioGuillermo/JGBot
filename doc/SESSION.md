# SESSION.md - JGBot Session Configuration Guide

This document describes the session configuration system for the JGBot project, covering both configured and unconfigured sessions.

## Overview

The JGBot uses two separate configuration files to manage session settings:

1. **Configured Sessions** (`config/session.json`) - Sessions with custom configurations
2. **Unconfigured Sessions** (`config/unconfig_session.json`) - Sessions using default settings

## File Locations

- **Configured Sessions**: `config/session.json`
- **Unconfigured Sessions**: `config/unconfig_session.json`

## Purpose of Each File

### Configured Sessions (`session.json`)

This file stores sessions that have been explicitly configured by the user or administrator. Each session in this file has custom settings and overrides the default behavior.

**Key Features:**

- Hot-reload support (changes are automatically detected and applied)
- Persistent custom configurations
- Session-specific settings like allowed status, history size, provider selection, etc.
- Tools and skills configuration per session

### Unconfigured Sessions (`unconfig_session.json`)

This file serves as a temporary storage for sessions that haven't been explicitly configured yet. When a new session is detected, it's initially created here with default settings.

**Key Features:**

- Default configuration for new sessions
- Temporary storage before user configuration
- Sessions are automatically moved to `session.json` when configured
- Allows users to see and manage unconfigured sessions

## Session Configuration Structure

### Basic Session Structure

```json
{
  "name": "Session Name",
  "id": "unique_session_id",
  "origin": "channel_chat_identifier",
  "allowed": true,
  "historySize": 50,
  "provider": "openai",
  "agentMaxIters": 3,
  "respond": {
    "always": true,
    "match": ".*"
  },
  "tools": [
    {
      "name": "message_reaction",
      "enabled": true
    },
    {
      "name": "javascript",
      "enabled": false
    },
    {
      "name": "skills",
      "enabled": false
    }
  ],
  "skills": [
    {
      "name": "test_skill",
      "enabled": false,
      "description": "A testing skill that greets the user."
    }
  ]
}
```

### Configuration Options

#### Basic Information

- **name** (string): Display name for the session
- **id** (string): Unique identifier for the session
- **origin** (string): Channel-specific chat identifier (e.g., Telegram chat ID, WhatsApp JID)

#### Session Control

- **allowed** (boolean): Whether the session is allowed to interact with the bot
  - Default: `false`
  - Example: `true`

#### Agent Configuration

- **historySize** (integer): Number of messages to keep in conversation history
  - Default: `50`
  - Example: `100`

- **provider** (string): LLM provider to use for this session
  - Options: `"openai"`, `"anthropic"`, `"google"`, `"ollama"`, `"mistral"`
  - Example: `"anthropic"`

- **agentMaxIters** (integer): Maximum iterations for the agent to respond
  - Default: `3`
  - Example: `5`

#### Response Configuration

- **respond.always** (boolean): Whether the bot should always respond
  - Default: `true`
  - Example: `false`

- **respond.match** (string): Regex pattern for messages to respond to
  - Default: `.*` (respond to all messages)
  - Example: `"^hello|hi$"` (only respond to greetings)

#### Tools Configuration

- **tools** (array): Available tools for the session
  - **message_reaction**: Message reaction tool
    - `enabled`: Whether to allow message reactions
  - **javascript**: JavaScript execution tool
    - `enabled`: Whether to allow JavaScript execution
  - **skills**: Skills execution tool
    - `enabled`: Whether to allow skills execution

#### Skills Configuration

- **skills** (array): Available skills for the session
  - **name**: Skill name
  - **enabled**: Whether the skill is enabled for this session
  - **description**: Description of the skill

## Hot-Reload Feature

The configured sessions file (`session.json`) supports hot-reloading:

### How It Works

1. **File Watching**: The system monitors `session.json` for changes
2. **Automatic Reload**: When changes are detected, the configuration is automatically reloaded
3. **Session Update**: All active sessions using the configuration are updated immediately
4. **No Restart Required**: Changes take effect without restarting the bot

### Supported Changes

- Adding new sessions
- Modifying existing session settings
- Enabling/disabling tools and skills
- Changing response patterns
- Updating allowed status

### Limitations

- Session removal requires active session cleanup
- Large file changes may cause brief delays
- Invalid JSON will cause reload errors

## Session Lifecycle

### 1. New Session Detection

When a new chat is detected:

- Session is created in `unconfig_session.json` with default settings
- User can see it in the unconfigured sessions list

### 2. Session Configuration

User can configure a session by:

- Moving it from `unconfig_session.json` to `session.json`
- Setting custom parameters
- Enabling/disabling specific tools and skills

### 3. Active Session

Once configured:

- Session uses custom settings from `session.json`
- Benefits from hot-reload for configuration updates
- Can be modified at any time

### 4. Session Removal

Sessions can be removed by:

- Deleting them from `session.json`
- System will automatically recreate them in `unconfig_session.json` if they become active again

## Default Configuration

When a session is first created in `unconfig_session.json`, it uses these defaults:

```json
{
  "name": "New Session",
  "id": "generated_id",
  "origin": "channel_origin",
  "allowed": false,
  "historySize": 50,
  "provider": "default_from_main_config",
  "agentMaxIters": 3,
  "respond": {
    "always": true,
    "match": ".*"
  },
  "tools": [
    {
      "name": "message_reaction",
      "enabled": true
    },
    {
      "name": "javascript",
      "enabled": false
    },
    {
      "name": "skills",
      "enabled": false
    }
  ],
  "skills": [
    {
      "name": "skill_name",
      "enabled": false,
      "description": "Skill description"
    }
  ]
}
```

## Management Commands

### View Unconfigured Sessions

```bash
cat config/unconfig_session.json
```

### View Configured Sessions

```bash
cat config/session.json
```

### Configure a Session

1. Edit `config/session.json` to add/modify session configuration
2. Save the file (hot-reload will automatically apply changes)

### Reset Session to Defaults

1. Remove the session from `session.json`
2. It will reappear in `unconfig_session.json` with defaults

## Best Practices

### Security Considerations

- Set `allowed: false` for sessions that shouldn't interact
- Disable JavaScript execution for untrusted sessions
- Use specific response patterns to limit interactions

### Performance Optimization

- Keep `historySize` reasonable for memory usage
- Disable unused tools to reduce processing overhead
- Use appropriate `agentMaxIters` to balance response quality vs performance

### Management Tips

- Regularly review and clean up unconfigured sessions
- Use descriptive session names for easier management
- Test response patterns before deploying to production
