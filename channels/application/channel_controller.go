package channelsapplication

import channelsdomain "JGBot/channels/domain"

type ChannelController struct {
	channels map[string]channelsdomain.Channel
}

func NewChannelController() channelsdomain.ChannelController {
	return &ChannelController{
		channels: make(map[string]channelsdomain.Channel),
	}
}

func (ctl *ChannelController) SetChannels(channels []channelsdomain.Channel) {
	for _, channel := range channels {
		ctl.channels[channel.GetName()] = channel
	}
}

func (ctl *ChannelController) GetChannel(name string) (channelsdomain.Channel, error) {
	channel, ok := ctl.channels[name]
	if !ok {
		return nil, channelsdomain.ErrChannelNotFound
	}
	return channel, nil
}

func (ctl *ChannelController) SetChannel(channel channelsdomain.Channel) {
	ctl.channels[channel.GetName()] = channel
}

func (ctl *ChannelController) DelChannel(name string) {
	delete(ctl.channels, name)
}

func (ctl *ChannelController) OnMessage(handler channelsdomain.MessageHandler) {
	for _, channel := range ctl.channels {
		channel.OnMessage(handler)
	}
}

func (ctl *ChannelController) Close() {
	for _, channel := range ctl.channels {
		channel.Close()
	}
}
