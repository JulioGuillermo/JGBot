package store

import (
	"JGBot/conf"
	sessiondomain "JGBot/session/domain"
	"JGBot/tools"
	"regexp"

	"github.com/fsnotify/fsnotify"
)

// FileConfigRepository implements sessiondomain.ConfigurationRepository using JSON files.
// This is a pure domain implementation that replaces the legacy sessionconf package.
type FileConfigRepository struct {
	configFile    string
	unconfigFile  string
	configs       []sessionConfig
	unconfigs     []sessionConfig
	configWatcher *tools.FileWatcher
}

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

	HistorySize   int `json:"HistorySize"`
	AgentMaxIters int `json:"AgentMaxIters"`
}

// Ensure FileConfigRepository implements sessiondomain.ConfigurationRepository
var _ sessiondomain.ConfigurationRepository = (*FileConfigRepository)(nil)

// NewFileConfigRepository creates a new FileConfigRepository
func NewFileConfigRepository() (*FileConfigRepository, error) {
	repo := &FileConfigRepository{
		configFile:   conf.SessionFile,
		unconfigFile: conf.UnconfigSessionFile,
		configs:      []sessionConfig{},
		unconfigs:    []sessionConfig{},
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	repo.watch()

	return repo, nil
}

func (r *FileConfigRepository) load() error {
	if err := tools.ReadJSONFile(r.configFile, &r.configs); err != nil {
		r.configs = []sessionConfig{}
	}
	if err := tools.ReadJSONFile(r.unconfigFile, &r.unconfigs); err != nil {
		r.unconfigs = []sessionConfig{}
	}
	return nil
}

func (r *FileConfigRepository) save() error {
	if err := tools.WriteJSONFile(r.configFile, r.configs); err != nil {
		return err
	}
	return tools.WriteJSONFile(r.unconfigFile, r.unconfigs)
}

func (r *FileConfigRepository) watch() {
	r.configWatcher, _ = tools.NewFileWatcher(r.configFile)
	r.configWatcher.OnChange = func(event fsnotify.Event) {
		r.load()
	}
	r.configWatcher.OnError = func(err error) {
		r.load()
	}
}

// Close closes the file watcher
func (r *FileConfigRepository) Close() {
	if r.configWatcher != nil {
		r.configWatcher.Close()
	}
}

func (r *FileConfigRepository) GetConfig(origin string) *sessiondomain.SessionConfiguration {
	for i := range r.configs {
		if r.configs[i].Origin == origin {
			return r.toDomain(&r.configs[i])
		}
	}
	return nil
}

func (r *FileConfigRepository) GetConfigByChannel(channel string, chatID uint) *sessiondomain.SessionConfiguration {
	// Not implemented in legacy code either
	return nil
}

func (r *FileConfigRepository) CreateConfig(chatName, sessionID, origin, channel string) *sessiondomain.SessionConfiguration {
	// Remove from unconfigs if present
	r.removeUnconfig(origin)

	conf := newSessionConfig(chatName, sessionID, origin)
	r.configs = append(r.configs, conf)
	r.save()
	return r.toDomain(&r.configs[len(r.configs)-1])
}

func (r *FileConfigRepository) CreateUnconfig(chatName, sessionID, origin, channel string) *sessiondomain.SessionConfiguration {
	// Remove from configs if present
	r.removeConfig(origin)

	conf := newSessionConfig(chatName, sessionID, origin)
	r.unconfigs = append(r.unconfigs, conf)
	r.save()
	return r.toDomain(&r.unconfigs[len(r.unconfigs)-1])
}

func (r *FileConfigRepository) removeConfig(origin string) {
	for i := range r.configs {
		if r.configs[i].Origin == origin {
			r.configs = append(r.configs[:i], r.configs[i+1:]...)
			return
		}
	}
}

func (r *FileConfigRepository) removeUnconfig(origin string) {
	for i := range r.unconfigs {
		if r.unconfigs[i].Origin == origin {
			r.unconfigs = append(r.unconfigs[:i], r.unconfigs[i+1:]...)
			return
		}
	}
}

// newSessionConfig creates a new sessionConfig with default values
func newSessionConfig(name, id, origin string) sessionConfig {
	return sessionConfig{
		Name:   name,
		ID:     id,
		Origin: origin,
		Admin:  "",

		Allowed:       false,
		HistorySize:   50,
		AgentMaxIters: 3,
	}
}

// toDomain converts sessionConfig to sessiondomain.SessionConfiguration
func (r *FileConfigRepository) toDomain(conf *sessionConfig) *sessiondomain.SessionConfiguration {
	if conf == nil {
		return nil
	}
	shouldRespond := func(msg string) bool {
		if conf.Respond.Always {
			return true
		}
		re, err := regexp.Compile(conf.Respond.Match)
		if err != nil {
			return false
		}
		return re.MatchString(msg)
	}
	return &sessiondomain.SessionConfiguration{
		Origin:        conf.Origin,
		Channel:       "",
		ChatID:        conf.ID,
		ChatName:      conf.Name,
		HistorySize:   conf.HistorySize,
		Admin:         conf.Admin,
		Allowed:       conf.Allowed,
		ShouldRespond: shouldRespond,
	}
}
