# JGBot Configuration Guide

This document provides a detailed overview of all configuration options available in JGBot.

## Configuration File Structure

The primary configuration is stored in `config/config.json`. The file is organized into several key sections:

### Overview

```json
{
  "Database": "path/to/database.db",
  "LogLevel": "Info",
  "Channels": {
    "Telegram": { ... },
    "Whatsapp": { ... }
  },
  "Providers": [
    { ... },
    { ... }
  ]
}
```

## Configuration Options

### Database Settings

- **`Database`** (string): Path to the SQLite database file.
  - **Default**: `"db/database.db"`
  - **Example**: `"data/mybot.db"`

### Logging Settings

- **`LogLevel`** (string): Defines the verbosity of the application logs.
  - **Options**: `"Debug"`, `"Info"`, `"Warn"`, `"Error"`
  - **Default**: `"Info"`
  - **Example**: `"Debug"`

### Channels Configuration

The `Channels` object contains settings for each supported communication platform.

#### Telegram Channel

```json
"Telegram": {
  "Enabled": true,
  "AutoEnableSession": false,
  "Config": {
    "Token": "your_telegram_bot_token"
  }
}
```

- **`Enabled`** (boolean): Set to `true` to activate the Telegram bot.
- **`AutoEnableSession`** (boolean): If `true`, new Telegram chats will automatically have sessions enabled without manual approval.
- **`Config.Token`** (string): Your official Telegram Bot API token. (Required if enabled).

#### WhatsApp Channel

```json
"Whatsapp": {
  "Enabled": true,
  "AutoEnableSession": false,
  "Config": {
    "DBPath": "path/to/whatsapp.db"
  }
}
```

- **`Enabled`** (boolean): Set to `true` to activate WhatsApp integration.
- **`AutoEnableSession`** (boolean): If `true`, new WhatsApp chats will automatically have sessions enabled.
- **`Config.DBPath`** (string): Path to the separate database used for WhatsApp session data.
  - **Default**: `"db/whatsapp.db"`

### LLM Providers Configuration

The `Providers` array allows you to configure multiple Large Language Model (LLM) backends.

#### Common Provider Fields

- **`Name`** (string): A unique identifier for the provider.
- **`Type`** (string): The provider type. Supported options: `"openai"`, `"anthropic"`, `"google"`, `"ollama"`, `"mistral"`.
- **`BaseUrl`** (string, optional): Custom API endpoint URL. If `null`, the default provider URL is used.
- **`ApiKey`** (string, optional): Your API key for the service.
- **`Model`** (string): The specific model name to use (e.g., `"gpt-4"`, `"claude-3-sonnet"`, `"gemini-pro"`).

#### Provider Examples

**OpenAI / Compatible:**
```json
{
  "Name": "openai",
  "Type": "openai",
  "ApiKey": "sk-...",
  "Model": "gpt-4"
}
```

**Anthropic:**
```json
{
  "Name": "anthropic",
  "Type": "anthropic",
  "ApiKey": "sk-ant-...",
  "Model": "claude-3-sonnet-20240229"
}
```

**Google Gemini:**
```json
{
  "Name": "google",
  "Type": "google",
  "ApiKey": "AIza...",
  "Model": "gemini-pro"
}
```

**Ollama (Local):**
```json
{
  "Name": "ollama",
  "Type": "ollama",
  "BaseUrl": "http://localhost:11434",
  "Model": "llama3"
}
```

## Setup and Management

### Initial Configuration

Upon the first execution, JGBot will:
1. Detect that `config/config.json` is missing.
2. Create a default configuration template.
3. Edit the configuration file to add your API tokens and other required information.
4. Restart the application.

### Validation

- Ensures all enabled channels have the necessary credentials.
- Verifies that file paths for databases are accessible.
- Checks that the provider configurations are structurally sound.

### Applying Changes

> [!NOTE]
> Most changes to `config/config.json` (such as adding providers or updating tokens) require an **application restart** to take effect. For session-specific overrides that support hot-reloading, please refer to [SESSION.md](SESSION.md).

