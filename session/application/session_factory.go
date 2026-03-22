package sessionapplication

import (
	agentpkg "JGBot/agent"
	channelsdomain "JGBot/channels/domain"
	"JGBot/database"
	sessiondomain "JGBot/session/domain"
	infAgent "JGBot/session/infrastructure/agent"
	infChannel "JGBot/session/infrastructure/channel"
	infPersistence "JGBot/session/infrastructure/persistence"
	infDb "JGBot/session/infrastructure/persistence/db"
	infStore "JGBot/session/infrastructure/store"
)

// NewHexSessionService creates a new SessionService with hexagonal architecture.
// It wires up all the adapters and returns the service ready to use.
func NewHexSessionService(
	channelCtl channelsdomain.ChannelController,
	agentCtl *agentpkg.AgentsCtl,
) (*SessionService, error) {
	// Initialize database using new infrastructure db
	if err := infDb.Migrate(database.DB); err != nil {
		return nil, err
	}

	// Create pure domain implementations (no legacy dependencies)
	// 1. Message Repository (persistence layer) - uses new infrastructure
	msgRepo := infPersistence.NewGormMessageRepository(database.DB)

	// 2. Configuration Repository - pure domain implementation
	configRepo, err := infStore.NewFileConfigRepository()
	if err != nil {
		return nil, err
	}

	// 3. SessionStore - implements agentdomain.SessionStore for agent package
	sessionStore, err := infStore.NewFileSessionStore(msgRepo)
	if err != nil {
		return nil, err
	}

	// 4. Channel Service (adapter wrapping existing channel controller)
	channelAdapter := infChannel.NewChannelAdapter(channelCtl)

	// 5. Agent Service (adapter bridging to existing agent)
	agentAdapter := infAgent.NewAgentAdapter(agentCtl, sessionStore)

	// Create application service with all dependencies injected
	sessionService := NewSessionService(
		msgRepo,
		configRepo,
		agentAdapter,
		channelAdapter,
	)

	// Register message handler to channel controller
	channelCtl.OnMessage(func(channel, origin string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string) {
		sessionService.OnNewMessage(channel, origin, chatID, chatName, senderID, senderName, messageID, message)
	})

	// Set agent dependencies
	agentCtl.SetDependencies(sessionStore, channelCtl)

	return sessionService, nil
}

// GetTaskHandler returns the TaskHandler Interface for auto-activation
func GetTaskHandler(svc *SessionService) sessiondomain.TaskHandler {
	return &taskHandlerAdapter{service: svc}
}

type taskHandlerAdapter struct {
	service *SessionService
}

func (a *taskHandlerAdapter) OnActivation(ctx *sessiondomain.ActivationContext) {
	a.service.OnAutoActivation(ctx)
}
