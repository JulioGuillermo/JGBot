package domain

type Status int

const (
	Normal = Status(iota)
	Writing
)
