package channelsdomain

type Status int

const (
	Normal = Status(iota)
	Writing
)
