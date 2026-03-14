package sc

import (
	"JGBot/conf"
)

func (c *SessionConf) ApplyDefConf(channel string) {
	if conf.Conf.DefConf != nil {
		c.applyDefConfig(conf.Conf.DefConf)
	}
	c.applyChannelDefConf(channel)
}

func (c *SessionConf) applyChannelDefConf(channel string) {
	chanConf := conf.Conf.GetChannelByName(channel)
	if chanConf == nil || chanConf.DefConf == nil {
		return
	}
	c.applyDefConfig(chanConf.DefConf)
}

func (c *SessionConf) applyDefConfig(config *conf.DefConf) {
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

	c.applyDefTools(config)
	c.applyDefSkills(config)
}

func (c *SessionConf) applyDefTools(config *conf.DefConf) {
	if config == nil || config.Tools == nil {
		return
	}

	for _, tool := range *config.Tools {
		c.applyDefTool(tool)
	}
}

func (c *SessionConf) applyDefTool(tool conf.Tool) {
	if tool.Name == nil {
		return
	}

	for _, t := range c.Tools {
		if t.Name != *tool.Name {
			continue
		}

		if tool.Enabled != nil {
			t.Enabled = *tool.Enabled
		}
		return
	}

	t := Tool{
		Name:    *tool.Name,
		Enabled: false,
	}

	if tool.Enabled != nil {
		t.Enabled = *tool.Enabled
	}

	c.Tools = append(c.Tools, t)
}

func (c *SessionConf) applyDefSkills(config *conf.DefConf) {
	if config == nil || config.Skills == nil {
		return
	}

	for _, skill := range *config.Skills {
		c.applyDefSkill(skill)
	}
}

func (c *SessionConf) applyDefSkill(skill conf.Skill) {
	if skill.Name == nil {
		return
	}

	for _, s := range c.Skills {
		if s.Name != *skill.Name {
			continue
		}

		if skill.Enabled != nil {
			s.Enabled = *skill.Enabled
		}

		if skill.Description != nil {
			s.Description = *skill.Description
		}

		return
	}

	s := Skill{
		Name:        *skill.Name,
		Enabled:     false,
		Description: "",
	}

	if skill.Enabled != nil {
		s.Enabled = *skill.Enabled
	}

	if skill.Description != nil {
		s.Description = *skill.Description
	}

	c.Skills = append(c.Skills, s)
}
