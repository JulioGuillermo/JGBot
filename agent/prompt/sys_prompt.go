package prompt

import (
	"JGBot/session/sessionconf/sc"
	"os"
)

func GetSystemPrompt(conf *sc.SessionConf) string {
	if conf.SystemPromptFile == "" {
		return DefaultSystemPrompt
	}

	data, err := os.ReadFile(conf.SystemPromptFile)
	if err == nil && len(data) > 0 {
		return string(data)
	}

	return DefaultSystemPrompt
}
