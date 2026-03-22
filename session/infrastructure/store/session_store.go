package store

import (
	agentdomain "JGBot/agent/domain"
	"JGBot/conf"
	sessiondomain "JGBot/session/domain"
	"JGBot/skill"
	"JGBot/tools"

	"github.com/fsnotify/fsnotify"
)

// FileSessionStore implements agentdomain.SessionStore using JSON files and GORM persistence.
// This is a pure domain implementation that replaces the legacy sessionconf and sessiondb packages.
type FileSessionStore struct {
	configFile    string
	unconfigFile  string
	configs       []agentdomain.SessionConfig
	unconfigs     []agentdomain.SessionConfig
	configWatcher *tools.FileWatcher
	msgRepo       sessiondomain.MessageRepository
}

// Ensure FileSessionStore implements agentdomain.SessionStore
var _ agentdomain.SessionStore = (*FileSessionStore)(nil)

// NewFileSessionStore creates a new FileSessionStore with the given message repository.
func NewFileSessionStore(msgRepo sessiondomain.MessageRepository) (*FileSessionStore, error) {
	// Type assert to the concrete type if it's not nil
	var repo sessiondomain.MessageRepository
	if msgRepo != nil {
		repo = msgRepo
	}

	store := &FileSessionStore{
		configFile:   conf.SessionFile,
		unconfigFile: conf.UnconfigSessionFile,
		configs:      []agentdomain.SessionConfig{},
		unconfigs:    []agentdomain.SessionConfig{},
		msgRepo:      repo,
	}

	if err := store.loadConfigs(); err != nil {
		return nil, err
	}

	store.watch()

	return store, nil
}

func (s *FileSessionStore) loadConfigs() error {
	if err := tools.ReadJSONFile(s.configFile, &s.configs); err != nil {
		s.configs = []agentdomain.SessionConfig{}
	}
	if err := tools.ReadJSONFile(s.unconfigFile, &s.unconfigs); err != nil {
		s.unconfigs = []agentdomain.SessionConfig{}
	}
	return nil
}

func (s *FileSessionStore) saveConfigs() error {
	if err := tools.WriteJSONFile(s.configFile, s.configs); err != nil {
		return err
	}
	return tools.WriteJSONFile(s.unconfigFile, s.unconfigs)
}

func (s *FileSessionStore) watch() {
	s.configWatcher, _ = tools.NewFileWatcher(s.configFile)
	s.configWatcher.OnChange = func(event fsnotify.Event) {
		s.loadConfigs()
	}
	s.configWatcher.OnError = func(err error) {
		s.loadConfigs()
	}
}

// Close closes the file watcher
func (s *FileSessionStore) Close() {
	if s.configWatcher != nil {
		s.configWatcher.Close()
	}
}

func (s *FileSessionStore) GetConfig(origin string) *agentdomain.SessionConfig {
	for i := range s.configs {
		if s.configs[i].Origin == origin {
			return &s.configs[i]
		}
	}
	return nil
}

func (s *FileSessionStore) GetConfigs() []*agentdomain.SessionConfig {
	result := make([]*agentdomain.SessionConfig, len(s.configs))
	for i := range s.configs {
		result[i] = &s.configs[i]
	}
	return result
}

func (s *FileSessionStore) CreateConfig(chatName, sessionID, origin, channel string) *agentdomain.SessionConfig {
	// Remove from unconfigs if present
	s.removeUnconfig(origin)

	conf := newAgentSessionConfig(chatName, sessionID, origin)
	s.configs = append(s.configs, conf)
	s.saveConfigs()
	return &s.configs[len(s.configs)-1]
}

func (s *FileSessionStore) CreateUnconfig(chatName, sessionID, origin, channel string) *agentdomain.SessionConfig {
	// Remove from configs if present
	s.removeConfig(origin)

	conf := newAgentSessionConfig(chatName, sessionID, origin)
	s.unconfigs = append(s.unconfigs, conf)
	s.saveConfigs()
	return &s.unconfigs[len(s.unconfigs)-1]
}

func (s *FileSessionStore) removeConfig(origin string) {
	for i := range s.configs {
		if s.configs[i].Origin == origin {
			s.configs = append(s.configs[:i], s.configs[i+1:]...)
			return
		}
	}
}

func (s *FileSessionStore) removeUnconfig(origin string) {
	for i := range s.unconfigs {
		if s.unconfigs[i].Origin == origin {
			s.unconfigs = append(s.unconfigs[:i], s.unconfigs[i+1:]...)
			return
		}
	}
}

func (s *FileSessionStore) GetHistory(channel string, chatID uint, limit int) ([]*agentdomain.SessionMessage, error) {
	msgs, err := s.msgRepo.GetHistory(channel, chatID, limit)
	if err != nil {
		return nil, err
	}
	result := make([]*agentdomain.SessionMessage, len(msgs))
	for i, msg := range msgs {
		result[i] = toAgentMessage(msg)
	}
	return result, nil
}

func (s *FileSessionStore) SaveMessage(channel string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string, role string, extra string) (*agentdomain.SessionMessage, error) {
	msg := &sessiondomain.Message{
		Channel:    channel,
		ChatID:     chatID,
		ChatName:   chatName,
		SenderID:   senderID,
		SenderName: senderName,
		MessageID:  messageID,
		Message:    message,
		Role:       role,
		Extra:      extra,
	}
	if err := s.msgRepo.Save(msg); err != nil {
		return nil, err
	}
	return toAgentMessage(msg), nil
}

func (s *FileSessionStore) ClearHistory(channel string, chatID uint) error {
	return s.msgRepo.ClearHistory(channel, chatID)
}

// newAgentSessionConfig creates a new SessionConfig with default values
func newAgentSessionConfig(name, id, origin string) agentdomain.SessionConfig {
	conf := agentdomain.SessionConfig{
		Name:   name,
		ID:     id,
		Origin: origin,
		Admin:  "",

		Allowed:       false,
		HistorySize:   50,
		AgentMaxIters: 3,
	}
	conf.Respond.Always = true
	conf.Respond.Match = ""

	// Default tools
	conf.Tools = []agentdomain.ToolConfig{
		{Name: "message_reaction", Enabled: true},
		{Name: "javascript", Enabled: false},
		{Name: "skills", Enabled: false},
		{Name: "subagent", Enabled: false},
		{Name: "cron", Enabled: false},
	}

	// Load skills from skill package
	conf.Skills = []agentdomain.SkillConfig{}
	for _, sk := range skill.Skills {
		conf.Skills = append(conf.Skills, agentdomain.SkillConfig{
			Name:        sk.Name,
			Enabled:     false,
			Description: sk.Description,
		})
	}

	return conf
}

// toAgentMessage converts a sessiondomain.Message to agentdomain.SessionMessage
func toAgentMessage(msg *sessiondomain.Message) *agentdomain.SessionMessage {
	if msg == nil {
		return nil
	}
	return &agentdomain.SessionMessage{
		ID:         msg.ID,
		CreatedAt:  msg.CreatedAt,
		UpdatedAt:  msg.UpdatedAt,
		Channel:    msg.Channel,
		ChatID:     msg.ChatID,
		ChatName:   msg.ChatName,
		SenderID:   msg.SenderID,
		SenderName: msg.SenderName,
		MessageID:  msg.MessageID,
		Message:    msg.Message,
		Role:       msg.Role,
		Extra:      msg.Extra,
	}
}
