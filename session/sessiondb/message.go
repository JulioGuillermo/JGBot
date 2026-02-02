package sessiondb

import (
	"JGBot/database"
	"fmt"
	"slices"
	"strings"

	"gorm.io/gorm"
)

type SessionMessage struct {
	gorm.Model

	Channel string

	ChatID   uint
	ChatName string

	SenderID   uint
	SenderName string

	MessageID uint
	Message   string

	Role  string
	Extra string
}

func Migrate() error {
	return database.DB.AutoMigrate(&SessionMessage{})
}

func NewSessionMessage(channel string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string, role string, extra string) (*SessionMessage, error) {
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
	err := msg.Save()
	return msg, err
}

func GetHistory(channel string, chatID uint, limit int) ([]*SessionMessage, error) {
	history := []*SessionMessage{}
	err := database.DB.Model(&SessionMessage{}).Where("channel = ? AND chat_id = ?", channel, chatID).Order("id DESC").Limit(limit).Find(&history).Error
	slices.Reverse(history)
	return history, err
}

func ClearHistory(channel string, chatID uint) error {
	return database.DB.Where("channel = ? AND chat_id = ?", channel, chatID).Delete(&SessionMessage{}).Error
}

func (m *SessionMessage) Save() error {
	return database.DB.Save(m).Error
}

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
