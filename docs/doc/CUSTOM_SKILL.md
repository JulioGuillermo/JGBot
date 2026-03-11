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
│   └── init.js           # The main entry point for the skill tool
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

#### init.js (Main Entry Point)

The `init.js` file contains the JavaScript implementation of your skill's tool. This file is executed as a module.

## Skill Metadata

### Required Fields

- `name`: Unique identifier for the skill (must match directory name)
- `description`: Brief description of the skill's purpose

### Optional Fields

- `license`: License for the skill (default: MIT)
- `metadata`: Additional metadata with:
  - `author`: Skill author name
  - `version`: Skill version

## JavaScript Environment

The JGBot runtime provides a specialized environment for skills.

> [!IMPORTANT]
> **Method Casing Rule**: Methods imported from Go into the JavaScript runtime always start with an **Uppercase** letter (e.g., `ReadFile`, `SetURL`). Standard JavaScript functions or variables follow standard JS casing conventions.

### Core Functions

- `GetArgs()`: Returns the arguments passed to the skill as an object. These arguments are passed by the AI agent.
- `print(...)`: Logs information to the bot's console/logs.

### Virtual File System (VFS)

Skills have access to two virtual file systems, isolated per session. A session is tied to a specific conversation (e.g., a Telegram chat or WhatsApp group).

- `VFPrivate`: Private storage. Only this specific skill within the current session has access to these files.
- `VFShared`: Shared storage. All skills in the current session have access to these files.

#### VFS Methods

Both `VFPrivate` and `VFShared` objects provide the following methods:

- `ReadFile(path)`: Reads a file and returns its content as a `Uint8Array`.
- `ReadStrFile(path)`: Reads a file and returns its content as a `string`.
- `WriteFile(path, content)`: Writes `Uint8Array` content to a file.
- `WriteStrFile(path, content)`: Writes `string` content to a file.
- `Exists(path)`: Checks if a file or directory exists.
- `DeleteFile(path)`: Deletes a file.
- `ReadDir(path)`: Lists the contents of a directory.
- `CreateDir(path)`: Creates a new directory.
- `DeleteDir(path)`: Deletes a directory.
- `Info(path)`: Returns metadata about a file/directory.
- `Join(...paths)`: Joins multiple path segments into one.
- `Split(path)`: Splits a path into its component segments.
- `Parent(path)`: Returns the parent directory of a path.
- `Name(path)`: Returns the base name of a path.

### HTTP Client

The `http` functionality is provided via `HttpRequest` and `HttpFormData` builders.

#### HttpRequest

Create a request using `HttpRequest()`. It supports a fluent API:

- `SetURL(url)`: Sets the target URL.
- `SetMethod(method)`: Sets the HTTP method (GET, POST, etc.).
- `SetHeader(key, ...values)`: Sets a header.
- `AddHeader(key, ...values)`: Adds values to an existing header.
- `RemoveHeader(key)`: Removes a header.
- `SetBody(bytes)`: Sets the request body from a `Uint8Array`.
- `SetBodyString(string)`: Sets the request body from a string.
- `SetBodyObj(obj)`: Serializes an object to JSON and sets it as the body.
- `SetBodyFormData(formData)`: Sets the body using a `HttpFormData` object.

**Execution methods:**
- `Fetch()`: Executes the request with the current configuration.
- `Get()`, `Post()`, `Put()`, `Delete()`, `Head()`, `Options()`, `Patch()`: Shorthands that set the method and then call `Fetch()`.

#### HttpResponse

The `Fetch()` method (and its shorthands) returns a response object with:

- `Status`: String status (e.g., "200 OK").
- `StatusCode`: Numeric status code.
- `Header`: Object containing response headers.
- `BodyBytes()`: Returns the response body as `Uint8Array`.
- `BodyString()`: Returns the response body as `string`.
- `CloseBody()`: Closes the response body (called automatically by `BodyBytes` and `BodyString`).

#### HttpFormData

Create multi-part form data using `HttpFormData()`:

- `AddField(key, value)`: Adds a text field.
- `AddFile(key, filename, bytes)`: Adds a file field.

## Creating a Custom Skill

### Basic Skill (No Tool)

If your skill only provides context or instructions to the AI without needing to execute code, just create `SKILL.md`.

### Skill with Custom Tool (JavaScript)

1. Create the skill directory: `skills/my_skill/`
2. Create `SKILL.md` to describe the tool's purpose.
3. Create `init.js` with your logic.

#### Example: Interaction with VFS and HTTP

```javascript
const args = GetArgs();

async function main() {
  const { city } = args;
  if (!city) return "Please provide a city name.";

  // Check cache in private VFS
  const cachePath = `weather_${city}.json`;
  if (VFPrivate.Exists(cachePath)) {
    return VFPrivate.ReadStrFile(cachePath);
  }

  // Fetch from API
  const response = HttpRequest()
    .SetURL(`https://api.example.com/weather?q=${city}`)
    .Get();

  if (response.StatusCode !== 200) {
    return `Error: Received status ${response.Status}`;
  }

  const data = response.BodyString();
  
  // Save to cache
  VFPrivate.WriteStrFile(cachePath, data);

  return data;
}

// The skill must export the result as default
export default main();
```

## Tool Input and Output

### Input (Arguments)
The `GetArgs()` function returns an object containing the parameters defined by the AI's understanding of the `SKILL.md` usage section. Ensure your `SKILL.md` clearly describes what inputs the tool expects.

### Output (Return Value)
- Your skill tool should return a **string**. This string will be fed back to the AI as the tool's result.
- You can export a string directly: `export default "Done!";`
- You can export a promise: `export default Promise((resolve) => resolve("Result"));`

- You can not export a function: `export default () => "Result";`
  Right now this will not work.

## Best Practices

1. **Single Responsibility**: Each skill should focus on one specific task.
2. **Clear Description**: Provide unambiguous descriptions in `SKILL.md` so the AI knows when to use the tool.
3. **Error Handling**: Always handle errors gracefully and return helpful error messages for the AI.

## Troubleshooting

- **Skill not loading**: Check `SKILL.md` YAML syntax.
- **Runtime Errors**: Check the bot output for JavaScript stack traces.
- **Undefined methods**: Verify you are using Uppercase for Go-exported methods.

