package session

import (
	"JGBot/agent"
	"JGBot/channels/channelctl"
	"JGBot/session/sessiondb"
	"log/slog"
)

type SessionCtl struct {
	logger     *slog.Logger
	channelCtl *channelctl.ChannelCtl
	agent      *agent.Agent
}

func NewSessionCtl(logger *slog.Logger, channelCtl *channelctl.ChannelCtl, agent *agent.Agent) (*SessionCtl, error) {
	sessiondb.Migrate()

	ctl := &SessionCtl{
		logger:     logger,
		channelCtl: channelCtl,
		agent:      agent,
	}

	channelCtl.OnMessage(ctl.OnNewMessage)

	return ctl, nil
}

func (s *SessionCtl) OnNewMessage(channel string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string) {
	if message == "" {
		s.logger.Info("Empty msg")
		return
	}

	history, err := sessiondb.GetHistory(channel, chatID, 50)
	if err != nil {
		s.logger.Error("Get history error", "err", err)
		return
	}
	msg, err := sessiondb.NewSessionMessage(channel, chatID, chatName, senderID, senderName, messageID, message, "user", "")
	if err != nil {
		s.logger.Error("New session message error", "err", err)
		return
	}
	if msg == nil {
		s.logger.Info("Empty msg")
		return
	}

	s.agent.Respond(
		history,
		msg,
		func(text, role, extra string) error {
			sessiondb.NewSessionMessage(channel, chatID, chatName, senderID, role, messageID, text, role, extra)
			if text == "" {
				return nil
			}
			return s.channelCtl.SendMessage(channel, chatID, text)
		},
		func(msg uint, reaction string) error {
			return s.channelCtl.ReactMessage(channel, chatID, msg, reaction)
		},
	)
}
