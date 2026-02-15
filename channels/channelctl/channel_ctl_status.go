package channelctl

import (
	"JGBot/channels"
	"fmt"
)

func (ctl *ChannelCtl) SupportStatus(channel string) bool {
	ch, ok := ctl.channels[channel]
	if !ok {
		return false
	}
	_, ok = ch.(channels.StatusChannel)
	return ok
}

func (ctl *ChannelCtl) Status(channel string, chatID uint, status channels.Status) error {
	ch, ok := ctl.channels[channel]
	if !ok {
		return fmt.Errorf("Not channel found with this name: %s", channel)
	}

	statusChannel, ok := ch.(channels.StatusChannel)
	if !ok {
		return channels.ErrNotSupported("status")
	}

	return statusChannel.Status(chatID, status)
}
