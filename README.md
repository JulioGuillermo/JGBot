# JGBot

A modular, AI-powered chatbot framework featuring multi-channel support, a custom skill system, and an extensible architecture.

## Features

- **Multi-Channel Support**: Seamless integration with Telegram and WhatsApp.
- **Custom Skills**: Extend functionality using JavaScript-based skills and tools.
- **AI Agents**: Robust conversation handling powered by LangChainGo, supporting multiple LLM providers.
- **Cron Jobs**: Schedule recurring tasks and messages with flexible cron expressions.
- **Timers & Alarms**: Set one-time timeouts or specific alarms to trigger bot actions.
- **Message Reactions**: Built-in support for reacting to messages with emojis.
- **Database Integration**: Reliable data persistence using SQLite with GORM.
- **Virtual File System (VFS)**: Isolated file storage for sessions, including private and shared access.
- **HTTP Client**: Built-in, fluent HTTP functionality for external API interactions.
- **Modular & Extensible**: Designed for easy customization and growth.


## Installation

### Prerequisites

- **Go**: Version 1.25.5 or higher.
- **Git**: For repository management.

### Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/JulioGuillermo/JGBot
   cd JGBot
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Build the application**:
   ```bash
   go build
   ```

4. **Run the application**:
   ```bash
   ./JGBot
   ```

*Note: On the first run, the application generates default configuration files and a new database file. You can customize the database path in the configuration.*

## Custom Skills

JGBot enables you to extend its capabilities through custom skills. For a deep dive into skill development, see the [Custom Skills Development Guide](doc/CUSTOM_SKILL.md).

### Quick Skill Example

Create a new skill directory in `skills/`:

```
skills/my_skill/
├── SKILL.md
└── init.js
```

**SKILL.md** (Metadata and Context):
```yaml
---
name: my_skill
description: A simple demonstration skill
---

# Hello Skill
This skill greets the user and echoes their input.

## Usage
When the user wants to say hello or test the bot, use this skill.
Inputs: {input}
```

**init.js** (Logic):
```javascript
const args = GetArgs();

async function main() {
  const { input } = args;
  return `You said: ${input}`;
}

// The skill must export the result as default
export default await main();
```

## Architecture

### Core Components

1. **Agent System**: Manages AI-driven logic and conversation flows.
2. **Channel Controller**: Abstracts communication with different platforms (Telegram, WhatsApp).
3. **Session Manager**: Maintains state and history for individual conversations.
4. **Skill System**: Discovers, loads, and executes custom JavaScript skills.
5. **Database**: Handles persistent storage for sessions and configurations.

### Directory Structure

```
JGBot/
├── agent/          # AI agent and tool implementations
├── channels/       # Channel-specific logic (Telegram, WhatsApp)
├── conf/           # Configuration management
├── config/         # Persistent configuration files (sessions, cron, timers)
├── cron/           # Cron job system implementation
├── database/       # Database models and GORM setup
├── doc/            # Detailed documentation
├── js/             # JavaScript runtime and Go-JS bridge
├── plugins/        # Future plugin system
├── session/        # Session logic and lifecycle
├── skill/          # Skill loader and executor
├── skills/         # Directory for user-defined skills
├── timer/          # Timer and alarm system implementation
├── main.go         # Entry point
└── go.mod          # Go module definitions
```

## Supported Channels

### Telegram
Full-featured support including message handling, reactions, and file uploads.

### WhatsApp
Integration for messaging across WhatsApp chats and groups.

## AI Providers

JGBot supports a wide range of LLM providers via LangChainGo:

- **OpenAI**: Compatible with GPT-4, GPT-3.5, and OpenAI-compliant APIs.
- **Anthropic**: Support for Claude 3 models.
- **Google**: Integration with Gemini Pro.
- **Ollama**: For running local models.
- **Mistral**: Support for Mistral AI models.

## Contributing

We welcome contributions!
1. Fork the repository.
2. Create a feature branch.
3. Commit your changes.
4. Submit a Pull Request.

## Support

- Explore the detailed guides in the `doc/` directory:
  - [Custom Skills Development Guide](doc/CUSTOM_SKILL.md)
  - [Available Skills List](doc/AVAILABLE_SKILLS.md)
  - [Session Configuration Guide](doc/SESSION.md)
  - [Scheduled Tasks Guide (Cron & Timers)](doc/SCHEDULED_TASKS.md)
  - [Configuration Guide](doc/CONF.md)
- Review existing GitHub issues or create a new one for bugs and feature requests.

