package sessioninfrastructure

import (
	"testing"

	channelsdomain "JGBot/channels/domain"
)

// mockChannelController implements channelsdomain.ChannelController for testing
type mockChannelController struct {
	channels      []channelsdomain.Channel
	getChannelErr error
	onMessageFn   func(handler channelsdomain.MessageHandler)
}

func (m *mockChannelController) SetChannels(channels []channelsdomain.Channel) {
	m.channels = channels
}

func (m *mockChannelController) GetChannel(name string) (channelsdomain.Channel, error) {
	for _, c := range m.channels {
		if c.GetName() == name {
			return c, nil
		}
	}
	return nil, m.getChannelErr
}

func (m *mockChannelController) SetChannel(channel channelsdomain.Channel) {}
func (m *mockChannelController) DelChannel(name string)                    {}

func (m *mockChannelController) OnMessage(handler channelsdomain.MessageHandler) {
	if m.onMessageFn != nil {
		m.onMessageFn(handler)
	}
}

func (m *mockChannelController) Close() {}

// mockChannel implements channelsdomain.Channel for testing
type mockChannel struct {
	nameValue       string
	autoEnableValue bool
	sendStatusFn    func(chatID uint, status channelsdomain.Status) error
	sendMessageFn   func(chatID uint, message string) error
	sendReactionFn  func(chatID uint, messageID uint, reaction string) error
	onMessageFn     func(handler channelsdomain.MessageHandler)
}

func (m *mockChannel) GetName() string { return m.nameValue }
func (m *mockChannel) OnMessage(handler channelsdomain.MessageHandler) {
	if m.onMessageFn != nil {
		m.onMessageFn(handler)
	}
}
func (m *mockChannel) SendStatus(chatID uint, status channelsdomain.Status) error {
	if m.sendStatusFn != nil {
		return m.sendStatusFn(chatID, status)
	}
	return nil
}
func (m *mockChannel) SendMessage(chatID uint, message string) error {
	if m.sendMessageFn != nil {
		return m.sendMessageFn(chatID, message)
	}
	return nil
}
func (m *mockChannel) SendMessageReaction(chatID uint, messageID uint, reaction string) error {
	if m.sendReactionFn != nil {
		return m.sendReactionFn(chatID, messageID, reaction)
	}
	return nil
}
func (m *mockChannel) AutoEnableSession() bool { return m.autoEnableValue }
func (m *mockChannel) Close()                  {}

func TestNewChannelAdapter(t *testing.T) {
	t.Parallel()

	ctrl := &mockChannelController{}
	adapter := NewChannelAdapter(ctrl)

	if adapter == nil {
		t.Error("NewChannelAdapter should not return nil")
	}
}

func TestChannelAdapter_SendMessage(t *testing.T) {
	t.Parallel()

	var sentChatID uint
	var sentMessage string

	mockCh := &mockChannel{
		nameValue: "telegram",
		sendMessageFn: func(chatID uint, message string) error {
			sentChatID = chatID
			sentMessage = message
			return nil
		},
	}

	ctrl := &mockChannelController{
		channels: []channelsdomain.Channel{mockCh},
	}

	adapter := NewChannelAdapter(ctrl)
	err := adapter.SendMessage("telegram", 123, "Hello")

	if err != nil {
		t.Errorf("SendMessage returned error: %v", err)
	}
	if sentChatID != 123 {
		t.Errorf("sentChatID = %v, want %v", sentChatID, 123)
	}
	if sentMessage != "Hello" {
		t.Errorf("sentMessage = %v, want %v", sentMessage, "Hello")
	}
}

func TestChannelAdapter_SendMessage_ChannelNotFound(t *testing.T) {
	t.Parallel()

	ctrl := &mockChannelController{
		getChannelErr: channelsdomain.ErrChannelNotFound,
	}

	adapter := NewChannelAdapter(ctrl)
	err := adapter.SendMessage("unknown", 123, "Hello")

	if err == nil {
		t.Error("SendMessage should return error for unknown channel")
	}
}

