package store

import (
	"JGBot/conf"
	sessiondomain "JGBot/session/domain"
)

// sessionConfig represents the session configuration stored in JSON
type sessionConfig struct {
	Name   string `json:"Name"`
	ID     string `json:"ID"`
	Origin string `json:"Origin"`
	Admin  string `json:"Admin"`

	Allowed bool `json:"Allowed"`
	Respond struct {
		Always bool   `json:"Always"`
		Match  string `json:"Match"`
	} `json:"Respond"`

	HistorySize      int     `json:"HistorySize"`
	Provider         string  `json:"Provider"`
	SystemPromptFile string  `json:"SystemPromptFile"`
	AgentMaxIters    int     `json:"AgentMaxIters"`
	ShowThink        bool    `json:"ShowThink"`
	Tools            []Tool  `json:"Tools"`
	Skills           []Skill `json:"Skills"`
}

type Tool struct {
	Name    string `json:"Name"`
	Enabled bool   `json:"Enabled"`
}

type Skill struct {
	Name        string `json:"Name"`
	Enabled     bool   `json:"Enabled"`
	Description string `json:"Description"`
}

// newSessionConfig creates a new sessionConfig with default values
func newSessionConfig(name, id, origin, channel string) sessionConfig {
	// 1. Create a domain config with basic defaults
	domainConf := sessiondomain.NewSessionConfiguration(name, id, origin, channel)

	// 3. Apply system-wide and channel-specific defaults from global config
	applyDefConf(domainConf, conf.Conf)

	tools := make([]Tool, 0, len(domainConf.Tools))
	for _, t := range domainConf.Tools {
		tools = append(tools, Tool{
			Name:    t.Name,
			Enabled: t.Enabled,
		})
	}
	skills := make([]Skill, 0, len(domainConf.Skills))
	for _, s := range domainConf.Skills {
		skills = append(skills, Skill{
			Name:        s.Name,
			Enabled:     s.Enabled,
			Description: s.Description,
		})
	}

	return sessionConfig{
		Name:   name,
		ID:     id,
		Origin: origin,
		Admin:  domainConf.Admin,

		Allowed: domainConf.Allowed,
		Respond: struct {
			Always bool   `json:"Always"`
			Match  string `json:"Match"`
		}{
			Always: domainConf.Respond.Always,
			Match:  domainConf.Respond.Match,
		},

		HistorySize:      domainConf.HistorySize,
		Provider:         domainConf.Provider,
		SystemPromptFile: domainConf.SystemPromptFile,
		AgentMaxIters:    domainConf.AgentMaxIters,
		ShowThink:        domainConf.ShowThink,
		Tools:            tools,
		Skills:           skills,
	}
}
