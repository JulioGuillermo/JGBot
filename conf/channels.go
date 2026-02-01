package conf

type Channel struct {
	Enabled bool
	Config  map[string]string
}

type Channels struct {
	Telegram Channel
	Whatsapp Channel
	// Others map[string]*Channel
}

func newChannels() *Channels {
	return &Channels{
		Telegram: Channel{
			Enabled: false,
			Config:  map[string]string{},
		},
		Whatsapp: Channel{
			Enabled: false,
			Config:  map[string]string{},
		},
	}
}
