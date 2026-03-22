package sessiondomain

type MessageRepository interface {
	Save(msg *Message) error
	GetHistory(channel string, chatID uint, limit int) ([]*Message, error)
	ClearHistory(channel string, chatID uint) error
}
