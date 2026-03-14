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

#### Channel Default Configuration (DefConf)

Each channel can have a default configuration that applies to all new sessions from that channel. This overrides the global default configuration.

```json
"Telegram": {
  "Enabled": true,
  "AutoEnableSession": false,
  "Config": {
    "Token": "your_telegram_bot_token"
  },
  "DefConf": {
    "allowed": true,
    "historySize": 50,
    "provider": "openai",
    "agentMaxIters": 10,
    "respond": {
      "always": true,
      "match": ".*"
    },
    "systemPromptFile": "",
    "tools": [
      { "name": "message_reaction", "enabled": true },
      { "name": "javascript", "enabled": false },
      { "name": "skills", "enabled": false },
      { "name": "subagent", "enabled": false },
      { "name": "cron", "enabled": false }
    ],
    "skills": []
  }
}
```

- **`DefConf`** (object, optional): Default configuration for sessions from this channel.
  - **`allowed`** (boolean): Default allowed status for new sessions.
  - **`historySize`** (integer): Default history size.
  - **`provider`** (string): Default LLM provider name.
  - **`agentMaxIters`** (integer): Default max agent iterations.
  - **`respond`** (object): Default respond settings.
  - **`systemPromptFile`** (string): Path to default system prompt file.
  - **`tools`** (array): Default enabled tools.
  - **`skills`** (array): Default enabled skills.

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

### Global Default Configuration (DefConf)

You can define a global default configuration that applies to all new sessions unless overridden by a channel-specific configuration.

```json
{
  "Database": "path/to/database.db",
  "LogLevel": "Info",
  "Channels": { ... },
  "Providers": [ ... ],
  "DefConf": {
    "allowed": false,
    "historySize": 50,
    "provider": "openai",
    "agentMaxIters": 3,
    "respond": {
      "always": true,
      "match": ""
    },
    "systemPromptFile": "",
    "tools": [
      { "name": "message_reaction", "enabled": true },
      { "name": "javascript", "enabled": false },
      { "name": "skills", "enabled": false },
      { "name": "subagent", "enabled": false },
      { "name": "cron", "enabled": false }
    ],
    "skills": []
  }
}
```

- **`DefConf`** (object, optional): Global default configuration for new sessions.
  - **`allowed`** (boolean): Default allowed status for new sessions.
  - **`historySize`** (integer): Default number of messages to include in LLM context.
  - **`provider`** (string): Default LLM provider name.
  - **`agentMaxIters`** (integer): Default maximum tool-use iterations per response.
  - **`respond`** (object): Default response triggering settings.
    - **`always`** (boolean): If `true`, respond to every message.
    - **`match`** (string): Regex pattern - only respond if message matches.
  - **`systemPromptFile`** (string): Path to default system prompt file.
  - **`tools`** (array): Default tool list with enabled/disabled status.
  - **`skills`** (array): Default skill list with enabled/disabled status.

### Default Configuration Priority

When a new session is created, the configuration is applied in the following order:

1. **Channel-specific DefConf** (highest priority) - defined in `Channels.Telegram.DefConf` or `Channels.Whatsapp.DefConf`
2. **Global DefConf** - defined in the root `DefConf` object
3. **Hardcoded defaults** (lowest priority) - `allowed: false`, `historySize: 50`, etc.

This allows you to set sensible defaults globally while still being able to customize defaults per channel (e.g., enable more features for Telegram sessions but fewer for WhatsApp).

## Persistent Task Storage
 
 Besides the main configuration, JGBot maintains several files to store the state of active tasks and sessions:
 
 - **`config/sessions.json`**: Stores session-specific configurations and permissions.
 - **`config/cron.json`**: Contains active cron job definitions.
 - **`config/timers.json`**: Stores pending timers and alarms.
 
 These files are managed automatically by the system when tasks are added or removed through AI tools.

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

