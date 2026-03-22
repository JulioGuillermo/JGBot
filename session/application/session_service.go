package sessionapplication

import (
	"JGBot/log"
	sessiondomain "JGBot/session/domain"
	"errors"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/agents"
)

type SessionService struct {
	repo         sessiondomain.MessageRepository
	configRepo   sessiondomain.ConfigurationRepository
	agentService sessiondomain.AgentService
	channelSvc   sessiondomain.ChannelService
}

func NewSessionService(
	repo sessiondomain.MessageRepository,
	configRepo sessiondomain.ConfigurationRepository,
	agentService sessiondomain.AgentService,
	channelSvc sessiondomain.ChannelService,
) *SessionService {
	return &SessionService{
		repo:         repo,
		configRepo:   configRepo,
		agentService: agentService,
		channelSvc:   channelSvc,
	}
}

func (s *SessionService) OnNewMessage(channel string, origin string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string) {
	if message == "" {
		log.Info("Empty msg")
		return
	}

	autoEnable, err := s.channelSvc.GetChannelAutoEnable(channel)
	if err != nil {
		log.Warn("Failed to get channel auto-enable status", "err", err)
		return
	}

	sessionConf, allowed := s.getSessionConf(origin, chatName, channel, chatID, autoEnable)
	if !allowed {
		return
	}

	if message == "/reset!" {
		err := s.repo.ClearHistory(channel, chatID)
		if err != nil {
			log.Error("Clear history error", "err", err)
			s.channelSvc.SendMessage(channel, chatID, fmt.Sprintf("Fail to clear history: %s", err.Error()))
		} else {
			s.channelSvc.SendMessage(channel, chatID, "History cleared")
		}
		return
	}

	history, err := s.repo.GetHistory(channel, chatID, sessionConf.HistorySize)
	if err != nil {
		log.Error("Get history error", "err", err)
		return
	}

	isAdmin, cleanMessage := s.isAdmin(sessionConf.Admin, message)
	if isAdmin {
		log.Info("Using ADMIN permissions")
	}

	msg := sessiondomain.NewMessage(channel, chatID, chatName, senderID, senderName, messageID, cleanMessage, "user", "")
	err = s.repo.Save(msg)
	if err != nil {
		log.Error("New session message error", "err", err)
		return
	}

	if !sessionConf.ShouldRespond(cleanMessage) {
		return
	}

	s.respond(
		sessionConf,
		channel,
		chatID,
		chatName,
		senderID,
		senderName,
		messageID,
		history,
		msg,
		isAdmin,
	)
}

// OnAutoActivation handles automatic task activations (cron, timer)
func (s *SessionService) OnAutoActivation(ctx *sessiondomain.ActivationContext) {
	msgContent := fmt.Sprintf("CRON EXECUTION: %s\n\nSCHEDULE: %s\n\nDESCRIPTION: %s\n\nMESSAGE: %s", ctx.Name, ctx.Schedule, ctx.Description, ctx.Message)

	autoEnable, err := s.channelSvc.GetChannelAutoEnable(ctx.Channel)
	if err != nil {
		log.Warn("Failed to get channel auto-enable status", "err", err)
		return
	}

	sessionConf, allowed := s.getSessionConf(ctx.Origin, ctx.ChatName, ctx.Channel, ctx.ChatID, autoEnable)
	if !allowed {
		return
	}

	history, err := s.repo.GetHistory(ctx.Channel, ctx.ChatID, sessionConf.HistorySize)
	if err != nil {
		log.Error("Get history error", "err", err)
		return
	}

	msg := sessiondomain.NewMessage(ctx.Channel, ctx.ChatID, ctx.ChatName, ctx.SenderID, "tool", ctx.MessageID, msgContent, "tool", "")
	err = s.repo.Save(msg)
	if err != nil {
		log.Error("New session message error", "err", err)
		return
	}

	s.respond(
		sessionConf,
		ctx.Channel,
		ctx.ChatID,
		ctx.ChatName,
		ctx.SenderID,
		"",
		ctx.MessageID,
		history,
		msg,
		false,
	)
}

func (s *SessionService) isAdmin(adminPermission string, message string) (bool, string) {
	admin := false
	if strings.HasPrefix(message, "/admin ") {
		admin = true
		message = strings.TrimPrefix(message, "/admin ")
	} else if message == "/admin" {
		admin = true
		message = ""
	}

	if adminPermission == "full" {
		return true, message
	}

	if adminPermission == "allow" {
		return admin, message
	}

	return false, message
}

func (s *SessionService) getSessionConf(origin, chatName, channel string, chatID uint, autoEnable bool) (*sessiondomain.SessionConfiguration, bool) {
	sessionConf := s.configRepo.GetConfig(origin)

	if sessionConf == nil && !autoEnable {
		log.Info("Not config session", "origin", origin)
		s.configRepo.CreateUnconfig(chatName, fmt.Sprintf("%s:%d", channel, chatID), origin, channel)
		return nil, false
	}

	if sessionConf == nil {
		log.Warn("Auto enable session", "origin", origin)
		sessionConf = s.configRepo.CreateConfig(chatName, fmt.Sprintf("%s:%d", channel, chatID), origin, channel)
	}

	if !sessionConf.Allowed {
		log.Info("Session not allowed", "origin", origin)
	}

	return sessionConf, sessionConf.Allowed
}

func (s *SessionService) respond(
	sessionConf *sessiondomain.SessionConfiguration,
	channel string,
	chatID uint,
	chatName string,
	senderID uint,
	senderName string,
	messageID uint,
	history []*sessiondomain.Message,
	msg *sessiondomain.Message,
	isAdmin bool,
) {
	onResponse := func(text, role, extra string) error {
		if senderName == "" {
			senderName = role
		}
		respMsg := sessiondomain.NewMessage(
			channel,
			chatID,
			chatName,
			senderID,
			senderName,
			messageID,
			text,
			role,
			extra,
		)
		s.repo.Save(respMsg)
		if text == "" {
			return nil
		}
		text, thinkings := sessiondomain.ExtractReasoning(text)
		if sessionConf.ShowThink && thinkings != "" {
			text = fmt.Sprintf(
				"```thinking\n%s\n```\n\n%s",
				thinkings,
				text,
			)
		}
		return s.channelSvc.SendMessage(channel, chatID, text)
	}

	onReact := func(msgID uint, reaction string) error {
		return s.channelSvc.SendReaction(
			channel,
			chatID,
			msgID,
			reaction,
		)
	}

	s.channelSvc.SendStatus(channel, chatID, "writing")
	defer s.channelSvc.SendStatus(channel, chatID, "normal")

	err := s.agentService.Respond(&sessiondomain.MessageContext{
		Origin:      sessionConf.Origin,
		Channel:     channel,
		ChatID:      chatID,
		ChatName:    chatName,
		Config:      sessionConf,
		History:     history,
		IncomingMsg: msg,
		IsAdmin:     isAdmin,
		OnResponse:  onResponse,
		OnReact:     onReact,
	})

	if err == nil {
		return
	}

	if errors.Is(err, agents.ErrNotFinished) {
		log.Error("Agent Max Iter Error", "err", err)
		s.channelSvc.SendMessage(channel, chatID, "[MAX ITER] I need a break 🤒...")
		return
	}
	log.Error("Agent respond error", "err", err)
	s.channelSvc.SendMessage(channel, chatID, "[ERROR] I probably made a mistake 😵\u200d💫...")
}
