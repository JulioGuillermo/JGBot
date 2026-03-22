package sessiondomain

// ActivationContext holds all context needed for task activation (cron, timer)
type ActivationContext struct {
	// Origin info
	Origin    string
	Channel   string
	ChatID    uint
	ChatName  string
	SenderID  uint
	MessageID uint

	// Task info
	Name        string
	Schedule    string
	Description string
	Message     string
}
