package admin

import (
	"JGBot/agent/tools"
	"JGBot/ctxs"
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

func (c *AdminSendMessageInitializerConf) sendMessage(rCtx *ctxs.RespondCtx, args AdminSendMessageArgs) string {
	if args.Session == "" || args.Message == "" {
		return "Error: Session (Origin) and Message are required for send_message."
	}

	if rCtx.ChannelCtl == nil {
		return "Error: Channel controller not initialized."
	}

	if rCtx.SessionStore == nil {
		return "Error: Session store not initialized."
	}

	sessionConf := rCtx.SessionStore.GetConfig(args.Session)
	if sessionConf == nil {
		return fmt.Sprintf("Error: Session with origin %s not found.", args.Session)
	}

	parts := strings.Split(sessionConf.ID, ":")
	if len(parts) != 2 {
		return fmt.Sprintf("Error: Invalid session ID format for %s: %s.", args.Session, sessionConf.ID)
	}

	channel := parts[0]
	chatID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return fmt.Sprintf("Error: Invalid chat ID in session: %s.", err.Error())
	}

	channelObj, err := rCtx.ChannelCtl.GetChannel(channel)
	if err != nil {
		return fmt.Sprintf("Error: Channel %s not found: %s.", channel, err.Error())
	}

	err = channelObj.SendMessage(uint(chatID), args.Message)
	if err != nil {
		return fmt.Sprintf("Error: Fail to send message: %s.", err.Error())
	}

	return fmt.Sprintf("Success: Message sent to %s (%s).", sessionConf.Name, args.Session)
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

			return c.sendMessage(rCtx, args), nil
		},
	}
}
