# JGBot

A modular AI-powered chatbot framework with multi-channel support, custom skills, and extensible architecture.

## Features

- **Multi-Channel Support**: Telegram, WhatsApp.
- **Custom Skills**: Extend bot functionality with skills and JavaScript skills tools.
- **AI Agents**: Powered by LangChainGo with support for multiple LLM providers.
- **Database Integration**: SQLite with GORM for data persistence.
- **Virtual File System**: Session-specific and shared file storage.
- **HTTP Client**: Built-in HTTP functionality for skills.
- **Modular Architecture**: Easy to extend and customize.

## Installation

### Prerequisites

- Go 1.25.5 or higher
- Git

### Setup

1. Clone the repository:

```bash
git clone https://github.com/JulioGuillermo/JGBot
cd JGBot
```

2. Install dependencies:

```bash
go mod tidy
```

3. Build the application:

```bash
go build
```

4. Run the application:

```bash
./JGBot
```

The first time you run the application, it will create all configuration files and new database file. You can specify the database file path in the configuration file.

## Custom Skills

JGBot supports custom skills that extend its capabilities. See [CUSTOM_SKILL.md](doc/CUSTOM_SKILL.md) for detailed documentation on creating custom skills.

### Quick Skill Example

Create a new skill in the `skills/` directory:

```
skills/my_skill/
├── SKILL.md
└── init.js
```

**SKILL.md**:

```yaml
---
name: my_skill
description: A custom skill example
---

# My Custom Skill

This skill demonstrates basic functionality.

## Inputs

{input}
```

**init.js**:

```javascript
const args = GetArgs();

async function main() {
  const { input } = args;
  return `You said: ${input}`;
}

await main();
```

## Architecture

### Core Components

1. **Agent System**: AI-powered conversation handling
2. **Channel Controller**: Manages different communication channels
3. **Session Manager**: Handles conversation sessions
4. **Skill System**: Loads and executes custom skills
5. **Database**: Persistent data storage

### Directory Structure

```
JGBot/
├── agent/          # AI agent implementation
├── channels/       # Channel implementations
├── conf/           # Configuration
├── database/       # Database models and connections
├── doc/            # Documentation
├── js/             # JavaScript runtime
├── plugins/        # Plugin system (coming soon)
├── session/        # Session management
├── skill/          # Skill system
├── skills/         # Custom skills directory
├── main.go         # Application entry point
└── go.mod          # Go module file
```

## Supported Channels

### Telegram

Built-in Telegram bot support with message handling, reactions, and file uploads.

### WhatsApp

WhatsApp integration for messaging functionality.

## AI Providers

JGBot supports multiple LLM providers:

- OpenAI (Can be used with any OpenAI compatible API)
- Anthropic
- Google
- Ollama (local models)
- Mistral

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## Support

For issues and questions:

- Check the documentation in `doc/`
- Review existing issues on GitHub
- Create a new issue if needed
