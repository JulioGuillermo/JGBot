package sessiondomain

import (
	"testing"
)

func TestActivationContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ctx  *ActivationContext
	}{
		{
			name: "full activation context",
			ctx: &ActivationContext{
				Origin:      "telegram",
				Channel:     "general",
				ChatID:      123,
				ChatName:    "Test Chat",
				SenderID:    456,
				MessageID:   789,
				Name:        "daily_reminder",
				Schedule:    "0 9 * * *",
				Description: "Daily reminder",
				Message:     "Don't forget to drink water!",
			},
		},
		{
			name: "minimal activation context",
			ctx: &ActivationContext{
				Origin:      "whatsapp",
				Channel:     "alerts",
				ChatID:      999,
				ChatName:    "Alert Group",
				SenderID:    0,
				MessageID:   0,
				Name:        "timer_task",
				Schedule:    "",
				Description: "",
				Message:     "Timer triggered",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.ctx.Origin == "" {
				t.Error("Origin should not be empty")
			}
			if tt.ctx.Channel == "" {
				t.Error("Channel should not be empty")
			}
			if tt.ctx.Name == "" {
				t.Error("Name should not be empty")
			}
		})
	}
}

func TestActivationContextEquality(t *testing.T) {
	t.Parallel()

	ctx1 := &ActivationContext{
		Origin:      "telegram",
		Channel:     "general",
		ChatID:      123,
		ChatName:    "Test",
		SenderID:    456,
		MessageID:   789,
		Name:        "task1",
		Schedule:    "0 9 * * *",
		Description: "Desc",
		Message:     "Msg",
	}

	ctx2 := &ActivationContext{
		Origin:      "telegram",
		Channel:     "general",
		ChatID:      123,
		ChatName:    "Test",
		SenderID:    456,
		MessageID:   789,
		Name:        "task1",
		Schedule:    "0 9 * * *",
		Description: "Desc",
		Message:     "Msg",
	}

	ctx3 := &ActivationContext{
		Origin:      "discord",
		Channel:     "general",
		ChatID:      123,
		ChatName:    "Test",
		SenderID:    456,
		MessageID:   789,
		Name:        "task1",
		Schedule:    "0 9 * * *",
		Description: "Desc",
		Message:     "Msg",
	}

	if ctx1.Origin != ctx2.Origin {
		t.Error("ctx1 and ctx2 should have same origin")
	}
	if ctx1.Name != ctx2.Name {
		t.Error("ctx1 and ctx2 should have same name")
	}
	if ctx1.Origin == ctx3.Origin {
		t.Error("ctx1 and ctx3 should have different origins")
	}
}

func TestMessageContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ctx  *MessageContext
	}{
		{
			name: "full message context",
			ctx: &MessageContext{
				Origin:   "telegram",
				Channel:  "general",
				ChatID:   123,
				ChatName: "Test Chat",
				History: []*Message{
					{Message: "Hello"},
					{Message: "How are you?"},
				},
				IncomingMsg: &Message{Message: "Hi there!"},
				IsAdmin:     false,
				Config:      &SessionConfiguration{Origin: "telegram"},
				OnResponse: func(text, role, extra string) error {
					return nil
				},
				OnReact: func(msgID uint, reaction string) error {
					return nil
				},
			},
		},
		{
			name: "minimal message context",
			ctx: &MessageContext{
				Origin:      "whatsapp",
				Channel:     "chat",
				ChatID:      456,
				ChatName:    "Group",
				History:     nil,
				IncomingMsg: &Message{Message: "Test"},
				IsAdmin:     true,
				Config:      nil,
				OnResponse: func(text, role, extra string) error {
					return nil
				},
				OnReact: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.ctx.Origin == "" {
				t.Error("Origin should not be empty")
			}
			if tt.ctx.Channel == "" {
				t.Error("Channel should not be empty")
			}
			if tt.ctx.IncomingMsg == nil {
				t.Error("IncomingMsg should not be nil")
			}
			if tt.ctx.OnResponse == nil {
				t.Error("OnResponse should not be nil")
			}
		})
	}
}

func TestMessageContextOnResponse(t *testing.T) {
	t.Parallel()

	var called bool
	var receivedText, receivedRole, receivedExtra string

	ctx := &MessageContext{
		OnResponse: func(text, role, extra string) error {
			called = true
			receivedText = text
			receivedRole = role
			receivedExtra = extra
			return nil
		},
	}

	err := ctx.OnResponse("Hello", "assistant", `{"key": "value"}`)
	if err != nil {
		t.Errorf("OnResponse returned error: %v", err)
	}
	if !called {
		t.Error("OnResponse was not called")
	}
	if receivedText != "Hello" {
		t.Errorf("receivedText = %v, want %v", receivedText, "Hello")
	}
	if receivedRole != "assistant" {
		t.Errorf("receivedRole = %v, want %v", receivedRole, "assistant")
	}
	if receivedExtra != `{"key": "value"}` {
		t.Errorf("receivedExtra = %v, want %v", receivedExtra, `{"key": "value"}`)
	}
}

func TestMessageContextOnReact(t *testing.T) {
	t.Parallel()

	var called bool
	var receivedMsgID uint
	var receivedReaction string

	ctx := &MessageContext{
		OnReact: func(msgID uint, reaction string) error {
			called = true
			receivedMsgID = msgID
			receivedReaction = reaction
			return nil
		},
	}

	err := ctx.OnReact(123, "👍")
	if err != nil {
		t.Errorf("OnReact returned error: %v", err)
	}
	if !called {
		t.Error("OnReact was not called")
	}
	if receivedMsgID != 123 {
		t.Errorf("receivedMsgID = %v, want %v", receivedMsgID, 123)
	}
	if receivedReaction != "👍" {
		t.Errorf("receivedReaction = %v, want %v", receivedReaction, "👍")
	}
}

func TestMessageContextOnReactNil(t *testing.T) {
	t.Parallel()

	ctx := &MessageContext{
		OnReact: nil,
	}

	// Should not panic when OnReact is nil
	if ctx.OnReact != nil {
		err := ctx.OnReact(123, "👍")
		if err != nil {
			t.Errorf("OnReact returned error: %v", err)
		}
	}
}
