---
name: javascript
description: Technical reference of the JavaScript runtime environment for agents.
metadata:
  author: JulioGuillermo
  version: "1.0"
---

# JavaScript Runtime Specification for Agents

This document specifies the available objects and functions in the sandboxed JavaScript (ES2023) environment used for script execution. This skill does not have its own tool; it provides the instructions for the global JavaScript tool.

## Core Runtime Rules

### Go-to-JS Bridge Naming Convention

> [!IMPORTANT]
> Methods and properties exported from the Go environment into JavaScript **must** be called using **Uppercase** initials (e.g., `ReadFile`, `SetURL`). JavaScript-native or user-defined functions follow standard camelCase.

## Global Functions and Objects

All methods are synchronous and can be called directly.

### Logging and Output

- `print(...args)`: Logs arguments to the session output (Use this instead of console`).
- `console`: Standard logging object with methods (Not recommended to use, use `print` instead):
  - `log(...args)`, `info(...args)`, `warn(...args)`, `error(...args)`, `debug(...args)`.

### Virtual File System: `VirtualFiles`

The `VirtualFiles` object provides access to the private session storage.

- **File Operations**:
  - `ReadFile(path)`: Returns `Uint8Array`.
  - `ReadStrFile(path)`: Returns `string`.
  - `WriteFile(path, Uint8Array)`: Writes binary data.
  - `WriteStrFile(path, string)`: Writes string data.
  - `DeleteFile(path)`: Removes a file.
  - `Exists(path)`: Returns `boolean`.
  - `Info(path)`: Returns metadata {Name, Size, IsDir}.
- **Directory Operations**:
  - `ReadDir(path)`: Returns array of file info objects {Name, Size, IsDir}.
  - `CreateDir(path)`: Creates a directory.
  - `DeleteDir(path)`: Removes a directory.
- **Path Utilities**:
  - `Join(...segments)`: Joins path segments.
  - `Split(path)`: Splits path into segments.
  - `Parent(path)`: Returns parent directory.
  - `Name(path)`: Returns file name.

### HTTP Client: `HttpRequest`

Network requests are constructed using a fluent builder.

- **HttpRequest Builder Methods**:
  - `SetURL(url)`: Sets the URL.
  - `SetMethod(method)`: Sets the HTTP method.
  - `SetHeader(key, ...values)`: Sets a header.
  - `AddHeader(key, ...values)`: Adds a header.
  - `RemoveHeader(key)`: Removes a header.
  - `SetBody(Uint8Array)`: Sets the request body.
  - `SetBodyString(string)`: Sets the request body as a string.
  - `SetBodyObj(obj)`: Sets the request body as an object.
  - `SetBodyFormData(HttpFormData)`: Sets the request body as form data.
- **Execution Methods**:
  - `Fetch()`: Executes and returns `HttpResponse`.
  - Shorthands: `Get()`, `Post()`, `Put()`, `Delete()`, `Patch()`, `Head()`, `Options()`.
- **HttpResponse Object**:
  - `Status` (string) - The status of the response.
  - `StatusCode` (int) - The status code of the response.
  - `Header` (object) - The headers of the response.
  - `BodyBytes()` (Uint8Array) - The body of the response as bytes.
  - `BodyString()` (string) - The body of the response as a string.
  - `CloseBody()` - Closes stream (automatic on body read).

### Form Data: `HttpFormData`

- `AddField(key, value)`
- `AddFile(key, filename, Uint8Array)`
