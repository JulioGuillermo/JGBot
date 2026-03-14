package channelctl

import (
	"JGBot/channels/application"
	"JGBot/channels/domain"
	"JGBot/channels/infrastructure/telegram"
	"JGBot/channels/infrastructure/whatsapp"
	"JGBot/conf"
)

func NewChannelCtl() (domain.ChannelController, error) {
	controller := application.NewChannelController()

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
