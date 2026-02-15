package channelctl

import (
	"JGBot/channels"
	"fmt"
)

func (ctl *ChannelCtl) SupportReaction(channel string) bool {
	ch, ok := ctl.channels[channel]
	if !ok {
		return false
	}
	_, ok = ch.(channels.ReactionChannel)
	return ok
}

func (ctl *ChannelCtl) ReactMessage(channel string, chatID uint, messageID uint, reaction string) error {
	ch, ok := ctl.channels[channel]
	if !ok {
		return fmt.Errorf("Not channel found with this name: %s", channel)
	}

	reactionChannel, ok := ch.(channels.ReactionChannel)
	if !ok {
		return channels.ErrNotSupported("reaction")
	}

	return reactionChannel.ReactMessage(chatID, messageID, reaction)
}
