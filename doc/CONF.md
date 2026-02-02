# CONF.md - JGBot Configuration Guide

This document describes all configuration options for the JGBot project.

## Configuration File Structure

The main configuration file is located at `config/config.json`. The configuration is divided into several main sections:

### Main Configuration

```json
{
  "database": "path/to/database.db",
  "logLevel": "Info",
  "channels": {
    "telegram": { ... },
    "whatsapp": { ... }
  },
  "providers": [
    { ... },
    { ... }
  ]
}
```

## Configuration Options

### Database Configuration

- **database** (string): Path to the SQLite database file
  - Default: `"db/database.db"`
  - Example: `"data/mybot.db"`

### Logging Configuration

- **logLevel** (string): Logging level for the application
  - Options: `"Debug"`, `"Info"`, `"Warn"`, `"Error"`
  - Default: `"Info"`
  - Example: `"Debug"`

### Channels Configuration

The `channels` object contains configurations for different communication channels:

#### Telegram Channel

```json
"telegram": {
  "enabled": true,
  "autoEnableSession": false,
  "config": {
    "token": "your_telegram_bot_token"
  }
}
```

- **enabled** (boolean): Whether the Telegram channel is enabled
  - Default: `false`
  - Example: `true`

- **autoEnableSession** (boolean): Whether to automatically enable sessions for new Telegram chats
  - Default: `false`
  - Example: `true`

- **config.token** (string): Telegram bot API token
  - Required: Yes (when enabled)
  - Example: `"123456789:ABCdefGHijKLmnoPqrsTuVwxyz123456"`

#### WhatsApp Channel

```json
"whatsapp": {
  "enabled": true,
  "autoEnableSession": false,
  "config": {
    "DBPath": "path/to/whatsapp.db"
  }
}
```

- **enabled** (boolean): Whether the WhatsApp channel is enabled
  - Default: `false`
  - Example: `true`

- **autoEnableSession** (boolean): Whether to automatically enable sessions for new WhatsApp chats
  - Default: `false`
  - Example: `true`

- **config.DBPath** (string): Path to the WhatsApp database file
  - Default: `"db/whatsapp.db"`
  - Example: `"data/whatsapp.db"`

### LLM Providers Configuration

The `providers` array contains configurations for different language model providers:

#### OpenAI Provider

```json
{
  "name": "openai",
  "type": "openai",
  "baseUrl": null,
  "apiKey": "sk-your-openai-api-key",
  "model": "gpt-4"
}
```

- **name** (string): Unique name for the provider
  - Required: Yes
  - Example: `"openai"`

- **type** (string): Provider type
  - Required: Yes
  - Options: `"openai"`, `"anthropic"`, `"google"`, `"ollama"`, `"mistral"`
  - Example: `"openai"`

- **baseUrl** (string, nullable): Custom API base URL (optional)
  - Default: `null` (uses default provider URL)
  - Example: `"https://api.openai.com/v1"`

- **apiKey** (string, nullable): API key for the provider
  - Required: Yes (for most providers)
  - Example: `"sk-your-api-key"`

- **model** (string, nullable): Model name to use
  - Required: Yes
  - Examples: `"gpt-4"`, `"gpt-3.5-turbo"`, `"claude-3-sonnet"`, `"gemini-pro"`

#### Anthropic Provider

```json
{
  "name": "anthropic",
  "type": "anthropic",
  "baseUrl": null,
  "apiKey": "sk-ant-api-key",
  "model": "claude-3-sonnet-20240229"
}
```

#### Google Provider

```json
{
  "name": "google",
  "type": "google",
  "baseUrl": null,
  "apiKey": "your-google-api-key",
  "model": "gemini-pro"
}
```

#### Ollama Provider

```json
{
  "name": "ollama",
  "type": "ollama",
  "baseUrl": "http://localhost:11434",
  "apiKey": null,
  "model": "llama2"
}
```

#### Mistral Provider

```json
{
  "name": "mistral",
  "type": "mistral",
  "baseUrl": null,
  "apiKey": "your-mistral-api-key",
  "model": "mistral-large-latest"
}
```

## Example Configuration File

```json
{
  "database": "data/bot.db",
  "logLevel": "Debug",
  "channels": {
    "telegram": {
      "enabled": true,
      "autoEnableSession": true,
      "config": {
        "token": "123456789:ABCdefGHijKLmnoPqrsTuVwxyz123456"
      }
    },
    "whatsapp": {
      "enabled": false,
      "autoEnableSession": false,
      "config": {
        "DBPath": "data/whatsapp.db"
      }
    }
  },
  "providers": [
    {
      "name": "openai",
      "type": "openai",
      "baseUrl": null,
      "apiKey": "sk-your-openai-api-key",
      "model": "gpt-4"
    },
    {
      "name": "anthropic",
      "type": "anthropic",
      "baseUrl": null,
      "apiKey": "sk-ant-api-key",
      "model": "claude-3-sonnet"
    }
  ]
}
```

## Configuration Management

### Automatic Configuration

When the application starts, it will:

1. Look for `config/config.json`
2. If it doesn't exist, create a default configuration
3. For missing required channel configurations, prompt the user to enter them
4. Save the updated configuration file

### Configuration Validation

The application validates configuration on startup:

- Checks that required channel configurations are present
- Validates API keys and model names
- Ensures file paths are accessible
- Verifies provider configurations are valid

### Configuration Persistence

- Configuration is automatically saved when modified
- Backups are not created automatically
- Configuration changes require application restart to take effect
