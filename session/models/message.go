package models

import "gorm.io/gorm"

type SessionMessage struct {
	gorm.Model

	Channel string

	ChatID   uint
	ChatName string

	SenderID   uint
	SenderName string

	MessageID uint
	Message   string

	FromAI bool
}
