package store

import (
	"JGBot/conf"
	sessiondomain "JGBot/session/domain"
	"strings"
)

// toDomain converts sessionConfig to sessiondomain.SessionConfiguration
func (r *FileConfigRepository) toDomain(confJSON *sessionConfig) *sessiondomain.SessionConfiguration {
	if confJSON == nil {
		return nil
	}

	// Extract channel from ID (format "channel:chatID")
	channel := ""
	parts := strings.Split(confJSON.ID, ":")
	if len(parts) > 0 {
		channel = parts[0]
	}

	tools := make([]sessiondomain.Tool, 0, len(confJSON.Tools))
	for _, t := range confJSON.Tools {
		tools = append(tools, sessiondomain.Tool{
			Name:    t.Name,
			Enabled: t.Enabled,
		})
	}
	skills := make([]sessiondomain.Skill, 0, len(confJSON.Skills))
	for _, s := range confJSON.Skills {
		skills = append(skills, sessiondomain.Skill{
			Name:        s.Name,
			Enabled:     s.Enabled,
			Description: s.Description,
		})
	}

	c := &sessiondomain.SessionConfiguration{
		Origin:   confJSON.Origin,
		Channel:  channel,
		ChatID:   confJSON.ID,
		ChatName: confJSON.Name,
		Admin:    confJSON.Admin,

		Allowed: confJSON.Allowed,
		Respond: sessiondomain.Respond{
			Always: confJSON.Respond.Always,
			Match:  confJSON.Respond.Match,
		},

		HistorySize:      confJSON.HistorySize,
		Provider:         confJSON.Provider,
		SystemPromptFile: confJSON.SystemPromptFile,
		AgentMaxIters:    confJSON.AgentMaxIters,
		ShowThink:        confJSON.ShowThink,
		Tools:            tools,
		Skills:           skills,
	}

	return c
}

func applyDefConf(c *sessiondomain.SessionConfiguration, config *conf.Config) {
	if config == nil {
		return
	}
	if config.DefConf != nil {
		applyDefConfig(c, config.DefConf)
	}
	chanConf := config.GetChannelByName(c.Channel)
	if chanConf != nil && chanConf.DefConf != nil {
		applyDefConfig(c, chanConf.DefConf)
	}
}

func applyDefConfig(c *sessiondomain.SessionConfiguration, config *conf.DefConf) {
	if config == nil {
		return
	}
	if config.Allowed != nil {
		c.Allowed = *config.Allowed
	}
	if config.Respond != nil {
		if config.Respond.Always != nil {
			c.Respond.Always = *config.Respond.Always
		}
		if config.Respond.Match != nil {
			c.Respond.Match = *config.Respond.Match
		}
	}
	if config.HistorySize != nil {
		c.HistorySize = *config.HistorySize
	}
	if config.Provider != nil {
		c.Provider = *config.Provider
	}
	if config.SystemPromptFile != nil {
		c.SystemPromptFile = *config.SystemPromptFile
	}
	if config.AgentMaxIters != nil {
		c.AgentMaxIters = *config.AgentMaxIters
	}
	if config.ShowThink != nil {
		c.ShowThink = *config.ShowThink
	}
	applyDefTools(c, config)
	applyDefSkills(c, config)
}

func applyDefTools(c *sessiondomain.SessionConfiguration, config *conf.DefConf) {
	if config == nil || config.Tools == nil {
		return
	}
	for _, tool := range *config.Tools {
		if tool.Name == nil {
			continue
		}
		found := false
		for i := range c.Tools {
			if c.Tools[i].Name == *tool.Name {
				if tool.Enabled != nil {
					c.Tools[i].Enabled = *tool.Enabled
				}
				found = true
				break
			}
		}
		if !found {
			t := sessiondomain.Tool{Name: *tool.Name, Enabled: false}
			if tool.Enabled != nil {
				t.Enabled = *tool.Enabled
			}
			c.Tools = append(c.Tools, t)
		}
	}
}

func applyDefSkills(c *sessiondomain.SessionConfiguration, config *conf.DefConf) {
	if config == nil || config.Skills == nil {
		return
	}
	for _, s := range *config.Skills {
		if s.Name == nil {
			continue
		}
		found := false
		for i := range c.Skills {
			if c.Skills[i].Name == *s.Name {
				if s.Enabled != nil {
					c.Skills[i].Enabled = *s.Enabled
				}
				if s.Description != nil {
					c.Skills[i].Description = *s.Description
				}
				found = true
				break
			}
		}
		if !found {
			newSkill := sessiondomain.Skill{Name: *s.Name, Enabled: false}
			if s.Enabled != nil {
				newSkill.Enabled = *s.Enabled
			}
			if s.Description != nil {
				newSkill.Description = *s.Description
			}
			c.Skills = append(c.Skills, newSkill)
		}
	}
}
