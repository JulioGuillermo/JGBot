package sessiondomain

import (
	"testing"
	"time"
)

func TestMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		msg  *Message
		want string
	}{
		{
			name: "user message",
			msg: &Message{
				ID:         1,
				Channel:    "telegram",
				ChatID:     123,
				ChatName:   "Test Chat",
				SenderID:   456,
				SenderName: "John",
				MessageID:  789,
				Message:    "Hello world",
				Role:       "user",
				Extra:      "",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			want: "SOURCE_CHANNEL: telegram\nUSER_NAME: John\nCHAT_NAME: Test Chat\nMESSAGE_ID: 789\n\nCONTENT: Hello world",
		},
		{
			name: "assistant message",
			msg: &Message{
				ID:         2,
				Channel:    "whatsapp",
				ChatID:     456,
				ChatName:   "Group Chat",
				SenderID:   0,
				SenderName: "AI",
				MessageID:  101,
				Message:    "Hello! How can I help?",
				Role:       "assistant",
				Extra:      "",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			want: "SOURCE_CHANNEL: whatsapp\nUSER_NAME: ASSISTANT\nCHAT_NAME: Group Chat\nMESSAGE_ID: 101\n\nCONTENT: Hello! How can I help?",
		},
		{
			name: "system message",
			msg: &Message{
				ID:         3,
				Channel:    "telegram",
				ChatID:     789,
				ChatName:   "Admin Chat",
				SenderID:   0,
				SenderName: "System",
				MessageID:  202,
				Message:    "System initialized",
				Role:       "system",
				Extra:      "",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			want: "SOURCE_CHANNEL: telegram\nUSER_NAME: SYSTEM\nCHAT_NAME: Admin Chat\nMESSAGE_ID: 202\n\nCONTENT: System initialized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.msg.String()
			if got != tt.want {
				t.Errorf("Message.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		channel    string
		chatID     uint
		chatName   string
		senderID   uint
		senderName string
		messageID  uint
		message    string
		role       string
		extra      string
		wantRole   string
		wantMsg    string
	}{
		{
			name:       "create user message",
			channel:    "telegram",
			chatID:     123,
			chatName:   "Test Chat",
			senderID:   456,
			senderName: "John",
			messageID:  789,
			message:    "Hello",
			role:       "user",
			extra:      "",
			wantRole:   "user",
			wantMsg:    "Hello",
		},
		{
			name:       "create tool message",
			channel:    "whatsapp",
			chatID:     321,
			chatName:   "Group",
			senderID:   111,
			senderName: "Bot",
			messageID:  222,
			message:    "Tool output",
			role:       "tool",
			extra:      `{"tool": "test"}`,
			wantRole:   "tool",
			wantMsg:    "Tool output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			msg := NewMessage(tt.channel, tt.chatID, tt.chatName, tt.senderID, tt.senderName, tt.messageID, tt.message, tt.role, tt.extra)

			if msg.Channel != tt.channel {
				t.Errorf("Channel = %v, want %v", msg.Channel, tt.channel)
			}
			if msg.ChatID != tt.chatID {
				t.Errorf("ChatID = %v, want %v", msg.ChatID, tt.chatID)
			}
			if msg.Role != tt.wantRole {
				t.Errorf("Role = %v, want %v", msg.Role, tt.wantRole)
			}
			if msg.Message != tt.wantMsg {
				t.Errorf("Message = %v, want %v", msg.Message, tt.wantMsg)
			}
			if msg.CreatedAt.IsZero() {
				t.Error("CreatedAt should not be zero")
			}
		})
	}
}

func TestMessageEquality(t *testing.T) {
	t.Parallel()

	now := time.Now()
	msg1 := &Message{
		ID:         1,
		Channel:    "telegram",
		ChatID:     123,
		ChatName:   "Test",
		SenderID:   456,
		SenderName: "John",
		MessageID:  789,
		Message:    "Hello",
		Role:       "user",
		Extra:      "",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	msg2 := &Message{
		ID:         1,
		Channel:    "telegram",
		ChatID:     123,
		ChatName:   "Test",
		SenderID:   456,
		SenderName: "John",
		MessageID:  789,
		Message:    "Hello",
		Role:       "user",
		Extra:      "",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	msg3 := &Message{
		ID:         2,
		Channel:    "whatsapp",
		ChatID:     123,
		ChatName:   "Test",
		SenderID:   456,
		SenderName: "John",
		MessageID:  789,
		Message:    "Hello",
		Role:       "user",
		Extra:      "",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if msg1.Channel != msg2.Channel {
		t.Error("msg1 and msg2 should have same channel")
	}
	if msg1.Message != msg2.Message {
		t.Error("msg1 and msg2 should have same message")
	}
	if msg1.Channel == msg3.Channel {
		t.Error("msg1 and msg3 should have different channels")
	}
}
