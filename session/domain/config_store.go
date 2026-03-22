package sessiondomain

// ConfigStore provides access to session configurations
type ConfigStore interface {
	GetConfig(origin string) *SessionConfiguration
	GetConfigs() []*SessionConfiguration
	CreateConfig(chatName, sessionID, origin, channel string) *SessionConfiguration
	CreateUnconfig(chatName, sessionID, origin, channel string) *SessionConfiguration
	Close()
}
