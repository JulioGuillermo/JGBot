package sessioninfrastructure

import (
	channelsdomain "JGBot/channels/domain"
)

// MockChannelController is a mock implementation of channelsdomain.ChannelController
type MockChannelController struct {
	Channels        []channelsdomain.Channel
	GetChannelFunc  func(name string) (channelsdomain.Channel, error)
	SetChannelsFunc func(channels []channelsdomain.Channel)
	SetChannelFunc  func(channel channelsdomain.Channel)
	DelChannelFunc  func(name string)
	OnMessageFunc   func(handler channelsdomain.MessageHandler)
	CloseFunc       func()
}

func (m *MockChannelController) SetChannels(channels []channelsdomain.Channel) {
	if m.SetChannelsFunc != nil {
		m.SetChannelsFunc(channels)
	}
	m.Channels = channels
}

func (m *MockChannelController) GetChannel(name string) (channelsdomain.Channel, error) {
	if m.GetChannelFunc != nil {
		return m.GetChannelFunc(name)
	}
	for _, c := range m.Channels {
		if c.GetName() == name {
			return c, nil
		}
	}
	return nil, nil
}

func (m *MockChannelController) SetChannel(channel channelsdomain.Channel) {
	if m.SetChannelFunc != nil {
		m.SetChannelFunc(channel)
	}
}

func (m *MockChannelController) DelChannel(name string) {
	if m.DelChannelFunc != nil {
		m.DelChannelFunc(name)
	}
}

func (m *MockChannelController) OnMessage(handler channelsdomain.MessageHandler) {
	if m.OnMessageFunc != nil {
		m.OnMessageFunc(handler)
	}
}

func (m *MockChannelController) Close() {
	if m.CloseFunc != nil {
		m.CloseFunc()
	}
}

// MockChannel is a mock implementation of channelsdomain.Channel
type MockChannel struct {
	NameValue               string
	AutoEnableValue         bool
	SendStatusFunc          func(chatID uint, status channelsdomain.Status) error
	SendMessageFunc         func(chatID uint, message string) error
	SendMessageReactionFunc func(chatID uint, messageID uint, reaction string) error
	OnMessageFunc           func(handler channelsdomain.MessageHandler)
	CloseFunc               func()
}

func (m *MockChannel) GetName() string {
	return m.NameValue
}

func (m *MockChannel) OnMessage(handler channelsdomain.MessageHandler) {
	if m.OnMessageFunc != nil {
		m.OnMessageFunc(handler)
	}
}

func (m *MockChannel) SendStatus(chatID uint, status channelsdomain.Status) error {
	if m.SendStatusFunc != nil {
		return m.SendStatusFunc(chatID, status)
	}
	return nil
}

func (m *MockChannel) SendMessage(chatID uint, message string) error {
	if m.SendMessageFunc != nil {
		return m.SendMessageFunc(chatID, message)
	}
	return nil
}

func (m *MockChannel) SendMessageReaction(chatID uint, messageID uint, reaction string) error {
	if m.SendMessageReactionFunc != nil {
		return m.SendMessageReactionFunc(chatID, messageID, reaction)
	}
	return nil
}

func (m *MockChannel) AutoEnableSession() bool {
	return m.AutoEnableValue
}

func (m *MockChannel) Close() {
	if m.CloseFunc != nil {
		m.CloseFunc()
	}
}
