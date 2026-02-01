package sc

type SessionConf struct {
	Name   string
	ID     string
	Origin string

	Allowed bool

	Respond Respond
}

func NewSessionConf(name, id, origin string) SessionConf {
	return SessionConf{
		Name:   name,
		ID:     id,
		Origin: origin,

		Allowed: false,
		Respond: Respond{
			Always: true,
		},
	}
}
