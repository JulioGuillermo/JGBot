package application

import "JGBot/channels/domain"

type ChannelController struct {
	channels map[string]domain.Channel
}

func NewChannelController() domain.ChannelController {
	return &ChannelController{
		channels: make(map[string]domain.Channel),
	}
}

func (ctl *ChannelController) SetChannels(channels []domain.Channel) {
	for _, channel := range channels {
		ctl.channels[channel.GetName()] = channel
	}
}

func (ctl *ChannelController) GetChannel(name string) (domain.Channel, error) {
	channel, ok := ctl.channels[name]
	if !ok {
		return nil, domain.ErrChannelNotFound
	}
	return channel, nil
}

func (ctl *ChannelController) SetChannel(channel domain.Channel) {
	ctl.channels[channel.GetName()] = channel
}

func (ctl *ChannelController) DelChannel(name string) {
	delete(ctl.channels, name)
}

func (ctl *ChannelController) OnMessage(handler domain.MessageHandler) {
	for _, channel := range ctl.channels {
		channel.OnMessage(handler)
	}
}

func (ctl *ChannelController) Close() {
	for _, channel := range ctl.channels {
		channel.Close()
	}
}
