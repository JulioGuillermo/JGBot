package admin

import (
	"JGBot/agent/tools"
	"JGBot/ctxs"
	"JGBot/session/sessionconf"
	"context"
	"fmt"
	"strings"
)

type AdminListSessionsArgs struct{}

type AdminListSessionsInitializerConf struct{}

func NewAdminListSessionsInitializerConf() *AdminListSessionsInitializerConf {
	return &AdminListSessionsInitializerConf{}
}

func (c *AdminListSessionsInitializerConf) Name() string {
	return "list_sessions"
}

func (c *AdminListSessionsInitializerConf) listSessions(sCtl *sessionconf.SessionCtl) string {
	sessions := sCtl.Config.Sessions
	if len(sessions) == 0 {
		return "The list of sessions is empty."
	}

	var sb strings.Builder
	sb.WriteString("Configured Sessions:\n")
	for _, s := range sessions {
		ch := "Unknown"
		if strings.Contains(s.Origin, ":") {
			ch = strings.Split(s.Origin, ":")[0]
		}
		fmt.Fprintf(&sb, "- Name: %s, Origin: %s, Channel: %s\n", s.Name, s.Origin, ch)
	}
	return sb.String()
}

func (c *AdminListSessionsInitializerConf) ToolInitializer(rCtx *ctxs.RespondCtx) tools.Tool {
	return &tools.ToolAutoArgs[AdminListSessionsArgs]{
		ToolName:        c.Name(),
		ToolDescription: "Administrative tools to list sessions.",
		IsAdmin:         rCtx.IsAdmin,
		ToolFunc: func(ctx context.Context, args AdminListSessionsArgs) (string, error) {
			if !rCtx.IsAdmin {
				return "Error: You do not have administrator permissions.", nil
			}

			return c.listSessions(rCtx.SessionCtl), nil
		},
	}
}
