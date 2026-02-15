package channels

type ReactionChannel interface {
	ReactMessage(chatID uint, messageID uint, reaction string) error
}
