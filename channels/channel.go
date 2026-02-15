package channels

import "fmt"

type Status int

const (
	Normal = Status(iota)
	Writing
)

type OnMessageHandler func(channel string, origin string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string)

type NamedChannel interface {
	GetName() string
}

type SessionControlChannel interface {
	AutoEnableSession() bool
}

type LifecycleChannel interface {
	Close()
}

type IncomingMessageChannel interface {
	OnMessage(handler OnMessageHandler)
}

type MessageSender interface {
	SendMessage(chatID uint, message string) error
}

func ErrNotSupported(feature string) error {
	return fmt.Errorf("%s not supported by this channel", feature)
}

type Channel interface {
	NamedChannel
	LifecycleChannel
	SessionControlChannel

	IncomingMessageChannel
	MessageSender
}
