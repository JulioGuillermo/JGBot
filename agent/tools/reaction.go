package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/callbacks"
)

type ReactionInput struct {
	MessageID uint   `json:"message_id"`
	Reaction  string `json:"reaction"`
}

type ReactionTool struct {
	CallbacksHandler callbacks.Handler
	OnReact          func(msg uint, reaction string) error
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
	input = GetInputString(input)

	var args ReactionInput
	if t.CallbacksHandler != nil {
		t.CallbacksHandler.HandleToolStart(ctx, input)
	}

	if err := json.Unmarshal([]byte(input), &args); err != nil {
		fmt.Println(input)
		if t.CallbacksHandler != nil {
			t.CallbacksHandler.HandleToolError(ctx, err)
		}
		return fmt.Sprintf("ERROR: Invalid JSON format: %s", err.Error()), nil
	}

	fmt.Println("### Reacting...", input, args.MessageID, args.Reaction)
	err := t.OnReact(args.MessageID, args.Reaction)
	if err != nil {
		if t.CallbacksHandler != nil {
			t.CallbacksHandler.HandleToolError(ctx, err)
		}
		return fmt.Sprintf("FAILURE: Could not apply reaction. Reason: %s.", err.Error()), nil
	}

	output := fmt.Sprintf("SUCCESS: Reaction %s applied to %d", args.Reaction, args.MessageID)
	if t.CallbacksHandler != nil {
		t.CallbacksHandler.HandleToolEnd(ctx, output)
	}

	return output, nil
}
