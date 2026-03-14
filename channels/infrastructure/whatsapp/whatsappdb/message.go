package whatsappdb

import (
	"JGBot/database"
	"errors"

	"gorm.io/gorm"
)

type WhatsAppMessage struct {
	gorm.Model
	ChatID   uint
	Chat     *WhatsAppChat
	SenderID uint
	Sender   *WhatsAppSender

	MessageID string
	Text      string
}

func GetMessage(id uint) (*WhatsAppMessage, error) {
	var message WhatsAppMessage
	message.ID = id
	err := database.DB.Preload("Sender").First(&message).Error
	return &message, err
}

func ReceivedMessage(chat *WhatsAppChat, sender *WhatsAppSender, messageID string, text string) (*WhatsAppMessage, error) {
	message, err := FindMessage(chat, sender, messageID)
	if err != nil {
		return nil, err
	}
	if message != nil {
		return message, err
	}
	message = &WhatsAppMessage{
		Chat:   chat,
		Sender: sender,

		MessageID: messageID,
		Text:      text,
	}
	err = message.Save()
	return message, err
}

func FindMessage(chat *WhatsAppChat, sender *WhatsAppSender, messageID string) (*WhatsAppMessage, error) {
	var message WhatsAppMessage
	err := database.DB.Where("chat_id", chat.ID).Where("sender_id", sender.ID).Where("message_id", messageID).First(&message).Error
	if err == nil {
		return &message, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return nil, err
}

func (m *WhatsAppMessage) Save() error {
	return database.DB.Save(m).Error
}
