package sessioninfrastructure

import (
	"testing"

	agentdomain "JGBot/agent/domain"
	sessiondomain "JGBot/session/domain"
)

func TestAgentAdapter_ImplementsInterface(t *testing.T) {
	t.Parallel()

	// Compile-time check that AgentAdapter implements domain.AgentService
	var _ sessiondomain.AgentService = (*AgentAdapter)(nil)
}

func TestAgentAdapter_toAgentMessage(t *testing.T) {
	t.Parallel()

	adapter := &AgentAdapter{}

	msg := &sessiondomain.Message{
		Channel:    "telegram",
		ChatID:     123,
		ChatName:   "Test Chat",
		SenderID:   456,
		SenderName: "John",
		MessageID:  789,
		Message:    "Hello",
		Role:       "user",
		Extra:      "",
	}

	agentMsg := adapter.toAgentMessage(msg)

	if agentMsg == nil {
		t.Fatal("toAgentMessage should not return nil")
	}
	if agentMsg.Channel != msg.Channel {
		t.Errorf("Channel = %v, want %v", agentMsg.Channel, msg.Channel)
	}
	if agentMsg.ChatID != msg.ChatID {
		t.Errorf("ChatID = %v, want %v", agentMsg.ChatID, msg.ChatID)
	}
	if agentMsg.Message != msg.Message {
		t.Errorf("Message = %v, want %v", agentMsg.Message, msg.Message)
	}
	if agentMsg.Role != msg.Role {
		t.Errorf("Role = %v, want %v", agentMsg.Role, msg.Role)
	}
}

func TestAgentAdapter_toAgentMessage_Nil(t *testing.T) {
	t.Parallel()

	adapter := &AgentAdapter{}
	agentMsg := adapter.toAgentMessage(nil)

	if agentMsg != nil {
		t.Error("toAgentMessage should return nil for nil input")
	}
}

func TestAgentAdapter_toAgentMessage_AllRoles(t *testing.T) {
	t.Parallel()

	adapter := &AgentAdapter{}

	roles := []string{"user", "assistant", "system", "tool"}

	for _, role := range roles {
		t.Run(role, func(t *testing.T) {
			msg := &sessiondomain.Message{
				Channel: "telegram",
				ChatID:  123,
				Message: "Test",
				Role:    role,
				Extra:   "",
			}

			agentMsg := adapter.toAgentMessage(msg)

			if agentMsg.Role != role {
				t.Errorf("Role = %v, want %v", agentMsg.Role, role)
			}
		})
	}
}

func TestAgentAdapter_toAgentMessage_PreservesAllFields(t *testing.T) {
	t.Parallel()

	adapter := &AgentAdapter{}

	msg := &sessiondomain.Message{
		ID:         100,
		Channel:    "telegram",
		ChatID:     123,
		ChatName:   "Test Chat",
		SenderID:   456,
		SenderName: "John",
		MessageID:  789,
		Message:    "Hello",
		Role:       "assistant",
		Extra:      `{"key": "value"}`,
	}

	agentMsg := adapter.toAgentMessage(msg)

	if agentMsg.ID != msg.ID {
		t.Errorf("ID = %v, want %v", agentMsg.ID, msg.ID)
	}
	if agentMsg.ChatName != msg.ChatName {
		t.Errorf("ChatName = %v, want %v", agentMsg.ChatName, msg.ChatName)
	}
	if agentMsg.SenderID != msg.SenderID {
		t.Errorf("SenderID = %v, want %v", agentMsg.SenderID, msg.SenderID)
	}
	if agentMsg.SenderName != msg.SenderName {
		t.Errorf("SenderName = %v, want %v", agentMsg.SenderName, msg.SenderName)
	}
	if agentMsg.MessageID != msg.MessageID {
		t.Errorf("MessageID = %v, want %v", agentMsg.MessageID, msg.MessageID)
	}
	if agentMsg.Extra != msg.Extra {
		t.Errorf("Extra = %v, want %v", agentMsg.Extra, msg.Extra)
	}

	// Ensure the types are correct
	var _ *agentdomain.SessionMessage = agentMsg
}
