package channelsdomain

type ChannelController interface {
	SetChannels(channels []Channel)
	GetChannel(name string) (Channel, error)
	SetChannel(channel Channel)
	DelChannel(name string)
	OnMessage(handler MessageHandler)
	Close()
}
