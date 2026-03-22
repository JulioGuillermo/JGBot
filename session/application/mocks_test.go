package sessionapplication

import (
	sessiondomain "JGBot/session/domain"
)

// MockMessageRepository is a mock implementation of sessiondomain.MessageRepository
type MockMessageRepository struct {
	SaveFunc       func(msg *sessiondomain.Message) error
	GetHistoryFunc func(channel string, chatID uint, limit int) ([]*sessiondomain.Message, error)
	ClearFunc      func(channel string, chatID uint) error
}

func (m *MockMessageRepository) Save(msg *sessiondomain.Message) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(msg)
	}
	return nil
}

func (m *MockMessageRepository) GetHistory(channel string, chatID uint, limit int) ([]*sessiondomain.Message, error) {
	if m.GetHistoryFunc != nil {
		return m.GetHistoryFunc(channel, chatID, limit)
	}
	return nil, nil
}

func (m *MockMessageRepository) ClearHistory(channel string, chatID uint) error {
	if m.ClearFunc != nil {
		return m.ClearFunc(channel, chatID)
	}
	return nil
}

// MockConfigurationRepository is a mock implementation of sessiondomain.ConfigurationRepository
type MockConfigurationRepository struct {
	GetConfigFunc      func(origin string) *sessiondomain.SessionConfiguration
	CreateConfigFunc   func(chatName, sessionID, origin, channel string) *sessiondomain.SessionConfiguration
	CreateUnconfigFunc func(chatName, sessionID, origin, channel string) *sessiondomain.SessionConfiguration
}

func (m *MockConfigurationRepository) GetConfig(origin string) *sessiondomain.SessionConfiguration {
	if m.GetConfigFunc != nil {
		return m.GetConfigFunc(origin)
	}
	return nil
}

func (m *MockConfigurationRepository) CreateConfig(chatName, sessionID, origin, channel string) *sessiondomain.SessionConfiguration {
	if m.CreateConfigFunc != nil {
		return m.CreateConfigFunc(chatName, sessionID, origin, channel)
	}
	return nil
}

func (m *MockConfigurationRepository) CreateUnconfig(chatName, sessionID, origin, channel string) *sessiondomain.SessionConfiguration {
	if m.CreateUnconfigFunc != nil {
		return m.CreateUnconfigFunc(chatName, sessionID, origin, channel)
	}
	return nil
}

// MockChannelService is a mock implementation of sessiondomain.ChannelService
type MockChannelService struct {
	SendMessageFunc          func(channelID string, chatID uint, message string) error
	SendReactionFunc         func(channelID string, chatID uint, messageID uint, reaction string) error
	SendStatusFunc           func(channelID string, chatID uint, status string) error
	GetChannelAutoEnableFunc func(channelID string) (bool, error)
}

func (m *MockChannelService) SendMessage(channelID string, chatID uint, message string) error {
	if m.SendMessageFunc != nil {
		return m.SendMessageFunc(channelID, chatID, message)
	}
	return nil
}

func (m *MockChannelService) SendReaction(channelID string, chatID uint, messageID uint, reaction string) error {
	if m.SendReactionFunc != nil {
		return m.SendReactionFunc(channelID, chatID, messageID, reaction)
	}
	return nil
}

func (m *MockChannelService) SendStatus(channelID string, chatID uint, status string) error {
	if m.SendStatusFunc != nil {
		return m.SendStatusFunc(channelID, chatID, status)
	}
	return nil
}

func (m *MockChannelService) GetChannelAutoEnable(channelID string) (bool, error) {
	if m.GetChannelAutoEnableFunc != nil {
		return m.GetChannelAutoEnableFunc(channelID)
	}
	return false, nil
}

// MockAgentService is a mock implementation of sessiondomain.AgentService
type MockAgentService struct {
	RespondFunc func(ctx *sessiondomain.MessageContext) error
}

func (m *MockAgentService) Respond(ctx *sessiondomain.MessageContext) error {
	if m.RespondFunc != nil {
		return m.RespondFunc(ctx)
	}
	return nil
}
