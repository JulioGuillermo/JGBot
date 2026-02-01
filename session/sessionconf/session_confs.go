package sessionconf

import (
	"JGBot/session/sessionconf/sc"
	"JGBot/tools"

	"github.com/fsnotify/fsnotify"
)

type SessionConfs struct {
	filePath string
	watcher  *tools.FileWatcher
	Sessions []sc.SessionConf
	OnChange func()
}

func NewSessionConfs(filePath string) (*SessionConfs, error) {
	sessions := &SessionConfs{
		filePath: filePath,
		Sessions: []sc.SessionConf{},
	}
	err := sessions.Load()
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (s *SessionConfs) Load() error {
	err := tools.ReadJSONFile(s.filePath, &s.Sessions)
	if s.OnChange != nil {
		s.OnChange()
	}
	return err
}

func (s *SessionConfs) Save() error {
	return tools.WriteJSONFile(s.filePath, s.Sessions)
}

func (s *SessionConfs) Watch() {
	s.watcher, _ = tools.NewFileWatcher(s.filePath)
	s.watcher.OnChange = func(event fsnotify.Event) {
		s.Load()
	}
	s.watcher.OnError = func(err error) {
		s.Load()
	}
}

func (s *SessionConfs) Close() {
	if s.watcher != nil {
		s.watcher.Close()
	}
}

func (s *SessionConfs) Add(session sc.SessionConf) {
	s.Sessions = append(s.Sessions, session)
	s.Save()
}

func (s *SessionConfs) GetOriginIndex(origin string) int {
	for i, conf := range s.Sessions {
		if conf.Origin == origin {
			return i
		}
	}
	return -1
}

func (s *SessionConfs) GetOrigin(origin string) *sc.SessionConf {
	idx := s.GetOriginIndex(origin)
	if idx == -1 {
		return nil
	}
	return &s.Sessions[idx]
}

func (s *SessionConfs) HasOrigin(origin string) bool {
	return s.GetOriginIndex(origin) != -1
}

func (s *SessionConfs) RemoveOrigin(origin string) {
	idx := s.GetOriginIndex(origin)
	if idx == -1 {
		return
	}
	s.Sessions = append(s.Sessions[:idx], s.Sessions[idx+1:]...)
	s.Save()
}
