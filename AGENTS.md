# AGENTS.md - JGBot Development Guide

This document provides guidelines for agentic coding agents working on the JGBot project, specifically tailored for the **Hexagonal Architecture** migration.

## Build/Test Commands

### Go Commands
```bash
# Build the main application
go build -o jgbot main.go

# Run tests (all tests)
go test ./...

# Run specific package tests
go test ./session/...
go test ./channels/...
go test ./agent/...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run a specific test function
go test -run TestJS ./test
```

### Testing Approach
- Tests are located within their respective packages or in the `test/` directory.
- Follow the hexagonal pattern:
  - `domain/`: Unit tests with minimal mocking.
  - `application/`: Integration tests using mocked domain interfaces.
  - `infrastructure/`: Adapter tests verifying external system interaction.
- Use standard Go testing framework (`testing` package).
- Test files follow Go naming convention: `*_test.go`.
- Use table-driven testing for multiple test cases.

## Hexagonal Architecture Guidelines

### Core Principles
- **SOLID**: Follow all five SOLID principles.
- **Dependency Inversion**: High-level modules (application) should not depend on low-level modules (infrastructure). Both should depend on abstractions (domain interfaces).
- **Domain First**: The domain layer must have **zero external dependencies** (except for standard library or basic utilities). It should not import `gorm`, `llms`, or platform-specific libraries.
- **Ports & Adapters**: Use interfaces in the domain to define "Ports", and implement them in the infrastructure as "Adapters".

### Package Structure
- `domain/`: Entities, value objects, and repository/service interfaces.
- `application/`: Application services that orchestrate domain logic.
- `infrastructure/`: Implementations of domain interfaces (GORM, JSON files, API clients).

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
    "github.com/tmc/langchaingo/llms"
    "gorm.io/gorm"
    "slog"
)

// Project imports last, with full module path
import (
    "JGBot/agent/domain"
    "JGBot/agent/handler"
    "JGBot/conf"
    "JGBot/database"
    "JGBot/log"
    "JGBot/session/domain"
)
```

### Naming Conventions
```go
// Variables and functions: camelCase
func getUserSession(id uint) (*Session, error)

// Types and interfaces: PascalCase
type SessionStore interface {
    GetConfig(origin string) *SessionConfig
}

// Constants: UPPER_SNAKE_CASE
const MaxHistorySize = 100
```

### Error Handling
```go
// Always return errors with context
func doSomething() error {
    err := repository.Save(entity)
    if err != nil {
        return fmt.Errorf("failed to save entity: %w", err)
    }
    return nil
}
```

### Logging
```go
// Use structured logging with slog
log.Info("Message processed", "channel", channel, "chat_id", chatID)
```

### Database Operations
- **NEVER** use `database.DB` directly in domain logic.
- Always implement a repository interface in `infrastructure/persistence` and inject it into the application service.

### Testing Patterns
```go
// Use domain types in tests
func TestAgentRun(t *testing.T) {
    tests := []struct {
        name    string
        history []*agentdomain.SessionMessage
        message *agentdomain.SessionMessage
        want    string
        wantErr bool
    }{
        {
            name:    "empty history",
            history: []*agentdomain.SessionMessage{},
            message: &agentdomain.SessionMessage{Message: "hello"},
            want:    "Hello",
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            agent := NewSubAgent(tt.name, provider, tt.history, tt.message)
            got, err := agent.Run()
            // ... checks
        })
    }
}
```

## Skill Development Guidelines
- Skills should remain in the `skills/` directory.
- Use JavaScript for skill logic.
- Ensure metadata in `SKILL.md` is accurate.

## Performance Considerations
- Use the `database` package's connection pooling.
- Implement caching in infrastructure adapters when necessary.
- Pass `context.Context` through all calls to support timeouts and cancellation.

## Security Guidelines
- Validate all inputs in the application layer.
- Ensure infrastructure adapters handle sensitive data (like API keys from config) securely.
- Never hardcode credentials; use `conf` or environment variables.
