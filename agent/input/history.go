package input

import (
	toolargs "JGBot/agent/tool_args"
	"JGBot/agent/tools"
	"JGBot/session/sessiondb"
	"encoding/json"
	"errors"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

type HistoryInput struct {
}

func NewHistoryInput() *HistoryInput {
	return &HistoryInput{}
}

func (HistoryInput) FormatMessages(values map[string]any) ([]llms.ChatMessage, error) {
	hist, ok := values["ChatHistory"]
	if !ok {
		return nil, errors.New("HistoryInput: ChatHistory required")
	}
	data, ok := hist.(string)
	if !ok {
		return nil, errors.New("HistoryInput: ChatHistory must be a json string")
	}

	var history []*sessiondb.SessionMessage
	err := json.Unmarshal([]byte(data), &history)
	if err != nil {
		return nil, err
	}

	messages := make([]llms.ChatMessage, 0, len(history))
	var chatMsg llms.ChatMessage
	for _, msg := range history {
		switch strings.ToLower(msg.Role) {
		case "tool":
			toolResult := tools.ToolResultFromJson(msg.Extra)
			if toolResult == nil {
				continue
			}
			var content string
			if toolResult.Error != "" {
				content = toolResult.Error
			} else {
				content = toolResult.Output
			}
			chatMsg = llms.ToolChatMessage{
				ID:      toolResult.ToolCall.ID,
				Content: content,
			}
		case "system":
			chatMsg = llms.SystemChatMessage{
				Content: msg.Message,
			}
		case "assistant":
			m := llms.AIChatMessage{
				Content: msg.Message,
			}

			toolCall := tools.ToolCallFromJson(msg.Extra)
			if toolCall != nil {
				content := toolargs.GetMsgValidArg(toolCall.Input)
				m.ToolCalls = []llms.ToolCall{{
					ID:   toolCall.ID,
					Type: "function",
					FunctionCall: &llms.FunctionCall{
						Name:      toolCall.Tool,
						Arguments: content,
					},
				}}
			}

			chatMsg = m
		default: // user
			chatMsg = llms.HumanChatMessage{
				Content: msg.String(),
			}
		}

		if chatMsg != nil {
			messages = append(messages, chatMsg)
		}
	}

	return messages, nil
}

func (HistoryInput) GetInputVariables() []string {
	return []string{
		"ChatHistory",
	}
}
