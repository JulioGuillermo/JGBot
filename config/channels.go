package config

type Channel struct {
	Enabled  bool
	LogLevel string
	Config   map[string]string
}

type Channels struct {
	Telegram Channel
	Whatsapp Channel
	// Others map[string]*Channel
}

func newChannels() *Channels {
	return &Channels{
		Telegram: Channel{
			Enabled:  false,
			LogLevel: "info",
			Config:   map[string]string{},
		},
		Whatsapp: Channel{
			Enabled:  false,
			LogLevel: "info",
			Config:   map[string]string{},
		},
	}
}
