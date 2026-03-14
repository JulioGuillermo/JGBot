package telegramdb

import (
	"JGBot/database"
	"errors"

	"gorm.io/gorm"
)

type TelegramMessage struct {
	gorm.Model
	ChatID   uint
	Chat     *TelegramChat
	SenderID uint
	Sender   *TelegramSender

	MessageID int
	Text      string
}

func GetMessage(id uint) (*TelegramMessage, error) {
	var message TelegramMessage
	message.ID = id
	err := database.DB.First(&message).Error
	return &message, err
}

func ReceivedMessage(chat *TelegramChat, sender *TelegramSender, messageID int, text string) (*TelegramMessage, error) {
	message, err := FindMessage(chat, sender, messageID)
	if err != nil {
		return nil, err
	}
	if message != nil {
		return message, err
	}
	message = &TelegramMessage{
		Chat:   chat,
		Sender: sender,

		MessageID: messageID,
		Text:      text,
	}
	err = message.Save()
	return message, err
}

func FindMessage(chat *TelegramChat, sender *TelegramSender, messageID int) (*TelegramMessage, error) {
	var message TelegramMessage
	err := database.DB.Where("chat_id", chat.ID).Where("sender_id", sender.ID).Where("message_id", messageID).First(&message).Error
	if err == nil {
		return &message, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return nil, err
}

func (m *TelegramMessage) Save() error {
	return database.DB.Save(m).Error
}
