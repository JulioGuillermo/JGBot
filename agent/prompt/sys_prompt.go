package prompt

import (
	agentdomain "JGBot/agent/domain"
	"JGBot/agent/toolconf/skill"
	"JGBot/agent/toolconf/tools_conf"
	"fmt"
	"os"
	"strings"
)

func getSystemPromptFile(conf *agentdomain.SessionConfig) string {
	if conf == nil || conf.SystemPromptFile == "" {
		return DefaultSystemPrompt
	}

	data, err := os.ReadFile(conf.SystemPromptFile)
	if err == nil && len(data) > 0 {
		return string(data)
	}

	return DefaultSystemPrompt
}

func JoinPromptWithSkill(prompt string, conf *agentdomain.SessionConfig) string {
	prompt = strings.TrimSpace(prompt)

	prompt += "\n\n" + GetToolPrompt(conf)

	if conf != nil {
		toolConf := getToolConf(conf, "skill")
		if toolConf != nil && toolConf.Enabled {
			prompt += "\n\n" + skill.GetSkillsPrompt(conf)
		}
	}

	return prompt
}

func GetSystemPrompt(conf *agentdomain.SessionConfig) string {
	return JoinPromptWithSkill(getSystemPromptFile(conf), conf)
}

func GetSubAgentPrompt(conf *agentdomain.SessionConfig) string {
	return JoinPromptWithSkill(SubAgentPrompt, conf)
}

// getToolConf is a helper to find a tool configuration by name
func getToolConf(conf *agentdomain.SessionConfig, name string) *agentdomain.ToolConfig {
	if conf == nil {
		return nil
	}
	for _, t := range conf.Tools {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

func GetToolPrompt(conf *agentdomain.SessionConfig) string {
	var sb strings.Builder
	sb.WriteString("**Available tools:**\n")
	sb.WriteString("You can use the following tools by a direct call:\n")
	for _, toolConf := range tools_conf.NativeTools {
		confTool := getToolConf(conf, toolConf.Name())
		if confTool == nil || !confTool.Enabled {
			continue
		}

		fmt.Fprintf(&sb, "- %s\n", toolConf.Name())
	}
	if subAgentConf := getToolConf(conf, "subagent"); subAgentConf != nil && subAgentConf.Enabled {
		fmt.Fprintf(&sb, "- subagent\n")
	}

	return sb.String()
}
