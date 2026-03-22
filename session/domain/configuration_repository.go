package sessiondomain

type ConfigurationRepository interface {
	GetConfig(origin string) *SessionConfiguration
	CreateConfig(chatName, sessionID, origin, channel string) *SessionConfiguration
	CreateUnconfig(chatName, sessionID, origin, channel string) *SessionConfiguration
}
