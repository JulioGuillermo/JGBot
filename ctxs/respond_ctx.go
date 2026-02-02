package ctxs

import (
	"JGBot/session/sessionconf/sc"
	"JGBot/session/sessiondb"
)

type OnResponse func(text, role, extra string) error
type OnReact func(msg uint, reaction string) error

type RespondCtx struct {
	SessionConf *sc.SessionConf
	History     []*sessiondb.SessionMessage
	Message     *sessiondb.SessionMessage
	OnResponse  OnResponse
	OnReact     OnReact
}
