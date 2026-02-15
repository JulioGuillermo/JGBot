package channels

type SimpleMenu struct {
	Title   string
	Message string
	Buttons []SimpleMenuButton
}

type SimpleMenuButton struct {
	ID    string
	Label string
}

type MenuCapableChannel interface {
	SendSimpleMenu(chatID uint, menu SimpleMenu) error
}
