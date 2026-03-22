package sessiondomain

// MessageContext holds context for a single message interaction
type MessageContext struct {
	// Chat info
	Origin   string
	Channel  string
	ChatID   uint
	ChatName string

	// Message data
	History     []*Message
	IncomingMsg *Message

	// Permission
	IsAdmin bool

	// Configuration
	Config *SessionConfiguration

	// Callbacks
	OnResponse func(text, role, extra string) error
	OnReact    func(msgID uint, reaction string) error
}
