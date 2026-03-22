package sessiondomain

import (
	"fmt"
	"strings"
	"time"
)

type Message struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

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

func (m *Message) String() string {
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

func NewMessage(channel string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string, role string, extra string) *Message {
	return &Message{
		Channel:    channel,
		ChatID:     chatID,
		ChatName:   chatName,
		SenderID:   senderID,
		SenderName: senderName,
		MessageID:  messageID,
		Message:    message,
		Role:       role,
		Extra:      extra,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
