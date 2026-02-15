package channelctl

import (
	"JGBot/channels"
	"fmt"
)

func (ctl *ChannelCtl) SupportMenus(channel string) bool {
	ch, ok := ctl.channels[channel]
	if !ok {
		return false
	}
	_, ok = ch.(channels.MenuCapableChannel)
	return ok
}

func (ctl *ChannelCtl) SendSimpleMenu(channel string, chatID uint, menu channels.SimpleMenu) error {
	ch, ok := ctl.channels[channel]
	if !ok {
		return fmt.Errorf("Not channel found with this name: %s", channel)
	}

	menuChannel, ok := ch.(channels.MenuCapableChannel)
	if !ok {
		return channels.ErrNotSupported("menu")
	}

	return menuChannel.SendSimpleMenu(chatID, menu)
}
