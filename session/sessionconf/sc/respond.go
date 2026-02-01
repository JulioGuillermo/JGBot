package sc

type Respond struct {
	Always bool
	Match  string
}

func (r Respond) Respond(text string) bool {
	if r.Always {
		return true
	}
	if r.Match != "" {
		return text == r.Match
	}
	return false
}
