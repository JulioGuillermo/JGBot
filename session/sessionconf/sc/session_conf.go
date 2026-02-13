package sc

import "JGBot/skill"

type SessionConf struct {
	Name   string
	ID     string
	Origin string
	Admin  string

	Allowed bool
	Respond Respond

	HistorySize int
	Provider    string

	SystemPromptFile string
	AgentMaxIters    int
	Tools            []Tool
	Skills           []Skill
}

func NewSessionConf(name, id, origin string) SessionConf {
	skConf := make([]Skill, 0)
	for _, skill := range skill.Skills {
		skConf = append(skConf, Skill{
			Name:        skill.Name,
			Enabled:     false,
			Description: skill.Description,
		})
	}

	return SessionConf{
		Name:   name,
		ID:     id,
		Origin: origin,
		Admin:  "",

		Allowed:       false,
		HistorySize:   50,
		AgentMaxIters: 3,
		Respond: Respond{
			Always: true,
		},
		SystemPromptFile: "",
		Tools: []Tool{
			{
				Name:    "message_reaction",
				Enabled: true,
			},
			{
				Name:    "javascript",
				Enabled: false,
			},
			{
				Name:    "skills",
				Enabled: false,
			},
			{
				Name:    "subagent",
				Enabled: false,
			},
			{
				Name:    "cron",
				Enabled: false,
			},
		},
		Skills: skConf,
	}
}

func (c *SessionConf) GetToolConf(name string) *Tool {
	for _, tool := range c.Tools {
		if tool.Name == name {
			return &tool
		}
	}
	return nil
}

func (c *SessionConf) GetSkillConf(name string) *Skill {
	for _, skill := range c.Skills {
		if skill.Name == name {
			return &skill
		}
	}
	return nil
}
