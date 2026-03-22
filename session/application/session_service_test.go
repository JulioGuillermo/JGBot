package sessionapplication

import (
	"errors"
	"testing"

	sessiondomain "JGBot/session/domain"

	"github.com/tmc/langchaingo/agents"
)

func TestNewSessionService(t *testing.T) {
	t.Parallel()

	repo := &MockMessageRepository{}
	configRepo := &MockConfigurationRepository{}
	agentSvc := &MockAgentService{}
	channelSvc := &MockChannelService{}

	svc := NewSessionService(repo, configRepo, agentSvc, channelSvc)

	if svc == nil {
		t.Error("NewSessionService should not return nil")
	}
	if svc.repo != repo {
		t.Error("repo should be set")
	}
	if svc.configRepo != configRepo {
		t.Error("configRepo should be set")
	}
	if svc.agentService != agentSvc {
		t.Error("agentService should be set")
	}
	if svc.channelSvc != channelSvc {
		t.Error("channelSvc should be set")
	}
}

func TestOnNewMessage_EmptyMessage(t *testing.T) {
	t.Parallel()

	repo := &MockMessageRepository{}
	configRepo := &MockConfigurationRepository{}
	agentSvc := &MockAgentService{}
	channelSvc := &MockChannelService{}

	svc := NewSessionService(repo, configRepo, agentSvc, channelSvc)

	// Should not panic with empty message
	svc.OnNewMessage("telegram", "origin", 123, "Chat", 456, "User", 789, "")
}

func TestOnNewMessage_GetChannelAutoEnableError(t *testing.T) {
	t.Parallel()

	channelSvc := &MockChannelService{
		GetChannelAutoEnableFunc: func(channelID string) (bool, error) {
			return false, errors.New("channel not found")
		},
	}

	svc := NewSessionService(&MockMessageRepository{}, &MockConfigurationRepository{}, &MockAgentService{}, channelSvc)

	// Should handle error gracefully
	svc.OnNewMessage("telegram", "origin", 123, "Chat", 456, "User", 789, "Hello")
}

func TestOnNewMessage_NoConfigAndNotAutoEnable(t *testing.T) {
	t.Parallel()

	channelSvc := &MockChannelService{
		GetChannelAutoEnableFunc: func(channelID string) (bool, error) {
			return false, nil // Not auto-enabled
		},
	}
	configRepo := &MockConfigurationRepository{
		GetConfigFunc: func(origin string) *sessiondomain.SessionConfiguration {
			return nil // No config
		},
		CreateUnconfigFunc: func(chatName, sessionID, origin, channel string) *sessiondomain.SessionConfiguration {
			return nil
		},
	}

	svc := NewSessionService(&MockMessageRepository{}, configRepo, &MockAgentService{}, channelSvc)
	svc.OnNewMessage("telegram", "origin", 123, "Chat", 456, "User", 789, "Hello")
}

func TestOnNewMessage_WithConfig(t *testing.T) {
	t.Parallel()

	var history []*sessiondomain.Message
	var agentCalled bool
	var incomingMsg *sessiondomain.Message
	var responseMsg *sessiondomain.Message

	channelSvc := &MockChannelService{
		GetChannelAutoEnableFunc: func(channelID string) (bool, error) {
			return true, nil
		},
		SendMessageFunc: func(channelID string, chatID uint, message string) error {
			return nil
		},
		SendStatusFunc: func(channelID string, chatID uint, status string) error {
			return nil
		},
	}

	configRepo := &MockConfigurationRepository{
		GetConfigFunc: func(origin string) *sessiondomain.SessionConfiguration {
			return &sessiondomain.SessionConfiguration{
				Origin:      "telegram",
				ChatID:      "123",
				ChatName:    "Chat",
				HistorySize: 50,
				Admin:       "",
				Allowed:     true,
				ShouldRespond: func(msg string) bool {
					return true
				},
			}
		},
	}

	repo := &MockMessageRepository{
		GetHistoryFunc: func(channel string, chatID uint, limit int) ([]*sessiondomain.Message, error) {
			return history, nil
		},
		SaveFunc: func(msg *sessiondomain.Message) error {
			if msg.Role == "user" {
				incomingMsg = msg
			} else {
				responseMsg = msg
			}
			return nil
		},
	}

	agentSvc := &MockAgentService{
		RespondFunc: func(ctx *sessiondomain.MessageContext) error {
			agentCalled = true
			// Call the response callback to simulate agent response
			ctx.OnResponse("Hello!", "assistant", "")
			return nil
		},
	}

	svc := NewSessionService(repo, configRepo, agentSvc, channelSvc)
	svc.OnNewMessage("telegram", "origin", 123, "Chat", 456, "User", 789, "Hi there")

	if !agentCalled {
		t.Error("Agent should have been called")
	}
	if incomingMsg == nil {
		t.Error("Incoming message should have been saved")
	}
	if incomingMsg != nil && incomingMsg.Role != "user" {
		t.Errorf("Incoming message role = %v, want user", incomingMsg.Role)
	}
	if responseMsg == nil {
		t.Error("Response message should have been saved")
	}
	if responseMsg != nil && responseMsg.Role != "assistant" {
		t.Errorf("Response message role = %v, want assistant", responseMsg.Role)
	}
}

