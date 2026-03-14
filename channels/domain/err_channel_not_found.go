package domain

var ErrChannelNotFound = &ChannelNotFoundError{}

type ChannelNotFoundError struct{}

func (e *ChannelNotFoundError) Error() string {
	return "channel not found"
}
