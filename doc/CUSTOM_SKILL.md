# Custom Skills Development Guide

This guide explains how to create custom skills for the JGBot project, including their structure, location, tool creation, and activation/deactivation.

## Overview

Skills are modular JavaScript extensions that extend the bot's capabilities. They can be simple text generators or complex tools with custom functionality.

## Skill Structure

### Directory Structure

Skills are located in the `skills/` directory. Each skill is a subdirectory containing:

```
skills/
├── skill_name/
│   ├── SKILL.md          # Skill metadata and content
│   └── init.js           # Optional JavaScript tool (if HasTool=true)
```

### Required Files

#### SKILL.md

Every skill must have a `SKILL.md` file with YAML frontmatter and content:

```yaml
---
name: skill_name
description: Brief description of what this skill does
license: MIT
metadata:
  author: your_name
  version: "1.0"
---

# Skill Content

Here you write the main content of your skill. This can be:
- Instructions for the AI on how to use this skill
- Examples of usage
- Context and constraints
- Any other relevant information

## Usage

The bot will use this content to understand when and how to activate your skill.
```

#### init.js (Optional)

If your skill has custom functionality, include an `init.js` file. This file contains the JavaScript implementation of your skill's tool.

## Skill Metadata

### Required Fields

- `name`: Unique identifier for the skill (must match directory name)
- `description`: Brief description of the skill's purpose

### Optional Fields

- `license`: License for the skill (default: MIT)
- `metadata`: Additional metadata with:
  - `author`: Skill author name
  - `version`: Skill version

## Creating a Custom Skill

### Basic Skill (No Tool)

1. Create skill directory:

```bash
mkdir skills/my_basic_skill
```

2. Create `SKILL.md`:

```yaml
---
name: my_basic_skill
description: A simple greeting skill
---

# Greeting Skill

This skill provides friendly greetings to users.

## Usage
When a user says hello or wants to greet someone, use this skill to generate a friendly response.
```

### Skill with Custom Tool

1. Create skill directory:

```bash
mkdir skills/my_advanced_skill
```

2. Create `SKILL.md`:

```yaml
---
name: my_advanced_skill
description: Advanced skill with custom functionality
---

# Advanced Calculator Skill

This skill performs mathematical calculations and data processing.

## Usage
When users need calculations, data analysis, or custom processing, use this skill.
```

3. Create `init.js`:

```javascript
// Get skill arguments
const args = GetArgs();

// Example: Calculator skill
function calculate(operation, a, b) {
  switch (operation) {
    case "add":
      return (parseFloat(a) + parseFloat(b)).toString();
    case "subtract":
      return (parseFloat(a) - parseFloat(b)).toString();
    case "multiply":
      return (parseFloat(a) * parseFloat(b)).toString();
    case "divide":
      return (parseFloat(a) / parseFloat(b)).toString();
    default:
      return "Unknown operation";
  }
}

// Main execution
async function main() {
  try {
    const { operation, a, b } = args;

    if (!operation || !a || !b) {
      return 'Error: Missing required arguments. Usage: {operation: "add|subtract|multiply|divide", a: number, b: number}';
    }

    const result = calculate(operation, a, b);
    return `Calculation result: ${result}`;
  } catch (error) {
    return `Error: ${error.message}`;
  }
}

// Execute the skill
// Export the final output of the skill as default
export default main();
```

You can use import and export statements, but all the code of the skill must be inside the skill directory. No external files are allowed.

```
skills/
├── skill_name/
│   ├── operations/
│   │   └── operations.js # Operations
│   ├── SKILL.md          # Skill metadata and content
│   ├── init.js           # Optional JavaScript tool (if HasTool=true)
```

```javascript
// operations/operations.js

// Example: Calculator skill
export function calculate(operation, a, b) {
  switch (operation) {
    case "add":
      return (parseFloat(a) + parseFloat(b)).toString();
    case "subtract":
      return (parseFloat(a) - parseFloat(b)).toString();
    case "multiply":
      return (parseFloat(a) * parseFloat(b)).toString();
    case "divide":
      return (parseFloat(a) / parseFloat(b)).toString();
    default:
      return "Unknown operation";
  }
}
```

```javascript
// init.js
import { calculate } from "/operations/operations.js";

const args = GetArgs();

async function main() {
  try {
    const { operation, a, b } = args;

    if (!operation || !a || !b) {
      return 'Error: Missing required arguments. Usage: {operation: "add|subtract|multiply|divide", a: number, b: number}';
    }

    const result = calculate(operation, a, b);
    return `Calculation result: ${result}`;
  } catch (error) {
    return `Error: ${error.message}`;
  }
}

export default main();
```

## JavaScript Environment

### Available Functions

- `GetArgs()`: Returns the arguments passed to the skill as an object
- `VFPrivate`: Virtual file system for private data, just this skill in the current session has access to
- `VFShared`: Virtual file system for shared data, all skill in this session have access to
- `http`: HTTP client for making web requests

### Virtual File System

Skills have access to two virtual file systems:

- `VFPrivate`: Session-specific private files
- `VFShared`: Shared files across all sessions

Virtual Files are isolate for each session. A session is a conversation from an origin, e.g. a WhatsApp or Telegram group or private chat. Each session has its own private (Just the owner skill) and shared (All skills) files.

Example usage:

