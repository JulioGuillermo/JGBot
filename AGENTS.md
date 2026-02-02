# AGENTS.md - JGBot Development Guide

This document provides guidelines for agentic coding agents working on the JGBot project.

## Build/Test Commands

### Go Commands
```bash
# Build the main application
go build -o jgbot main.go

# Run tests (all tests)
go test ./...

# Run specific test file
go test ./test

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run a specific test function
go test -run TestJS ./test

# Build and run
go run main.go
```

### Testing Approach
- Tests are located in the `test/` directory
- Plugin tests use `plugins/test/` directory
- Use standard Go testing framework (`testing` package)
- Test files follow Go naming convention: `*_test.go`
- Use table-driven testing for multiple test cases

## Code Style Guidelines

### Imports
```go
// Standard library imports first, sorted alphabetically
import (
    "context"
    "fmt"
    "os"
    "slices"
    "strings"
    "time"

    // Third-party imports next, sorted alphabetically
    "github.com/fastschema/qjs"
    "github.com/tmc/langchaingo/agents"
    "github.com/tmc/langchaingo/llms"
    "gorm.io/gorm"
    "slog"
)

// Project imports last, with full module path
import (
    "JGBot/agent/handler"
    "JGBot/agent/tools"
    "JGBot/conf"
    "JGBot/database"
    "JGBot/log"
    "JGBot/session/sessiondb"
)
```

### Package Structure
- Follow the existing directory structure
- Each package has a single responsibility
- Use descriptive package names (lowercase)
- Keep packages focused and cohesive

### Naming Conventions
```go
// Variables and functions: camelCase
func getUserSession(id uint) (*Session, error)
var messageQueue chan Message

// Types and interfaces: PascalCase
type Agent struct {
    Ctx      context.Context
    Name     string
    Handler  *handler.AgentHandler
    Provider llms.Model
}

// Constants: UPPER_SNAKE_CASE
const (
    MaxRetries    = 3
    DefaultTimeout = 30 * time.Second
)

// Private fields: camelCase with underscore prefix
type agent struct {
    ctx      context.Context
    name     string
    provider llms.Model
}
```

### Error Handling
```go
// Always return errors with context
func doSomething() error {
    err := database.InitConnection()
    if err != nil {
        return fmt.Errorf("failed to initialize database: %w", err)
    }
    return nil
}

// Check errors immediately
func processMessage(msg string) error {
    if msg == "" {
        return errors.New("message cannot be empty")
    }
    // ... rest of function
}

// Use custom error types for specific errors
type ConfigError struct {
    Field string
    Err   error
}

func (e *ConfigError) Error() string {
    return fmt.Sprintf("config error on field %q: %v", e.Field, e.Err)
}
```

### Logging
```go
// Use structured logging with slog
log.Info("System starting", "version", "1.0.0")
log.Error("Database connection failed", "error", err, "attempt", attempt)
log.Debug("Processing message", "message_id", msgID, "channel", channel)

// Use appropriate log levels
log.Debug("Detailed debug information")
log.Info("General information")
log.Warn("Warning conditions")
log.Error("Error conditions")
```

### Context Usage
```go
// Always pass context through the call chain
func processRequest(ctx context.Context, req Request) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // Process request
    }
}
```

### Database Operations
```go
// Use GORM for database operations
func saveMessage(msg *SessionMessage) error {
    return database.DB.Create(msg).Error
}

// Use transactions for atomic operations
func transferFunds(ctx context.Context, from, to User, amount float64) error {
    return database.DB.Transaction(func(tx *gorm.DB) error {
        // Update sender balance
        if err := tx.Model(&from).Update("balance", from.Balance-amount).Error; err != nil {
            return err
        }
        // Update receiver balance
        if err := tx.Model(&to).Update("balance", to.Balance+amount).Error; err != nil {
            return err
        }
        return nil
    })
}
```

