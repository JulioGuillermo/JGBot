package sessionconf

type SessionConf struct {
	Name         string
	ID           string
	Origin       string
	OriginRegExp bool

	Allowed bool
}

type SessionConfs struct {
	filePath  string
	hotReload bool
	Sessions  []SessionConf
}

func NewSessionConfs(filePath string, hotReload bool) (*SessionConfs, error) {
	sessions := &SessionConfs{
		filePath:  filePath,
		hotReload: hotReload,
	}

	return sessions, nil
}

func (s *SessionConfs) Load() error {
}
