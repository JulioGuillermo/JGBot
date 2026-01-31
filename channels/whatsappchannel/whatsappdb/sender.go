package whatsappdb

import (
	"JGBot/database"
	"errors"

	"go.mau.fi/whatsmeow/types"
	"gorm.io/gorm"
)

type WhatsAppSender struct {
	JID

	Name string
}

func GetSender(id uint) (*WhatsAppSender, error) {
	var sender WhatsAppSender
	sender.ID = id
	err := database.DB.First(&sender).Error
	return &sender, err
}

func ReceivedSender(senderJID *types.JID, name string) (*WhatsAppSender, error) {
	sender, err := FindSender(senderJID)
	if err != nil {
		return nil, err
	}
	if sender != nil {
		err = sender.UpdateIfChange(name)
		return sender, err
	}
	sender = &WhatsAppSender{
		Name: name,
	}
	sender.FromJID(senderJID)
	err = sender.Save()
	return sender, err
}

func FindSender(senderJID *types.JID) (*WhatsAppSender, error) {
	var sender WhatsAppSender
	err := database.DB.
		Where("user", senderJID.User).
		Where("raw_agent", senderJID.RawAgent).
		Where("device", senderJID.Device).
		Where("integrator", senderJID.Integrator).
		Where("server", senderJID.Server).
		First(&sender).Error
	if err == nil {
		return &sender, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return nil, err
}

func (s *WhatsAppSender) UpdateIfChange(name string) error {
	if !s.Change(name) {
		return nil
	}
	s.Name = name
	return s.Save()
}

func (s *WhatsAppSender) Change(name string) bool {
	return s.Name != name
}

func (s *WhatsAppSender) Save() error {
	return database.DB.Save(s).Error
}
