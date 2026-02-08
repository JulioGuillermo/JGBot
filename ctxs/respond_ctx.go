package ctxs

import (
	"JGBot/session/sessionconf/sc"
	"JGBot/session/sessiondb"
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
		OnResponse:  c.OnResponse,
		OnReact:     c.OnReact,
	}
}