```javascript
// Read from private files
const privateContent = await VFPrivate.ReadFile("data.json");

// Write to private files
await VFPrivate.WriteFile("result.txt", "New data");

// Read from shared files
const sharedData = await VFShared.ReadFile("config.json");
```

### HTTP Client

Access to HTTP functionality is available through the `http` object:

```javascript
const formData = HttpFormData()
  .AddField("name", "test")
  .AddFile("file", "test.txt", new Uint8Array([1, 2, 3, 4, 5]));

const response = HttpRequest()
  .SetURL("https://httpbin.org/get")
  .SetBodyFormData(formData)
  .Get();

console.log(response.BodyString());
```

## Skill Activation and Deactivation

### Automatic Loading

Skills are automatically loaded on startup from the `skills/` directory. The system:

1. Scans all directories in `skills/`
2. Loads valid skills (those with proper `SKILL.md` files)
3. Checks for `init.js` files to determine if they have tools
4. Makes them available to the bot

### Activation

Skills are activated when the AI determines they are relevant to the current conversation. The bot uses the content in `SKILL.md` to understand when to use each skill.

### Deactivation

Skills remain active for the duration of a conversation session. There is no explicit deactivation mechanism, but the AI will stop using a skill when it's no longer relevant.

### Managing Skills

#### To Disable a Skill

1. Rename the skill directory (prefix with underscore):

```bash
mv skills/my_skill skills/_my_skill
```

2. Or remove the skill entirely:

```bash
rm -rf skills/my_skill
```

#### To Enable a Skill

1. Ensure the skill directory name doesn't start with underscore
2. Make sure `SKILL.md` is properly formatted
3. Restart the bot to reload skills

## Best Practices

### Skill Design

1. **Single Responsibility**: Each skill should focus on one specific task
2. **Clear Description**: Provide unambiguous descriptions in `SKILL.md`
3. **Error Handling**: Always handle errors gracefully in JavaScript tools
4. **Input Validation**: Validate all inputs in your JavaScript code

### Performance Considerations

1. Keep JavaScript tools lightweight
2. Use async/await for non-blocking operations
3. Cache data when appropriate to avoid repeated requests
4. Implement timeouts for external requests

### Security

1. Validate all user inputs
2. Sanitize any data that might be displayed
3. Be careful with file system operations
4. Implement proper error handling without exposing sensitive information

## Example Skills

### Weather Skill

```yaml
---
name: weather_skill
description: Get current weather information for a location
---

# Weather Skill

This skill retrieves current weather information for specified locations.

## Usage
When users ask about weather conditions, use this skill to get current weather data.
```

```javascript
// init.js
const args = GetArgs();

async function main() {
  try {
    const { location } = args;

    if (!location) {
      return 'Error: Location is required. Usage: {location: "city name"}';
    }

    const response = await http.get(
      `https://api.weather.com/v1/weather?location=${encodeURIComponent(location)}`,
    );
    const weather = JSON.parse(response.body);

    return `Current weather in ${location}: ${weather.temperature}°C, ${weather.description}`;
  } catch (error) {
    return `Error fetching weather: ${error.message}`;
  }
}

await main();
```

### File Processing Skill

```yaml
---
name: file_processor
description: Process and analyze uploaded files
---

# File Processing Skill

This skill processes uploaded files, extracting information and performing analysis.

## Usage
When users upload files and need analysis or processing, use this skill.
```

```javascript
// init.js
const args = GetArgs();

async function main() {
  try {
    const { filename, action } = args;

    if (!filename || !action) {
      return 'Error: Filename and action are required. Usage: {filename: "file.txt", action: "read|analyze|process"}';
    }

    const filepath = path.join(VFPrivate.path, filename);

    switch (action) {
      case "read":
        const content = await VFPrivate.readFile(filename);
        return `File content: ${content}`;

      case "analyze":
        const data = await VFPrivate.readFile(filename);
        const wordCount = data.split(" ").length;
        const lineCount = data.split("\n").length;
        return `Analysis: ${wordCount} words, ${lineCount} lines`;

      default:
        return "Error: Unknown action. Use: read, analyze, or process";
    }
  } catch (error) {
    return `Error processing file: ${error.message}`;
  }
}

await main();
```

## Troubleshooting

### Common Issues

1. **Skill not loading**: Check `SKILL.md` format and ensure required fields are present
2. **JavaScript errors**: Check console output for syntax or runtime errors
3. **Permission issues**: Ensure proper file permissions on skill directories
4. **Missing dependencies**: Skills run in a restricted environment; external dependencies are limited

### Debugging

1. Check the bot logs for skill loading errors
2. Test JavaScript tools independently if possible
3. Use console.log statements for debugging (output appears in bot logs)
4. Verify file paths and permissions for virtual file system operations

## Advanced Features

### Context Access

Skills can access conversation context through the `args` parameter, which includes:

- User messages
- Conversation history
- Session data
- Configuration settings

### Multi-step Operations

Skills can perform complex, multi-step operations:

```javascript
async function main() {
  // Step 1: Get data
  const data = await fetchExternalData();

  // Step 2: Process data
  const processed = processData(data);

  // Step 3: Save results
  await VFPrivate.writeFile("result.json", JSON.stringify(processed));

  // Step 4: Return summary
  return `Processed ${processed.length} items`;
}
```

### Integration with Other Skills

Skills can be designed to work together by following consistent patterns and sharing data through the virtual file system.
