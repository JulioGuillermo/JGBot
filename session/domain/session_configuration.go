package sessiondomain

type SessionConfiguration struct {
	Origin      string
	Channel     string
	ChatID      string // Stored as string "channel:id" usually, or just ID
	ChatName    string
	HistorySize int
	Admin       string
	Allowed     bool
	// Respond policy is part of config, logic to check if we should respond
	ShouldRespond func(message string) bool
}
