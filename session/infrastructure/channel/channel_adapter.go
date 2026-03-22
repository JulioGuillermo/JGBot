package sessioninfrastructure

import (
	channelsdomain "JGBot/channels/domain"
	sessiondomain "JGBot/session/domain"
)

type ChannelAdapter struct {
	channelCtl channelsdomain.ChannelController
}

func NewChannelAdapter(channelCtl channelsdomain.ChannelController) *ChannelAdapter {
	return &ChannelAdapter{
		channelCtl: channelCtl,
	}
}

func (a *ChannelAdapter) SendMessage(channelID string, chatID uint, message string) error {
	channel, err := a.channelCtl.GetChannel(channelID)
	if err != nil {
		return err
	}
	return channel.SendMessage(chatID, message)
}

func (a *ChannelAdapter) SendReaction(channelID string, chatID uint, messageID uint, reaction string) error {
	channel, err := a.channelCtl.GetChannel(channelID)
	if err != nil {
		return err
	}
	return channel.SendMessageReaction(chatID, messageID, reaction)
}

func (a *ChannelAdapter) SendStatus(channelID string, chatID uint, status string) error {
	channel, err := a.channelCtl.GetChannel(channelID)
	if err != nil {
		return err
	}
	var s channelsdomain.Status
	switch status {
	case "writing":
		s = channelsdomain.Writing
	default:
		s = channelsdomain.Normal
	}
	return channel.SendStatus(chatID, s)
}

func (a *ChannelAdapter) GetChannelAutoEnable(channelID string) (bool, error) {
	channel, err := a.channelCtl.GetChannel(channelID)
	if err != nil {
		return false, err
	}
	return channel.AutoEnableSession(), nil
}

// Ensure ChannelAdapter implements sessiondomain.ChannelService
var _ sessiondomain.ChannelService = (*ChannelAdapter)(nil)
