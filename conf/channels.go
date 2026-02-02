package conf

type Channel struct {
	Enabled           bool
	AutoEnableSession bool
	Config            map[string]string
}

type Channels struct {
	Telegram Channel
	Whatsapp Channel
	// Others map[string]*Channel
}

func newChannels() *Channels {
	return &Channels{
		Telegram: Channel{
			Enabled:           false,
			AutoEnableSession: false,
			Config:            map[string]string{},
		},
		Whatsapp: Channel{
			Enabled:           false,
			AutoEnableSession: false,
			Config:            map[string]string{},
		},
	}
}
