package telegramdb

import (
	"JGBot/database"
	"errors"

	"gorm.io/gorm"
)

type TelegramChat struct {
	gorm.Model
	ChatID   int64 `gorm:"unique"`
	ChatName string
}

func GetChat(id uint) (*TelegramChat, error) {
	var chat TelegramChat
	chat.ID = id
	err := database.DB.First(&chat).Error
	return &chat, err
}

func ReceivedChat(chatID int64, chatName string) (*TelegramChat, error) {
	chat, err := FindChat(chatID)
	if err != nil {
		return nil, err
	}
	if chat != nil {
		err = chat.UpdateIfChange(chatName)
		return chat, err
	}
	chat = &TelegramChat{
		ChatID:   chatID,
		ChatName: chatName,
	}
	err = chat.Save()
	return chat, err
}

func FindChat(chatID int64) (*TelegramChat, error) {
	var chat TelegramChat
	err := database.DB.Where("chat_id", chatID).First(&chat).Error
	if err == nil {
		return &chat, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, err
}

func (c *TelegramChat) UpdateIfChange(chatName string) error {
	if !c.Change(chatName) {
		return nil
	}
	c.ChatName = chatName
	return c.Save()
}

func (c *TelegramChat) Change(chatName string) bool {
	return c.ChatName != chatName
}

func (c *TelegramChat) Save() error {
	return database.DB.Save(c).Error
}
