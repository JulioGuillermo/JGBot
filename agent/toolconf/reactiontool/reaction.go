package reactiontool

import (
	"JGBot/agent/tools"
	"JGBot/session/sessionconf/sc"
	"JGBot/session/sessiondb"
	"context"
	"fmt"
)

type ReactionArgs struct {
	MessageID uint   `json:"message_id" description:"The identifier of the message to react to."`
	Reaction  string `json:"reaction" description:"The emoji representing the reaction, e.g., '👍', '👎'."`
}

type ReactionInitializerConf struct{}

func (c *ReactionInitializerConf) Name() string {
	return "message_reaction"
}

func (c *ReactionInitializerConf) ToolInitializer(sessionConf *sc.SessionConf, history []*sessiondb.SessionMessage, message *sessiondb.SessionMessage, onResponse func(text, role, extra string) error, onReact func(msg uint, reaction string) error) tools.Tool {
	return &tools.ToolAutoArgs[ReactionArgs]{
		ToolName:        c.Name(),
		ToolDescription: "Use this tool to add or update a reaction to a specific message.",
		ToolFunc: func(ctx context.Context, args ReactionArgs) (string, error) {
			err := onReact(args.MessageID, args.Reaction)
			if err != nil {
				return "", fmt.Errorf("FAILURE: Could not apply reaction. Reason: %s.", err.Error())
			}

			return fmt.Sprintf("SUCCESS: Reaction %s applied to %d", args.Reaction, args.MessageID), nil
		},
	}
}
