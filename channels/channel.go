package channels

type OnMessageHandler func(channel string, origin string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string)

type Channel interface {
	GetName() string
	OnMessage(handler OnMessageHandler)
	SendMessage(chatID uint, message string) error
	ReactMessage(chatID uint, messageID uint, reaction string) error
	Close()
}
