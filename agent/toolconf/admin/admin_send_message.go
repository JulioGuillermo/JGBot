package admin

import (
	"JGBot/agent/tools"
	"JGBot/channels/channelctl"
	"JGBot/ctxs"
	"JGBot/session/sessionconf"
	"context"
	"fmt"
	"strconv"
	"strings"
)

type AdminSendMessageArgs struct {
	Session string `json:"session" description:"The origin of the target session."`
	Message string `json:"message" description:"The message to send to the target session."`
}

type AdminSendMessageInitializerConf struct{}

func NewAdminSendMessageInitializerConf() *AdminSendMessageInitializerConf {
	return &AdminSendMessageInitializerConf{}
}

func (c *AdminSendMessageInitializerConf) Name() string {
	return "send_message"
}

func (c *AdminSendMessageInitializerConf) sendMessage(sCtl *sessionconf.SessionCtl, cCtl *channelctl.ChannelCtl, origin, message string) string {
	if origin == "" || message == "" {
		return "Error: Session (Origin) and Message are required for send_message."
	}

	if cCtl == nil {
		return "Error: Channel controller not initialized."
	}

	if sCtl == nil {
		return "Error: Session controller not initialized."
	}

	sessionConf := sCtl.GetConfigOrigin(origin)
	if sessionConf == nil {
		return fmt.Sprintf("Error: Session with origin %s not found.", origin)
	}

	parts := strings.Split(sessionConf.ID, ":")
	if len(parts) != 2 {
		return fmt.Sprintf("Error: Invalid session ID format for %s: %s.", origin, sessionConf.ID)
	}

	channel := parts[0]
	chatID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return fmt.Sprintf("Error: Invalid chat ID in session: %s.", err.Error())
	}

	err = cCtl.SendMessage(channel, uint(chatID), message)
	if err != nil {
		return fmt.Sprintf("Error: Fail to send message: %s.", err.Error())
	}

	return fmt.Sprintf("Success: Message sent to %s (%s).", sessionConf.Name, origin)
}

func (c *AdminSendMessageInitializerConf) ToolInitializer(rCtx *ctxs.RespondCtx) tools.Tool {
	return &tools.ToolAutoArgs[AdminSendMessageArgs]{
		ToolName:        c.Name(),
		ToolDescription: "Administrative tools to send messages across sections.",
		IsAdmin:         rCtx.IsAdmin,
		ToolFunc: func(ctx context.Context, args AdminSendMessageArgs) (string, error) {
			if !rCtx.IsAdmin {
				return "Error: You do not have administrator permissions.", nil
			}

			return c.sendMessage(rCtx.SessionCtl, rCtx.ChannelCtl, args.Session, args.Message), nil
		},
	}
}
