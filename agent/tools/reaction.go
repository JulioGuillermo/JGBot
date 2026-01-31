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
	onReact  func(msg uint, reaction string) error
	onResult func(toolResult ToolResult)
}

func NewReactionTool(onReact func(msg uint, reaction string) error, onResult func(toolResult ToolResult)) tools.Tool {
	return ReactionTool{
		onReact:  onReact,
		onResult: onResult,
	}
}

func (t ReactionTool) Name() string {
	return "message_reaction"
}

func (t ReactionTool) Description() string {
	return "Use this tool to add or update a reaction to a specific message. " +
		"The input must be a valid JSON object containing: " +
		"\"message_id\" int (the identifier of the message) and " +
		"\"reaction\" string (the emoji representing the reaction, e.g., '👍', '👎')."
}

func (t ReactionTool) Call(ctx context.Context, input string) (string, error) {
	var args ReactionInput

	if err := json.Unmarshal([]byte(input), &args); err != nil {
		errMsg := fmt.Sprintf("ERROR: Invalid JSON format: %s", err.Error())
		t.onResult(ToolResult{
			Error: errMsg,
		})
		return errMsg, nil
	}

	err := t.onReact(args.MessageID, args.Reaction)
	if err != nil {
		errMsg := fmt.Sprintf("FAILURE: Could not apply reaction. Reason: %s.", err.Error())
		t.onResult(ToolResult{
			Error: errMsg,
		})
		return errMsg, nil
	}

	output := fmt.Sprintf("SUCCESS: Reaction %s applied to %d", args.Reaction, args.MessageID)
	t.onResult(ToolResult{
		Output: output,
	})

	return output, nil
}
