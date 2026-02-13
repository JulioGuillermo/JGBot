package ctxs

import (
	"JGBot/channels/channelctl"
	"JGBot/session/sessionconf"
	"JGBot/session/sessionconf/sc"
	"JGBot/session/sessiondb"
	"fmt"
	"strconv"
	"strings"
)

type OnResponse func(text, role, extra string) error
type OnReact func(msg uint, reaction string) error
type GetHistory func() ([]*sessiondb.SessionMessage, error)

type RespondCtx struct {
	Origin      string
	Channel     string
	ChatID      uint
	ChatName    string
	SessionConf *sc.SessionConf
	History     []*sessiondb.SessionMessage
	Message     *sessiondb.SessionMessage
	IsAdmin     bool
	SessionCtl  *sessionconf.SessionCtl
	ChannelCtl  *channelctl.ChannelCtl
	OnResponse  OnResponse
	OnReact     OnReact
	GetHistory  GetHistory
}

func (c *RespondCtx) Copy() *RespondCtx {
	history := make([]*sessiondb.SessionMessage, len(c.History))
	copy(history, c.History)
	return &RespondCtx{
		SessionConf: c.SessionConf,
		History:     history,
		Message:     c.Message,
		IsAdmin:     c.IsAdmin,
		SessionCtl:  c.SessionCtl,
		ChannelCtl:  c.ChannelCtl,
		OnResponse:  c.OnResponse,
		OnReact:     c.OnReact,
	}
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

	sessionConf := c.SessionCtl.GetConfigOrigin(origin)
	if sessionConf == nil {
		return nil, fmt.Errorf("Error: Session with origin %s not found.", origin)
	}

	Partses := strings.Split(sessionConf.ID, ":")
	if len(Partses) != 2 {
		return nil, fmt.Errorf("Error: Invalid target session ID format: %s.", sessionConf.ID)
	}

	chatID, err := strconv.ParseUint(Partses[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Error: Invalid target chat ID: %s.", err.Error())
	}

	ctx := c.Copy()
	ctx.Origin = origin
	ctx.SessionConf = sessionConf
	ctx.Channel = Partses[0]
	ctx.ChatID = uint(chatID)
	ctx.ChatName = sessionConf.Name
	return ctx, nil
}
