package sessionconf

import (
	"JGBot/conf"
	"JGBot/session/sessionconf/sc"
)

type SessionCtl struct {
	Config   *SessionConfs
	Unconfig *SessionConfs
}

func NewSessionCtl() (*SessionCtl, error) {
	ctl := &SessionCtl{}
	var err error

	ctl.Config, err = NewSessionConfs(conf.SessionFile)
	if err != nil {
		return nil, err
	}

	ctl.Unconfig, err = NewSessionConfs(conf.UnconfigSessionFile)
	if err != nil {
		return nil, err
	}

	ctl.Config.OnChange = ctl.onConfigChange
	ctl.Config.Watch()

	return ctl, nil
}

func (s *SessionCtl) GetConfigOrigin(origin string) *sc.SessionConf {
	return s.Config.GetOrigin(origin)
}

func (s *SessionCtl) AddUnconfig(name, id, origin, channel string) *sc.SessionConf {
	session := s.getNewConfig(name, id, origin, channel)
	s.Unconfig.SetSession(session)
	return &session
}

func (s *SessionCtl) AddConfig(name, id, origin, channel string) *sc.SessionConf {
	session := s.getNewConfig(name, id, origin, channel)
	s.Config.SetSession(session)
	return &session
}

func (s *SessionCtl) getNewConfig(name, id, origin, channel string) sc.SessionConf {
	session := sc.NewSessionConf(name, id, origin)

	conf := getDefConfig(channel)
	if conf == nil {
		return session
	}

	session.Allowed = conf.Allowed
	session.Respond = sc.Respond(conf.Respond)
	session.HistorySize = conf.HistorySize
	session.Provider = conf.Provider
	session.SystemPromptFile = conf.SystemPromptFile
	session.AgentMaxIters = conf.AgentMaxIters
	if len(conf.Tools) > 0 {
		session.Tools = convertTools(conf.Tools)
	}
	if len(conf.Skills) > 0 {
		session.Skills = convertSkills(conf.Skills)
	}

	return session
}

func (s *SessionCtl) onConfigChange() {
	for _, conf := range s.Config.Sessions {
		s.Unconfig.RemoveOrigin(conf.Origin)
	}
}

func (s *SessionCtl) Close() {
	s.Config.Close()
	s.Unconfig.Close()
}
