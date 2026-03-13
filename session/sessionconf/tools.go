package sessionconf

import (
	"JGBot/conf"
	"JGBot/session/sessionconf/sc"
)

func getDefConfig(channel string) *conf.DefConf {
	// Check channel-specific default config
	if conf.Conf.Channels.Telegram.DefConf != nil && channel == "telegram" {
		return conf.Conf.Channels.Telegram.DefConf
	}
	if conf.Conf.Channels.Whatsapp.DefConf != nil && channel == "whatsapp" {
		return conf.Conf.Channels.Whatsapp.DefConf
	}

	// Fall back to global default config
	return conf.Conf.DefConf
}

func convertTools(tools []conf.Tool) []sc.Tool {
	result := make([]sc.Tool, len(tools))
	for i, t := range tools {
		result[i] = sc.Tool(t)
	}
	return result
}

func convertSkills(skills []conf.Skill) []sc.Skill {
	result := make([]sc.Skill, len(skills))
	for i, s := range skills {
		result[i] = sc.Skill(s)
	}
	return result
}
