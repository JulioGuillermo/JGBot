package sessiondomain

type ChannelService interface {
	SendMessage(channelID string, chatID uint, message string) error
	SendReaction(channelID string, chatID uint, messageID uint, reaction string) error
	SendStatus(channelID string, chatID uint, status string) error
	GetChannelAutoEnable(channelID string) (bool, error)
}
