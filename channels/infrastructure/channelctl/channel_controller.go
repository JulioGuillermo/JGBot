package channelctl

import (
	channelsapplication "JGBot/channels/application"
	channelsdomain "JGBot/channels/domain"
	"JGBot/channels/infrastructure/telegram"
	"JGBot/channels/infrastructure/whatsapp"
	"JGBot/conf"
)

func NewChannelCtl() (channelsdomain.ChannelController, error) {
	controller := channelsapplication.NewChannelController()

	if conf.Conf.Channels.Telegram.Enabled {
		channel, err := telegram.NewTelegramChannel()
		if err != nil {
			return nil, err
		}
		controller.SetChannel(channel)
	}

	if conf.Conf.Channels.Whatsapp.Enabled {
		channel, err := whatsapp.NewWhatsAppChannel()
		if err != nil {
			return nil, err
		}
		controller.SetChannel(channel)
	}

	return controller, nil
}
