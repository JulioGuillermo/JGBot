package telegramdb

import (
	"JGBot/database"
	"errors"

	"gorm.io/gorm"
)

type TelegramSender struct {
	gorm.Model
	SenderID  int64 `gorm:"unique"`
	FirstName string
	LastName  string
	Username  string
}

func GetSender(id uint) (*TelegramSender, error) {
	var sender TelegramSender
	sender.ID = id
	err := database.DB.First(&sender).Error
	return &sender, err
}

func ReceivedSender(senderID int64, firstName, lastName, username string) (*TelegramSender, error) {
	sender, err := FindSender(senderID)
	if err != nil {
		return nil, err
	}
	if sender != nil {
		err = sender.UpdateIfChange(firstName, lastName, username)
		return sender, err
	}
	sender = &TelegramSender{
		SenderID:  senderID,
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
	}
	err = sender.Save()
	return sender, err
}

func FindSender(senderID int64) (*TelegramSender, error) {
	var sender TelegramSender
	err := database.DB.Where("sender_id", senderID).First(&sender).Error
	if err == nil {
		return &sender, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return nil, err
}

func (s *TelegramSender) UpdateIfChange(firstName, lastName, username string) error {
	if !s.Change(firstName, lastName, username) {
		return nil
	}
	s.FirstName = firstName
	s.LastName = lastName
	s.Username = username
	return s.Save()
}

func (s *TelegramSender) Change(firstName, lastName, username string) bool {
	return s.FirstName != firstName || s.LastName != lastName || s.Username != username
}

func (s *TelegramSender) Save() error {
	return database.DB.Save(s).Error
}

func (s *TelegramSender) String() string {
	if s.FirstName != "" && s.LastName != "" {
		return s.FirstName + " " + s.LastName
	}
	if s.FirstName != "" {
		return s.FirstName
	}
	if s.LastName != "" {
		return s.LastName
	}
	return s.Username
}
