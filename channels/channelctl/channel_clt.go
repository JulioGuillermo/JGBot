package channelctl

import (
	"JGBot/channels/domain"
	"JGBot/channels/infrastructure/telegram"
	"JGBot/channels/infrastructure/whatsapp"
	"JGBot/conf"
	"fmt"
)

type ChannelCtl struct {
	channels map[string]domain.Channel
}

func InitChannelCtl() (*ChannelCtl, error) {
	ctl := &ChannelCtl{
		channels: make(map[string]domain.Channel),
	}

	if conf.Conf.Channels.Telegram.Enabled {
		channel, err := telegram.NewTelegramChannel()
		if err != nil {
			return nil, err
		}
		ctl.channels[channel.GetName()] = channel
	}

	if conf.Conf.Channels.Whatsapp.Enabled {
		channel, err := whatsapp.NewWhatsAppChannel()
		if err != nil {
			return nil, err
		}
		ctl.channels[channel.GetName()] = channel
	}

	return ctl, nil
}

func (ctl *ChannelCtl) OnMessage(handler domain.MessageHandler) {
	for _, channel := range ctl.channels {
		channel.OnMessage(handler)
	}
}

func (ctl *ChannelCtl) AutoEnableSession(channel string) bool {
	ch, ok := ctl.channels[channel]
	if !ok {
		return false
	}
	return ch.AutoEnableSession()
}

func (ctl *ChannelCtl) Status(channel string, chatID uint, status domain.Status) error {
	if ctl == nil {
		return nil
	}

	ch, ok := ctl.channels[channel]
	if !ok {
		return fmt.Errorf("Not channel found with this name: %s", channel)
	}

	return ch.SendStatus(chatID, status)
}

func (ctl *ChannelCtl) SendMessage(channel string, chatID uint, message string) error {
	ch, ok := ctl.channels[channel]
	if !ok {
		return fmt.Errorf("Not channel found with this name: %s", channel)
	}

	return ch.SendMessage(chatID, message)
}

func (ctl *ChannelCtl) ReactMessage(channel string, chatID uint, messageID uint, reaction string) error {
	ch, ok := ctl.channels[channel]
	if !ok {
		return fmt.Errorf("Not channel found with this name: %s", channel)
	}

	return ch.SendMessageReaction(chatID, messageID, reaction)
}

func (ctl *ChannelCtl) Close() {
	for _, channel := range ctl.channels {
		channel.Close()
	}
}
