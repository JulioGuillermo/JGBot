package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/tools"
)

type ReactionInput struct {
	MessageID uint   `json:"message_id"`
	Reaction  string `json:"reaction"`
}

type ReactionTool struct {
	onReact func(msg uint, reaction string) error
}

func NewReactionTool(onReact func(msg uint, reaction string) error) tools.Tool {
	return ReactionTool{
		onReact: onReact,
	}
}

func (t ReactionTool) Name() string {
	return "message_reaction"
}

func (t ReactionTool) Description() string {
	return "Use this tool to add or update a reaction to a specific message. " +
		"The input must be a valid JSON object containing: " +
		"\"message_id\" int (the identifier of the message) and " +
		"\"reaction\" string (the emoji representing the reaction, e.g., '👍', '👎', '😄')."
}

func (t ReactionTool) Call(ctx context.Context, input string) (string, error) {
	var args ReactionInput

	if err := json.Unmarshal([]byte(input), &args); err != nil {
		return "ERROR: Invalid JSON format. Please provide 'message_id' and 'reaction' as keys.", nil
	}

	err := t.onReact(args.MessageID, args.Reaction)
	if err != nil {
		return fmt.Sprintf("FAILURE: Could not apply reaction. Reason: %v. Please inform the user or try a different Message ID.", err), nil
	}

	return fmt.Sprintf("SUCCESS: Reaction %s applied to %d", args.Reaction, args.MessageID), nil
}
