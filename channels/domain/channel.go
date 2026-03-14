package domain

type Channel interface {
	GetName() string
	OnMessage(handler MessageHandler)
	SendStatus(chatID uint, status Status) error
	SendMessage(chatID uint, message string) error
	SendMessageReaction(chatID uint, messageID uint, reaction string) error
	AutoEnableSession() bool
	Close()
}
