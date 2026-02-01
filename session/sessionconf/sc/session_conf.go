package sc

type SessionConf struct {
	Name   string
	ID     string
	Origin string

	Allowed bool
	Respond Respond

	HistorySize int
	Provider    string
}

func NewSessionConf(name, id, origin string) SessionConf {
	return SessionConf{
		Name:   name,
		ID:     id,
		Origin: origin,

		Allowed:     false,
		HistorySize: 50,
		Respond: Respond{
			Always: true,
		},
	}
}