func TestChannelAdapter_SendReaction(t *testing.T) {
	t.Parallel()

	var sentChatID, sentMessageID uint
	var sentReaction string

	mockCh := &mockChannel{
		nameValue: "telegram",
		sendReactionFn: func(chatID uint, messageID uint, reaction string) error {
			sentChatID = chatID
			sentMessageID = messageID
			sentReaction = reaction
			return nil
		},
	}

	ctrl := &mockChannelController{
		channels: []channelsdomain.Channel{mockCh},
	}

	adapter := NewChannelAdapter(ctrl)
	err := adapter.SendReaction("telegram", 123, 456, "👍")

	if err != nil {
		t.Errorf("SendReaction returned error: %v", err)
	}
	if sentChatID != 123 {
		t.Errorf("sentChatID = %v, want %v", sentChatID, 123)
	}
	if sentMessageID != 456 {
		t.Errorf("sentMessageID = %v, want %v", sentMessageID, 456)
	}
	if sentReaction != "👍" {
		t.Errorf("sentReaction = %v, want %v", sentReaction, "👍")
	}
}

func TestChannelAdapter_SendStatus_Writing(t *testing.T) {
	t.Parallel()

	var sentStatus channelsdomain.Status

	mockCh := &mockChannel{
		nameValue: "telegram",
		sendStatusFn: func(chatID uint, status channelsdomain.Status) error {
			sentStatus = status
			return nil
		},
	}

	ctrl := &mockChannelController{
		channels: []channelsdomain.Channel{mockCh},
	}

	adapter := NewChannelAdapter(ctrl)
	err := adapter.SendStatus("telegram", 123, "writing")

	if err != nil {
		t.Errorf("SendStatus returned error: %v", err)
	}
	if sentStatus != channelsdomain.Writing {
		t.Errorf("sentStatus = %v, want %v (Writing)", sentStatus, channelsdomain.Writing)
	}
}

func TestChannelAdapter_SendStatus_Normal(t *testing.T) {
	t.Parallel()

	var sentStatus channelsdomain.Status

	mockCh := &mockChannel{
		nameValue: "telegram",
		sendStatusFn: func(chatID uint, status channelsdomain.Status) error {
			sentStatus = status
			return nil
		},
	}

	ctrl := &mockChannelController{
		channels: []channelsdomain.Channel{mockCh},
	}

	adapter := NewChannelAdapter(ctrl)
	err := adapter.SendStatus("telegram", 123, "normal")

	if err != nil {
		t.Errorf("SendStatus returned error: %v", err)
	}
	if sentStatus != channelsdomain.Normal {
		t.Errorf("sentStatus = %v, want %v (Normal)", sentStatus, channelsdomain.Normal)
	}
}

func TestChannelAdapter_GetChannelAutoEnable(t *testing.T) {
	t.Parallel()

	mockCh := &mockChannel{
		nameValue:       "telegram",
		autoEnableValue: true,
	}

	ctrl := &mockChannelController{
		channels: []channelsdomain.Channel{mockCh},
	}

	adapter := NewChannelAdapter(ctrl)
	enabled, err := adapter.GetChannelAutoEnable("telegram")

	if err != nil {
		t.Errorf("GetChannelAutoEnable returned error: %v", err)
	}
	if !enabled {
		t.Error("GetChannelAutoEnable should return true")
	}
}

func TestChannelAdapter_GetChannelAutoEnable_Disabled(t *testing.T) {
	t.Parallel()

	mockCh := &mockChannel{
		nameValue:       "telegram",
		autoEnableValue: false,
	}

	ctrl := &mockChannelController{
		channels: []channelsdomain.Channel{mockCh},
	}

	adapter := NewChannelAdapter(ctrl)
	enabled, err := adapter.GetChannelAutoEnable("telegram")

	if err != nil {
		t.Errorf("GetChannelAutoEnable returned error: %v", err)
	}
	if enabled {
		t.Error("GetChannelAutoEnable should return false")
	}
}
