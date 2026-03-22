package db

import (
	"fmt"
	"strings"

	sessiondomain "JGBot/session/domain"
	"gorm.io/gorm"
)

// SessionMessage is the database model for messages
type SessionMessage struct {
	gorm.Model

	Channel  string
	ChatID   uint
	ChatName string

	SenderID   uint
	SenderName string

	MessageID uint
	Message   string

	Role  string
	Extra string
}

// Migrate runs the database migration for SessionMessage
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&SessionMessage{})
}

// Save persists the message to the database
func (m *SessionMessage) Save(db *gorm.DB) error {
	return db.Save(m).Error
}

// String returns a human-readable representation of the message
func (m *SessionMessage) String() string {
	name := m.SenderName
	role := strings.ToUpper(m.Role)
	if role == "TOOL" || role == "SYSTEM" || role == "ASSISTANT" {
		name = role
	}
	return fmt.Sprintf(
		"SOURCE_CHANNEL: %s\nUSER_NAME: %s\nCHAT_NAME: %s\nMESSAGE_ID: %d\n\nCONTENT: %s",
		m.Channel, name, m.ChatName, m.MessageID, m.Message,
	)
}

// SaveMessage creates and saves a new message
func SaveMessage(db *gorm.DB, channel string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string, role string, extra string) (*SessionMessage, error) {
	msg := &SessionMessage{
		Channel:    channel,
		ChatID:     chatID,
		ChatName:   chatName,
		SenderID:   senderID,
		SenderName: senderName,
		MessageID:  messageID,
		Message:    message,
		Role:       role,
		Extra:      extra,
	}
	err := msg.Save(db)
	return msg, err
}

// GetHistory retrieves message history for a channel and chat
func GetHistory(db *gorm.DB, channel string, chatID uint, limit int) ([]*SessionMessage, error) {
	history := []*SessionMessage{}
	err := db.Model(&SessionMessage{}).Where("channel = ? AND chat_id = ?", channel, chatID).Order("id DESC").Limit(limit).Find(&history).Error

	// Reverse to get chronological order
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	return history, err
}

// ClearHistory deletes all messages for a channel and chat
func ClearHistory(db *gorm.DB, channel string, chatID uint) error {
	return db.Where("channel = ? AND chat_id = ?", channel, chatID).Delete(&SessionMessage{}).Error
}

// ToDomain converts SessionMessage to domain.Message
func (m *SessionMessage) ToDomain() *sessiondomain.Message {
	if m == nil {
		return nil
	}
	return &sessiondomain.Message{
		ID:         m.ID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		Channel:    m.Channel,
		ChatID:     m.ChatID,
		ChatName:   m.ChatName,
		SenderID:   m.SenderID,
		SenderName: m.SenderName,
		MessageID:  m.MessageID,
		Message:    m.Message,
		Role:       m.Role,
		Extra:      m.Extra,
	}
}
