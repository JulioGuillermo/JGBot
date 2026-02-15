package channels

type StatusChannel interface {
	Status(chatID uint, status Status) error
}