func TestOnNewMessage_ResetCommand(t *testing.T) {
	t.Parallel()

	var clearCalled bool

	channelSvc := &MockChannelService{
		GetChannelAutoEnableFunc: func(channelID string) (bool, error) {
			return true, nil
		},
		SendMessageFunc: func(channelID string, chatID uint, message string) error {
			return nil
		},
	}

	configRepo := &MockConfigurationRepository{
		GetConfigFunc: func(origin string) *sessiondomain.SessionConfiguration {
			return &sessiondomain.SessionConfiguration{
				Origin:  "telegram",
				Allowed: true,
				ShouldRespond: func(msg string) bool {
					return true
				},
			}
		},
	}

	repo := &MockMessageRepository{
		ClearFunc: func(channel string, chatID uint) error {
			clearCalled = true
			return nil
		},
	}

	svc := NewSessionService(repo, configRepo, &MockAgentService{}, channelSvc)
	svc.OnNewMessage("telegram", "origin", 123, "Chat", 456, "User", 789, "/reset!")

	if !clearCalled {
		t.Error("ClearHistory should have been called")
	}
}

func TestOnNewMessage_ShouldNotRespond(t *testing.T) {
	t.Parallel()

	var agentCalled bool

	channelSvc := &MockChannelService{
		GetChannelAutoEnableFunc: func(channelID string) (bool, error) {
			return true, nil
		},
	}

	configRepo := &MockConfigurationRepository{
		GetConfigFunc: func(origin string) *sessiondomain.SessionConfiguration {
			return &sessiondomain.SessionConfiguration{
				Origin:  "telegram",
				Allowed: true,
				ShouldRespond: func(msg string) bool {
					return false // Should not respond
				},
			}
		},
	}

	repo := &MockMessageRepository{}

	agentSvc := &MockAgentService{
		RespondFunc: func(ctx *sessiondomain.MessageContext) error {
			agentCalled = true
			return nil
		},
	}

	svc := NewSessionService(repo, configRepo, agentSvc, channelSvc)
	svc.OnNewMessage("telegram", "origin", 123, "Chat", 456, "User", 789, "skip")

	if agentCalled {
		t.Error("Agent should not have been called when ShouldRespond returns false")
	}
}

func TestOnNewMessage_MaxIterationsError(t *testing.T) {
	t.Parallel()

	var sentMessage string

	channelSvc := &MockChannelService{
		GetChannelAutoEnableFunc: func(channelID string) (bool, error) {
			return true, nil
		},
		SendMessageFunc: func(channelID string, chatID uint, message string) error {
			sentMessage = message
			return nil
		},
		SendStatusFunc: func(channelID string, chatID uint, status string) error {
			return nil
		},
	}

	configRepo := &MockConfigurationRepository{
		GetConfigFunc: func(origin string) *sessiondomain.SessionConfiguration {
			return &sessiondomain.SessionConfiguration{
				Origin:  "telegram",
				Allowed: true,
				ShouldRespond: func(msg string) bool {
					return true
				},
			}
		},
	}

	repo := &MockMessageRepository{}

	agentSvc := &MockAgentService{
		RespondFunc: func(ctx *sessiondomain.MessageContext) error {
			return agents.ErrNotFinished
		},
	}

	svc := NewSessionService(repo, configRepo, agentSvc, channelSvc)
	svc.OnNewMessage("telegram", "origin", 123, "Chat", 456, "User", 789, "test")

	if sentMessage == "" {
		t.Error("Error message should have been sent")
	}
}

func TestOnNewMessage_GeneralError(t *testing.T) {
	t.Parallel()

	var sentMessage string

	channelSvc := &MockChannelService{
		GetChannelAutoEnableFunc: func(channelID string) (bool, error) {
			return true, nil
		},
		SendMessageFunc: func(channelID string, chatID uint, message string) error {
			sentMessage = message
			return nil
		},
		SendStatusFunc: func(channelID string, chatID uint, status string) error {
			return nil
		},
	}

	configRepo := &MockConfigurationRepository{
		GetConfigFunc: func(origin string) *sessiondomain.SessionConfiguration {
			return &sessiondomain.SessionConfiguration{
				Origin:  "telegram",
				Allowed: true,
				ShouldRespond: func(msg string) bool {
					return true
				},
			}
		},
	}

	repo := &MockMessageRepository{}

	agentSvc := &MockAgentService{
		RespondFunc: func(ctx *sessiondomain.MessageContext) error {
			return errors.New("some error")
		},
	}

	svc := NewSessionService(repo, configRepo, agentSvc, channelSvc)
	svc.OnNewMessage("telegram", "origin", 123, "Chat", 456, "User", 789, "test")

	if sentMessage == "" {
		t.Error("Error message should have been sent")
	}
}

