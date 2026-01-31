package whatsappdb

import (
	"JGBot/database"
	"errors"

	"go.mau.fi/whatsmeow/types"
	"gorm.io/gorm"
)

type WhatsAppChat struct {
	JID

	Name string
}

func GetChat(id uint) (*WhatsAppChat, error) {
	var chat WhatsAppChat
	chat.ID = id
	err := database.DB.First(&chat).Error
	return &chat, err
}

func ReceivedChat(chatJID *types.JID, name string) (*WhatsAppChat, error) {
	chat, err := FindChat(chatJID)
	if err != nil {
		return nil, err
	}
	if chat != nil {
		err = chat.UpdateIfChange(name)
		return chat, err
	}
	chat = &WhatsAppChat{
		Name: name,
	}
	chat.FromJID(chatJID)
	err = chat.Save()
	return chat, err
}

func FindChat(chatJID *types.JID) (*WhatsAppChat, error) {
	var chat WhatsAppChat
	err := database.DB.
		Where("user", chatJID.User).
		Where("raw_agent", chatJID.RawAgent).
		Where("device", chatJID.Device).
		Where("integrator", chatJID.Integrator).
		Where("server", chatJID.Server).
		First(&chat).Error
	if err == nil {
		return &chat, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, err
}

func (c *WhatsAppChat) UpdateIfChange(name string) error {
	if !c.Change(name) {
		return nil
	}
	c.Name = name
	return c.Save()
}

func (c *WhatsAppChat) Change(name string) bool {
	return c.Name != name
}

func (c *WhatsAppChat) Save() error {
	return database.DB.Save(c).Error
}
