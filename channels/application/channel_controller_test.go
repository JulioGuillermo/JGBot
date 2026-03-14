package channelsapplication

import (
	"errors"
	"testing"

	channelsdomain "JGBot/channels/domain"
)

// mockChannel implements channelsdomain.Channel for testing
type mockChannel struct {
	name            string
	autoEnable      bool
	onMessageCalled bool
	closeCalled     bool
	handler         channelsdomain.MessageHandler
}

func (m *mockChannel) GetName() string         { return m.name }
func (m *mockChannel) AutoEnableSession() bool { return m.autoEnable }
func (m *mockChannel) OnMessage(handler channelsdomain.MessageHandler) {
	m.handler = handler
	m.onMessageCalled = true
}
func (m *mockChannel) SendStatus(chatID uint, status channelsdomain.Status) error { return nil }
func (m *mockChannel) SendMessage(chatID uint, message string) error              { return nil }
func (m *mockChannel) SendMessageReaction(chatID uint, messageID uint, reaction string) error {
	return nil
}
func (m *mockChannel) Close() { m.closeCalled = true }

func TestChannelController_SetChannels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		channels []channelsdomain.Channel
		wantLen  int
	}{
		{
			name:     "empty slice",
			channels: []channelsdomain.Channel{},
			wantLen:  0,
		},
		{
			name:     "single channel",
			channels: []channelsdomain.Channel{&mockChannel{name: "channel1"}},
			wantLen:  1,
		},
		{
			name:     "multiple channels",
			channels: []channelsdomain.Channel{&mockChannel{name: "channel1"}, &mockChannel{name: "channel2"}},
			wantLen:  2,
		},
		{
			name:     "duplicate names overwrites",
			channels: []channelsdomain.Channel{&mockChannel{name: "channel1"}, &mockChannel{name: "channel1"}},
			wantLen:  1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctl := NewChannelController()
			ctl.SetChannels(tt.channels)

			if len(ctl.(*ChannelController).channels) != tt.wantLen {
				t.Errorf("SetChannels() got len = %d, want %d", len(ctl.(*ChannelController).channels), tt.wantLen)
			}
		})
	}
}

func TestChannelController_GetChannel(t *testing.T) {
	t.Parallel()

	existingChannel := &mockChannel{name: "channel1"}
	ctl := NewChannelController()
	ctl.SetChannels([]channelsdomain.Channel{existingChannel})

	tests := []struct {
		name    string
		channel string
		want    channelsdomain.Channel
		wantErr error
	}{
		{
			name:    "existing channel",
			channel: "channel1",
			want:    existingChannel,
			wantErr: nil,
		},
		{
			name:    "non-existing channel",
			channel: "channel2",
			want:    nil,
			wantErr: channelsdomain.ErrChannelNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ctl.GetChannel(tt.channel)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetChannel() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelController_SetChannel(t *testing.T) {
	t.Parallel()

	ctl := NewChannelController()

	newChannel := &mockChannel{name: "channel1"}
	ctl.SetChannel(newChannel)

	got, err := ctl.GetChannel("channel1")
	if err != nil {
		t.Fatalf("GetChannel() unexpected error: %v", err)
	}
	if got.GetName() != "channel1" {
		t.Errorf("SetChannel() channel name = %s, want channel1", got.GetName())
	}
}

func TestChannelController_DelChannel(t *testing.T) {
	t.Parallel()

	ctl := NewChannelController()
	ctl.SetChannel(&mockChannel{name: "channel1"})

	// Verify channel exists
	_, err := ctl.GetChannel("channel1")
	if err != nil {
		t.Fatal("channel should exist before deletion")
	}

	// Delete channel
	ctl.DelChannel("channel1")

	// Verify channel is gone
	_, err = ctl.GetChannel("channel1")
	if !errors.Is(err, channelsdomain.ErrChannelNotFound) {
		t.Error("DelChannel() should remove channel")
	}
}

func TestChannelController_OnMessage(t *testing.T) {
	t.Parallel()

	ch1 := &mockChannel{name: "channel1"}
	ch2 := &mockChannel{name: "channel2"}

	ctl := NewChannelController()
	ctl.SetChannels([]channelsdomain.Channel{ch1, ch2})

	handler := func(channel string, origin string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string) {
	}
	ctl.OnMessage(handler)

	if !ch1.onMessageCalled || !ch2.onMessageCalled {
		t.Error("OnMessage() should register handler on all channels")
	}
}

func TestChannelController_Close(t *testing.T) {
	t.Parallel()

	ch1 := &mockChannel{name: "channel1"}
	ch2 := &mockChannel{name: "channel2"}

	ctl := NewChannelController()
	ctl.SetChannels([]channelsdomain.Channel{ch1, ch2})

	ctl.Close()

	if !ch1.closeCalled || !ch2.closeCalled {
		t.Error("Close() should call Close() on all channels")
	}
}
