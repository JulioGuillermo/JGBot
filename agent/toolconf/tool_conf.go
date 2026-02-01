package toolconf

import (
	"JGBot/agent/tools"
	"JGBot/session/sessionconf/sc"
	"JGBot/session/sessiondb"
)

type ToolInitializer func(sessionConf *sc.SessionConf, history []*sessiondb.SessionMessage, message *sessiondb.SessionMessage, onResponse func(text, role, extra string) error, onReact func(msg uint, reaction string) error) tools.Tool

type ToolInitializerConf interface {
	Name() string
	ToolInitializer(sessionConf *sc.SessionConf, history []*sessiondb.SessionMessage, message *sessiondb.SessionMessage, onResponse func(text, role, extra string) error, onReact func(msg uint, reaction string) error) tools.Tool
}
