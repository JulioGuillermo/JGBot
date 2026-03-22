package ctxs

import (
	agentdomain "JGBot/agent/domain"
	channelsdomain "JGBot/channels/domain"
	"fmt"
	"strconv"
	"strings"
)

type OnResponse func(text, role, extra string) error
type OnReact func(msg uint, reaction string) error
type GetHistory func() ([]*agentdomain.SessionMessage, error)

type RespondCtx struct {
	Origin       string
	Channel      string
	ChatID       uint
	ChatName     string
	SessionConf  *agentdomain.SessionConfig
	History      []*agentdomain.SessionMessage
	Message      *agentdomain.SessionMessage
	IsAdmin      bool
	SessionStore agentdomain.SessionStore
	ChannelCtl   channelsdomain.ChannelController
	OnResponse   OnResponse
	OnReact      OnReact
	GetHistory   GetHistory
}

func (c *RespondCtx) Copy() *RespondCtx {
	history := make([]*agentdomain.SessionMessage, len(c.History))
	copy(history, c.History)
	return &RespondCtx{
		Origin:       c.Origin,
		Channel:      c.Channel,
		ChatID:       c.ChatID,
		ChatName:     c.ChatName,
		SessionConf:  c.SessionConf,
		History:      history,
		Message:      c.Message,
		IsAdmin:      c.IsAdmin,
		SessionStore: c.SessionStore,
		ChannelCtl:   c.ChannelCtl,
		OnResponse:   c.OnResponse,
		OnReact:      c.OnReact,
		GetHistory:   c.GetHistory,
	}
}

func (c *RespondCtx) Status(status channelsdomain.Status) {
	if c.ChannelCtl == nil {
		return
	}
	channel, _ := c.ChannelCtl.GetChannel(c.Channel)
	if channel == nil {
		return
	}
	channel.SendStatus(c.ChatID, status)
}

func (c *RespondCtx) GetOrigin(origin string) string {
	if c.IsAdmin && origin != "" {
		return origin
	}
	return c.Origin
}

func (c *RespondCtx) GetSessionCtx(origin string) (*RespondCtx, error) {
	if !c.IsAdmin || origin == "" {
		return c.Copy(), nil
	}

	sessionConf := c.SessionStore.GetConfig(origin)
	if sessionConf == nil {
		return nil, fmt.Errorf("Error: Session with origin %s not found.", origin)
	}

	parts := strings.Split(sessionConf.ID, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Error: Invalid target session ID format: %s.", sessionConf.ID)
	}

	chatID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Error: Invalid target chat ID: %s.", err.Error())
	}

	ctx := c.Copy()
	ctx.Origin = origin
	ctx.SessionConf = sessionConf
	ctx.Channel = parts[0]
	ctx.ChatID = uint(chatID)
	ctx.ChatName = sessionConf.Name
	return ctx, nil
}
