package sessioninfrastructure

import (
	agentpkg "JGBot/agent"
	agentdomain "JGBot/agent/domain"
	"JGBot/ctxs"
	sessiondomain "JGBot/session/domain"
)

type AgentAdapter struct {
	agentCtl     *agentpkg.AgentsCtl
	sessionStore agentdomain.SessionStore
}

func NewAgentAdapter(agentCtl *agentpkg.AgentsCtl, sessionStore agentdomain.SessionStore) *AgentAdapter {
	return &AgentAdapter{
		agentCtl:     agentCtl,
		sessionStore: sessionStore,
	}
}

func (a *AgentAdapter) Respond(ctx *sessiondomain.MessageContext) error {
	// Get the full config from session store
	sessionConf := a.sessionStore.GetConfig(ctx.Origin)

	// Convert domain messages to agent domain messages
	history := make([]*agentdomain.SessionMessage, len(ctx.History))
	for i, msg := range ctx.History {
		history[i] = a.toAgentMessage(msg)
	}

	message := a.toAgentMessage(ctx.IncomingMsg)

	respCtx := &ctxs.RespondCtx{
		Origin:       ctx.Origin,
		Channel:      ctx.Channel,
		ChatID:       ctx.ChatID,
		ChatName:     ctx.ChatName,
		SessionConf:  sessionConf,
		History:      history,
		Message:      message,
		IsAdmin:      ctx.IsAdmin,
		SessionStore: a.sessionStore,
		OnResponse:   ctx.OnResponse,
		OnReact:      ctx.OnReact,
		GetHistory: func() ([]*agentdomain.SessionMessage, error) {
			return a.sessionStore.GetHistory(ctx.Channel, ctx.ChatID, sessionConf.HistorySize)
		},
	}

	return a.agentCtl.Respond(respCtx)
}

func (a *AgentAdapter) toAgentMessage(msg *sessiondomain.Message) *agentdomain.SessionMessage {
	if msg == nil {
		return nil
	}
	return &agentdomain.SessionMessage{
		ID:         msg.ID,
		CreatedAt:  msg.CreatedAt,
		UpdatedAt:  msg.UpdatedAt,
		Channel:    msg.Channel,
		ChatID:     msg.ChatID,
		ChatName:   msg.ChatName,
		SenderID:   msg.SenderID,
		SenderName: msg.SenderName,
		MessageID:  msg.MessageID,
		Message:    msg.Message,
		Role:       msg.Role,
		Extra:      msg.Extra,
	}
}

// Ensure AgentAdapter implements sessiondomain.AgentService
var _ sessiondomain.AgentService = (*AgentAdapter)(nil)
