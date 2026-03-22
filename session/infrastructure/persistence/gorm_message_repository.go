package persistence

import (
	sessiondomain "JGBot/session/domain"
	"JGBot/session/infrastructure/persistence/db"
	"gorm.io/gorm"
)

// GormMessageRepository implements MessageRepository using GORM
type GormMessageRepository struct {
	db *gorm.DB
}

// NewGormMessageRepository creates a new GormMessageRepository
func NewGormMessageRepository(db *gorm.DB) *GormMessageRepository {
	return &GormMessageRepository{db: db}
}

// Save saves a message to the database
func (r *GormMessageRepository) Save(msg *sessiondomain.Message) error {
	_, err := db.SaveMessage(
		r.db,
		msg.Channel,
		msg.ChatID,
		msg.ChatName,
		msg.SenderID,
		msg.SenderName,
		msg.MessageID,
		msg.Message,
		msg.Role,
		msg.Extra,
	)
	return err
}

// GetHistory retrieves message history
func (r *GormMessageRepository) GetHistory(channel string, chatID uint, limit int) ([]*sessiondomain.Message, error) {
	dbMsgs, err := db.GetHistory(r.db, channel, chatID, limit)
	if err != nil {
		return nil, err
	}

	msgs := make([]*sessiondomain.Message, len(dbMsgs))
	for i, dbMsg := range dbMsgs {
		msgs[i] = dbMsg.ToDomain()
	}
	return msgs, nil
}

// ClearHistory deletes all messages for a channel and chat
func (r *GormMessageRepository) ClearHistory(channel string, chatID uint) error {
	return db.ClearHistory(r.db, channel, chatID)
}

// Ensure GormMessageRepository implements sessiondomain.MessageRepository
var _ sessiondomain.MessageRepository = (*GormMessageRepository)(nil)
