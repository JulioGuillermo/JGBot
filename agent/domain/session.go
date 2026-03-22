package agentdomain

import (
	"time"
)

// SessionMessage represents a chat message in the session
type SessionMessage struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Channel    string
	ChatID     uint
	ChatName   string
	SenderID   uint
	SenderName string
	MessageID  uint
	Message    string
	Role       string
	Extra      string
}

// String returns a human-readable representation
func (m *SessionMessage) String() string {
	return m.Message
}

// SessionConfig represents session configuration
type SessionConfig struct {
	Name   string
	ID     string
	Origin string
	Admin  string

	Allowed bool
	Respond struct {
		Always bool
		Match  string
	}

	HistorySize      int
	Provider         string
	SystemPromptFile string
	AgentMaxIters    int
	Tools            []ToolConfig
	Skills           []SkillConfig
}

// ToolConfig represents a tool configuration
type ToolConfig struct {
	Name    string
	Enabled bool
}

// SkillConfig represents a skill configuration
type SkillConfig struct {
	Name        string
	Enabled     bool
	Description string
}

// SessionStore provides access to session data
type SessionStore interface {
	// Config operations
	GetConfig(origin string) *SessionConfig
	GetConfigs() []*SessionConfig
	CreateConfig(chatName, sessionID, origin, channel string) *SessionConfig
	CreateUnconfig(chatName, sessionID, origin, channel string) *SessionConfig

	// Message operations
	GetHistory(channel string, chatID uint, limit int) ([]*SessionMessage, error)
	SaveMessage(channel string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string, role string, extra string) (*SessionMessage, error)
	ClearHistory(channel string, chatID uint) error
}
