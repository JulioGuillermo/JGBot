package prompt

import (
	"JGBot/agent/toolconf/skill"
	"JGBot/agent/toolconf/tools_conf"
	"JGBot/session/sessionconf/sc"
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

func JoinPromptWithSkill(prompt string, conf *sc.SessionConf) string {
	prompt = strings.TrimSpace(prompt)

	prompt += "\n\n" + GetToolPrompt(conf)

	toolConf := conf.GetToolConf("skill")
	if toolConf != nil && toolConf.Enabled {
		prompt += "\n\n" + skill.GetSkillsPrompt(conf)
	}

	return prompt
}

func GetSystemPrompt(conf *sc.SessionConf) string {
	return JoinPromptWithSkill(getSystemPromptFile(conf), conf)
}

func GetSubAgentPrompt(conf *sc.SessionConf) string {
	return JoinPromptWithSkill(SubAgentPrompt, conf)
}

func GetToolPrompt(conf *sc.SessionConf) string {
	var sb strings.Builder
	sb.WriteString("**Available tools:**\n")
	sb.WriteString("You can use the following tools by a direct call:\n")
	for _, toolConf := range tools_conf.NativeTools {
		confTool := conf.GetToolConf(toolConf.Name())
		if confTool == nil || !confTool.Enabled {
			continue
		}

		fmt.Fprintf(&sb, "- %s\n", toolConf.Name())
	}
	if subAgentConf := conf.GetToolConf("subagent"); subAgentConf != nil && subAgentConf.Enabled {
		fmt.Fprintf(&sb, "- subagent\n")
	}

	return sb.String()
}