func TestIsAdmin(t *testing.T) {
	t.Parallel()

	svc := NewSessionService(&MockMessageRepository{}, &MockConfigurationRepository{}, &MockAgentService{}, &MockChannelService{})

	tests := []struct {
		name             string
		adminPermission  string
		message          string
		wantIsAdmin      bool
		wantCleanMessage string
	}{
		{
			name:             "full permission",
			adminPermission:  "full",
			message:          "any message",
			wantIsAdmin:      true,
			wantCleanMessage: "any message",
		},
		{
			name:             "allow permission with admin prefix",
			adminPermission:  "allow",
			message:          "/admin secret command",
			wantIsAdmin:      true,
			wantCleanMessage: "secret command",
		},
		{
			name:             "allow permission without admin prefix",
			adminPermission:  "allow",
			message:          "regular message",
			wantIsAdmin:      false,
			wantCleanMessage: "regular message",
		},
		{
			name:             "no permission",
			adminPermission:  "",
			message:          "/admin anything",
			wantIsAdmin:      false,
			wantCleanMessage: "anything", // prefix is stripped but admin is false
		},
		{
			name:             "admin command alone",
			adminPermission:  "allow",
			message:          "/admin",
			wantIsAdmin:      true,
			wantCleanMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			isAdmin, cleanMsg := svc.isAdmin(tt.adminPermission, tt.message)
			if isAdmin != tt.wantIsAdmin {
				t.Errorf("isAdmin = %v, want %v", isAdmin, tt.wantIsAdmin)
			}
			if cleanMsg != tt.wantCleanMessage {
				t.Errorf("cleanMessage = %v, want %v", cleanMsg, tt.wantCleanMessage)
			}
		})
	}
}

func TestOnAutoActivation(t *testing.T) {
	t.Parallel()

	var agentCalled bool

	channelSvc := &MockChannelService{
		GetChannelAutoEnableFunc: func(channelID string) (bool, error) {
			return true, nil
		},
		SendMessageFunc: func(channelID string, chatID uint, message string) error {
			return nil
		},
		SendStatusFunc: func(channelID string, chatID uint, status string) error {
			return nil
		},
	}

	configRepo := &MockConfigurationRepository{
		GetConfigFunc: func(origin string) *sessiondomain.SessionConfiguration {
			return &sessiondomain.SessionConfiguration{
				Origin:  "telegram",
				Allowed: true,
				ShouldRespond: func(msg string) bool {
					return true
				},
			}
		},
	}

	repo := &MockMessageRepository{}

	agentSvc := &MockAgentService{
		RespondFunc: func(ctx *sessiondomain.MessageContext) error {
			agentCalled = true
			return nil
		},
	}

	svc := NewSessionService(repo, configRepo, agentSvc, channelSvc)

	activationCtx := &sessiondomain.ActivationContext{
		Origin:      "telegram",
		Channel:     "general",
		ChatID:      123,
		ChatName:    "Test Chat",
		SenderID:    456,
		MessageID:   789,
		Name:        "daily_task",
		Schedule:    "0 9 * * *",
		Description: "Daily task",
		Message:     "Hello from cron!",
	}

	svc.OnAutoActivation(activationCtx)

	if !agentCalled {
		t.Error("Agent should have been called for auto activation")
	}
}

func TestOnAutoActivation_NotAllowed(t *testing.T) {
	t.Parallel()

	var agentCalled bool

	channelSvc := &MockChannelService{
		GetChannelAutoEnableFunc: func(channelID string) (bool, error) {
			return true, nil
		},
	}

	configRepo := &MockConfigurationRepository{
		GetConfigFunc: func(origin string) *sessiondomain.SessionConfiguration {
			return &sessiondomain.SessionConfiguration{
				Allowed: false,
			}
		},
	}

	agentSvc := &MockAgentService{
		RespondFunc: func(ctx *sessiondomain.MessageContext) error {
			agentCalled = true
			return nil
		},
	}

	svc := NewSessionService(&MockMessageRepository{}, configRepo, agentSvc, channelSvc)

	activationCtx := &sessiondomain.ActivationContext{
		Origin:   "telegram",
		Channel:  "general",
		ChatID:   123,
		ChatName: "Test",
	}

	svc.OnAutoActivation(activationCtx)

	if agentCalled {
		t.Error("Agent should not have been called when session is not allowed")
	}
}
