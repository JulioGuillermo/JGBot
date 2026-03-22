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

	// 1. Check auto-enable (if needed by config logic) or just get channel
	// In old code: channelObj, err := s.channelCtl.GetChannel(channel)
	// New code: config logic needs to know if allowed.
	// We might need to query channel service for auto-enable flag.
	autoEnable, err := s.channelSvc.GetChannelAutoEnable(channel)
	if err != nil {
		log.Warn("Failed to get channel auto-enable status", "err", err)
		return
	}

	// 2. Get Config
	sessionConf := s.configRepo.GetConfig(origin)
	if sessionConf == nil && autoEnable {
		log.Warn("Auto enable session", "origin", origin)
		sessionConf = s.configRepo.CreateConfig(chatName, fmt.Sprintf("%s:%d", channel, chatID), origin, channel)
	} else if sessionConf == nil {
		log.Info("Not config session", "origin", origin)
		s.configRepo.CreateUnconfig(chatName, fmt.Sprintf("%s:%d", channel, chatID), origin, channel)
		return
	}

	if !sessionConf.Allowed {
		log.Info("Session not allowed", "origin", origin)
		return
	}

	// 3. Handle /reset!
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

	// 4. Get History
	history, err := s.repo.GetHistory(channel, chatID, sessionConf.HistorySize)
	if err != nil {
		log.Error("Get history error", "err", err)
		return
	}

	// 5. Admin Check
	isAdmin, cleanMessage := s.isAdmin(sessionConf.Admin, message)
	if isAdmin {
		log.Info("Using ADMIN permissions")
	}

	// 6. Save Incoming Message
	msg := sessiondomain.NewMessage(channel, chatID, chatName, senderID, senderName, messageID, cleanMessage, "user", "")
	err = s.repo.Save(msg)
	if err != nil {
		log.Error("New session message error", "err", err)
		return
	}

	// 7. Check Respond Policy
	if !sessionConf.ShouldRespond(cleanMessage) {
		return
	}

	// 8. Prepare Callbacks
	onResponse := func(text, role, extra string) error {
		respMsg := sessiondomain.NewMessage(channel, chatID, chatName, senderID, senderName, messageID, text, role, extra)
		s.repo.Save(respMsg)
		if text == "" {
			return nil
		}
		return s.channelSvc.SendMessage(channel, chatID, text)
	}

	onReact := func(msgID uint, reaction string) error {
		return s.channelSvc.SendReaction(channel, chatID, msgID, reaction)
	}

	// 9. Send "Writing" status
	s.channelSvc.SendStatus(channel, chatID, "writing")
	defer s.channelSvc.SendStatus(channel, chatID, "normal") // Assuming "normal" creates "stop writing" effect or idle

	// 10. Call Agent
	err = s.agentService.Respond(&sessiondomain.MessageContext{
		Origin:      origin,
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
	s.channelSvc.SendMessage(channel, chatID, "[ERROR] I probably made a mistake 😵‍💫...")
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

// OnAutoActivation handles automatic task activations (cron, timer)
func (s *SessionService) OnAutoActivation(ctx *sessiondomain.ActivationContext) {
	msgContent := fmt.Sprintf("CRON EXECUTION: %s\n\nSCHEDULE: %s\n\nDESCRIPTION: %s\n\nMESSAGE: %s", ctx.Name, ctx.Schedule, ctx.Description, ctx.Message)

	autoEnable, err := s.channelSvc.GetChannelAutoEnable(ctx.Channel)
	if err != nil {
		log.Warn("Failed to get channel auto-enable status", "err", err)
		return
	}

	sessionConf := s.configRepo.GetConfig(ctx.Origin)
	if sessionConf == nil && autoEnable {
		log.Warn("Auto enable session", "origin", ctx.Origin)
		sessionConf = s.configRepo.CreateConfig(ctx.ChatName, fmt.Sprintf("%s:%d", ctx.Channel, ctx.ChatID), ctx.Origin, ctx.Channel)
	} else if sessionConf == nil {
		log.Info("Not config session", "origin", ctx.Origin)
		s.configRepo.CreateUnconfig(ctx.ChatName, fmt.Sprintf("%s:%d", ctx.Channel, ctx.ChatID), ctx.Origin, ctx.Channel)
		return
	}
	if !sessionConf.Allowed {
		log.Info("Session not allowed", "origin", ctx.Origin)
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

	onResponse := func(text, role, extra string) error {
		respMsg := sessiondomain.NewMessage(ctx.Channel, ctx.ChatID, ctx.ChatName, ctx.SenderID, role, ctx.MessageID, text, role, extra)
		s.repo.Save(respMsg)
		if text == "" {
			return nil
		}
		return s.channelSvc.SendMessage(ctx.Channel, ctx.ChatID, text)
	}

	onReact := func(msgID uint, reaction string) error {
		return s.channelSvc.SendReaction(ctx.Channel, ctx.ChatID, msgID, reaction)
	}

	s.channelSvc.SendStatus(ctx.Channel, ctx.ChatID, "writing")
	defer s.channelSvc.SendStatus(ctx.Channel, ctx.ChatID, "normal")

	err = s.agentService.Respond(&sessiondomain.MessageContext{
		Origin:      ctx.Origin,
		Channel:     ctx.Channel,
		ChatID:      ctx.ChatID,
		ChatName:    ctx.ChatName,
		Config:      sessionConf,
		History:     history,
		IncomingMsg: msg,
		IsAdmin:     false,
		OnResponse:  onResponse,
		OnReact:     onReact,
	})

	if err == nil {
		return
	}

	if errors.Is(err, agents.ErrNotFinished) {
		log.Error("(AUTO ACT) Agent Max Iter Error", "err", err)
		s.channelSvc.SendMessage(ctx.Channel, ctx.ChatID, "(AUTO ACT) [MAX ITER] I need a break 🤒...")
		return
	}
	log.Error("Agent respond error", "err", err)
	s.channelSvc.SendMessage(ctx.Channel, ctx.ChatID, "(AUTO ACT) [ERROR] I probably made a mistake 😵‍💫...")
}
