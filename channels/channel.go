package channels

type Status int

const (
	Normal = Status(iota)
	Writing
)

type OnMessageHandler func(channel string, origin string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string)

type Channel interface {
	GetName() string
	OnMessage(handler OnMessageHandler)
	Status(chatID uint, status Status) error
	SendMessage(chatID uint, message string) error
	ReactMessage(chatID uint, messageID uint, reaction string) error
	AutoEnableSession() bool
	Close()
}
