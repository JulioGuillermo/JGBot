package prompt

import (
	"JGBot/log"
	"JGBot/session/sessionconf/sc"
	"JGBot/skill"
	"fmt"
	"os"
	"strings"
)

func getSystemPromptFile(conf *sc.SessionConf) string {
	if conf.SystemPromptFile == "" {
		return DefaultSystemPrompt
	}

	data, err := os.ReadFile(conf.SystemPromptFile)
	if err == nil && len(data) > 0 {
		return string(data)
	}

	return DefaultSystemPrompt
}

func getSkillsPrompt(conf *sc.SessionConf) string {
	var sb strings.Builder
	sb.WriteString("## Available skills:\n")
	for _, skillConf := range conf.Skills {
		if !skillConf.Enabled {
			continue
		}

		skill, ok := skill.Skills[skillConf.Name]
		if !ok {
			log.Warn("Skill not found", "skill", skillConf.Name)
			continue
		}

		if skill.HasTool {
			fmt.Fprintf(&sb, "- %s: [Skill Tool available] %s\n", skill.Name, skill.Description)
		} else {
			fmt.Fprintf(&sb, "- %s: %s\n", skill.Name, skill.Description)
		}
	}
	return sb.String()
}

func GetSystemPrompt(conf *sc.SessionConf) string {
	prompt := getSystemPromptFile(conf)
	prompt = strings.TrimSpace(prompt)

	toolConf := conf.GetToolConf("skill")
	if toolConf != nil && toolConf.Enabled {
		prompt += "\n\n" + getSkillsPrompt(conf)
	}

	return prompt
}