### Type Safety
```go
// Use interfaces for abstractions
type MessageProcessor interface {
    Process(msg string) error
}

// Use specific types instead of generic ones
type UserID uint
type ChannelID string
type MessageID uint

// Use value types for small, fixed sets
type LogLevel string

const (
    LogLevelDebug LogLevel = "debug"
    LogLevelInfo  LogLevel = "info"
    LogLevelWarn  LogLevel = "warn"
    LogLevelError LogLevel = "error"
)
```

### Function Design
```go
// Keep functions focused and small
func loadSkills() ([]*Skill, error) {
    skillDirs, err := readSkillDir()
    if err != nil {
        return nil, err
    }
    
    skills := make([]*Skill, 0, len(skillDirs))
    for _, dir := range skillDirs {
        skill := LoadSkill(dir)
        if skill != nil {
            skills = append(skills, skill)
        }
    }
    return skills, nil
}

// Use named return values for documentation
func calculateTotal(items []Item) (total float64, err error) {
    for _, item := range items {
        total += item.Price * float64(item.Quantity)
    }
    return total, nil
}
```

### Resource Management
```go
// Use defer for cleanup
func processFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    // Process file
    return nil
}

// Use context for timeouts and cancellation
func fetchWithTimeout(ctx context.Context, url string) ([]byte, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    // ... rest of function
}
```

### Configuration Management
```go
// Use the conf package for configuration
func initConfig() error {
    if err := conf.InitConfig(); err != nil {
        return fmt.Errorf("failed to initialize config: %w", err)
    }
    return nil
}

// Access configuration through global variable
func getDatabasePath() string {
    return conf.Conf.Database
}
```

### Testing Patterns
```go
// Use table-driven tests
func TestAgentRun(t *testing.T) {
    tests := []struct {
        name    string
        history []*sessiondb.SessionMessage
        message *sessiondb.SessionMessage
        want    string
        wantErr bool
    }{
        {
            name:    "empty history",
            history: []*sessiondb.SessionMessage{},
            message: &sessiondb.SessionMessage{Message: "hello"},
            want:    "Hello",
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            agent := NewTestAgent()
            got, err := agent.Run(tt.history, tt.message)
            if (err != nil) != tt.wantErr {
                t.Errorf("Agent.Run() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Agent.Run() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Skill Development Guidelines

### Skill Structure
- Skills are located in the `skills/` directory
- Each skill has a `SKILL.md` metadata file
- Skills implement the `Skill` interface
- Use JavaScript for skill logic

### Skill Metadata
```yaml
---
name: skill_name
description: Brief description of the skill
license: MIT
metadata:
  author: author_name
  version: "1.0"
---
```

### Error Handling in Skills
```go
func ExecSkillTool(name string, args SkillArgs) (string, error) {
    output, err := runners.RunModule(
        SkillToolFile,
        path.Join(SkillDir, name),
        // ... exec options
    )
    if err != nil {
        return "", fmt.Errorf("skill execution failed: %w", err)
    }
    return output.Result, nil
}
```

## Channel Development

### Channel Interface
- Implement the `Channel` interface
- Handle incoming and outgoing messages
- Support reactions and status updates
- Use database for persistence

### Message Handling
```go
type Channel interface {
    SendMessage(channel string, chatID uint, message string) error
    ReactMessage(channel string, chatID uint, messageID uint, reaction string) error
    OnMessage(handler func(channel string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string))
}
```

## Agent Development

### Agent Configuration
- Support multiple LLM providers (OpenAI, Anthropic, Google, Ollama, Mistral)
- Use configurable temperature and max tokens
- Implement retry logic for API calls

### Tool Integration
- Tools implement the `tools.Tool` interface
- Use proper argument validation
- Handle errors gracefully
- Provide meaningful descriptions

## Performance Considerations

- Use connection pooling for database operations
- Implement rate limiting for API calls
- Cache frequently accessed data
- Use appropriate data structures for performance
- Profile and optimize critical paths

## Security Guidelines

- Validate all inputs
- Sanitize user-generated content
- Use secure configuration storage
- Implement proper error handling without exposing sensitive information
- Follow principle of least privilege

## Code Quality

- Use Go fmt for code formatting
- Implement comprehensive testing
- Document public APIs
- Use linters and static analysis
- Keep dependencies up to date
- Follow semantic versioning