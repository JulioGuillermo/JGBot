package session

import (
	"JGBot/channels"
	"JGBot/ctxs"
	"JGBot/log"
	"JGBot/session/sessiondb"
	"errors"
	"fmt"

	"github.com/tmc/langchaingo/agents"
)

func (s *SessionCtl) OnAutoActivation(
	origin string,
	channel string,
	chatID uint,
	chatName string,
	senderID uint,
	messageID uint,

	name,
	schedule,
	description,
	message string,
) {
	msgContent := fmt.Sprintf("CRON EXECUTION: %s\n\nSCHEDULE: %s\n\nDESCRIPTION: %s\n\nMESSAGE: %s", name, schedule, description, message)

	sessionConf := s.sessionCtl.GetConfigOrigin(origin)
	if sessionConf == nil {
		if s.channelCtl.AutoEnableSession(channel) {
			log.Warn("Auto enable session", "origin", origin)
			s.sessionCtl.AddConfig(chatName, fmt.Sprintf("%s:%d", channel, chatID), origin)
		} else {
			log.Info("Not config session", "origin", origin)
			s.sessionCtl.AddUnconfig(chatName, fmt.Sprintf("%s:%d", channel, chatID), origin)
		}
		return
	} else if !sessionConf.Allowed {
		log.Info("Session not allowed", "origin", origin)
		return
	}

	history, err := sessiondb.GetHistory(channel, chatID, sessionConf.HistorySize)
	if err != nil {
		log.Error("Get history error", "err", err)
		return
	}
	msg, err := sessiondb.NewSessionMessage(channel, chatID, chatName, senderID, "tool", messageID, msgContent, "tool", "")
	if err != nil {
		log.Error("New session message error", "err", err)
		return
	}
	if msg == nil {
		log.Info("Empty msg")
		return
	}

	respCtx := &ctxs.RespondCtx{
		Origin:      origin,
		Channel:     channel,
		ChatID:      chatID,
		ChatName:    chatName,
		SessionConf: sessionConf,
		History:     history,
		Message:     msg,
		OnResponse: func(text, role, extra string) error {
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
		OnReact: func(msg uint, reaction string) error {
			return s.channelCtl.ReactMessage(channel, chatID, msg, reaction)
		},
		GetHistory: func() ([]*sessiondb.SessionMessage, error) {
			return sessiondb.GetHistory(channel, chatID, sessionConf.HistorySize)
		},
	}

	s.channelCtl.Status(channel, chatID, channels.Writing)
	defer s.channelCtl.Status(channel, chatID, channels.Normal)

	err = s.agent.Respond(respCtx)
	if err == nil {
		return
	}
	if errors.Is(err, agents.ErrNotFinished) {
		log.Error("(AUTO ACT) Agent Max Iter Error", "err", err)
		s.channelCtl.SendMessage(channel, chatID, "(AUTO ACT) [MAX ITER] I need a break 🤒...")
		return
	}
	log.Error("Agent respond error", "err", err)
	s.channelCtl.SendMessage(channel, chatID, "(AUTO ACT) [ERROR] I probably made a mistake 😵‍💫...")
}
