package persistence

import (
	"testing"

	sessiondomain "JGBot/session/domain"
)

func TestGormMessageRepository_ImplementsInterface(t *testing.T) {
	t.Parallel()

	// Compile-time check that GormMessageRepository implements domain.MessageRepository
	var _ sessiondomain.MessageRepository = (*GormMessageRepository)(nil)
}
