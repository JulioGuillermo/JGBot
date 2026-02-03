package tools_conf

import (
	"JGBot/agent/tools"
	"JGBot/ctxs"
)

type ToolInitializerConf interface {
	Name() string
	ToolInitializer(ctx *ctxs.RespondCtx) tools.Tool
}
