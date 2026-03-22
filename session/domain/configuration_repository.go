package sessiondomain

type ConfigurationRepository interface {
	GetConfig(origin string) *SessionConfiguration
	GetConfigByChannel(channel string, chatID uint) *SessionConfiguration
	CreateConfig(chatName, sessionID, origin, channel string) *SessionConfiguration
	CreateUnconfig(chatName, sessionID, origin, channel string) *SessionConfiguration
}
