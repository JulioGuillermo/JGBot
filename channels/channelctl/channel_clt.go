package channelctl

import (
	"JGBot/channels"
	"JGBot/channels/telegramchannel"
	"JGBot/channels/whatsappchannel"
	"JGBot/conf"
	"fmt"
)

type ChannelCtl struct {
	channels map[string]channels.Channel
}

func InitChannelCtl() (*ChannelCtl, error) {
	ctl := &ChannelCtl{
		channels: make(map[string]channels.Channel),
	}

	if conf.Conf.Channels.Telegram.Enabled {
		channel, err := telegramchannel.NewTelegramChannel()
		if err != nil {
			return nil, err
		}
		ctl.channels[channel.GetName()] = channel
	}

	if conf.Conf.Channels.Whatsapp.Enabled {
		channel, err := whatsappchannel.NewWhatsAppChannel()
		if err != nil {
			return nil, err
		}
		ctl.channels[channel.GetName()] = channel
	}

	return ctl, nil
}

func (ctl *ChannelCtl) OnMessage(handler channels.OnMessageHandler) {
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

func (ctl *ChannelCtl) SendMessage(channel string, chatID uint, message string) error {
	ch, ok := ctl.channels[channel]
	if !ok {
		return fmt.Errorf("Not channel found with this name: %s", channel)
	}

	return ch.SendMessage(chatID, message)
}

func (ctl *ChannelCtl) Close() {
	for _, channel := range ctl.channels {
		channel.Close()
	}
}
