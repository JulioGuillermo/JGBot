package sessiondomain

import (
	"JGBot/log"
	"JGBot/skill"
	"regexp"
)

type Respond struct {
	Always bool
	Match  string
}

func (r Respond) ShouldRespond(text string) bool {
	if r.Always {
		return true
	}
	if r.Match == "" {
		return false
	}
	re, err := regexp.Compile(r.Match)
	if err != nil {
		log.Error("Failed to compile respond match regex", "error", err, "match", r.Match)
		return false
	}
	return re.MatchString(text)
}

type Tool struct {
	Name    string
	Enabled bool
}

type Skill struct {
	Name        string
	Enabled     bool
	Description string
}

type SessionConfiguration struct {
	Origin   string
	Channel  string
	ChatID   string
	ChatName string
	Admin    string

	Allowed bool
	Respond Respond

	HistorySize int
	Provider    string

	SystemPromptFile string
	AgentMaxIters    int
	ShowThink        bool

	Tools  []Tool
	Skills []Skill
}

func NewSessionConfiguration(name, id, origin, channel string) *SessionConfiguration {
	skConf := make([]Skill, 0)
	for _, s := range skill.Skills {
		skConf = append(skConf, Skill{
			Name:        s.Name,
			Enabled:     false,
			Description: s.Description,
		})
	}

	return &SessionConfiguration{
		Origin:   origin,
		Channel:  channel,
		ChatID:   id,
		ChatName: name,
		Admin:    "",

		Allowed:       false,
		HistorySize:   50,
		AgentMaxIters: 3,
		ShowThink:     false,
		Respond: Respond{
			Always: true,
		},
		SystemPromptFile: "",
		Tools: []Tool{
			{Name: "message_reaction", Enabled: true},
			{Name: "javascript", Enabled: false},
			{Name: "skills", Enabled: false},
			{Name: "subagent", Enabled: false},
			{Name: "cron", Enabled: false},
		},
		Skills: skConf,
	}
}

func (c *SessionConfiguration) GetToolConf(name string) *Tool {
	for i := range c.Tools {
		if c.Tools[i].Name == name {
			return &c.Tools[i]
		}
	}
	return nil
}

func (c *SessionConfiguration) GetSkillConf(name string) *Skill {
	for i := range c.Skills {
		if c.Skills[i].Name == name {
			return &c.Skills[i]
		}
	}
	return nil
}

func (c *SessionConfiguration) ShouldRespond(text string) bool {
	return c.Respond.ShouldRespond(text)
}
