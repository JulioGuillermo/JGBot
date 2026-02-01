package session

import (
	"JGBot/agent"
	"JGBot/channels/channelctl"
	"JGBot/log"
	"JGBot/session/sessionconf"
	"JGBot/session/sessiondb"
	"fmt"
)

type SessionCtl struct {
	channelCtl *channelctl.ChannelCtl
	agent      *agent.Agent
	sessionCtl *sessionconf.SessionCtl
}

func NewSessionCtl(channelCtl *channelctl.ChannelCtl, agent *agent.Agent) (*SessionCtl, error) {
	sessiondb.Migrate()

	sessionCtl, err := sessionconf.NewSessionCtl()
	if err != nil {
		return nil, err
	}

	ctl := &SessionCtl{
		channelCtl: channelCtl,
		agent:      agent,
		sessionCtl: sessionCtl,
	}

	channelCtl.OnMessage(ctl.OnNewMessage)

	return ctl, nil
}

func (s *SessionCtl) OnNewMessage(channel string, origin string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string) {
	if message == "" {
		log.Info("Empty msg")
		return
	}

	sessionConf := s.sessionCtl.GetConfigOrigin(origin)
	if sessionConf == nil {
		log.Info("Not config session", "origin", origin)
		s.sessionCtl.AddUnconfig(chatName, fmt.Sprintf("%s:%d", channel, chatID), origin)
		return
	} else if !sessionConf.Allowed {
		log.Info("Session not allowed", "origin", origin)
		return
	}

	history, err := sessiondb.GetHistory(channel, chatID, 50)
	if err != nil {
		log.Error("Get history error", "err", err)
		return
	}
	msg, err := sessiondb.NewSessionMessage(channel, chatID, chatName, senderID, senderName, messageID, message, "user", "")
	if err != nil {
		log.Error("New session message error", "err", err)
		return
	}
	if msg == nil {
		log.Info("Empty msg")
		return
	}

	if !sessionConf.Respond.Respond(message) {
		return
	}

	s.agent.Respond(
		sessionConf,
		history,
		msg,
		func(text, role, extra string) error {
			sessiondb.NewSessionMessage(
				channel,
				chatID,
				chatName,
				senderID,
				role,
				messageID,
				text,
				role,
				extra,
			)
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
