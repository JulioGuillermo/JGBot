package sc

import (
	"JGBot/log"
	"regexp"
)

type Respond struct {
	Always bool
	Match  string
}

func (r Respond) Respond(text string) bool {
	if r.Always {
		return true
	}

	re, err := regexp.Compile(r.Match)
	if err != nil {
		log.Error("Fail to compile match", "error", err)
		return false
	}

	return re.MatchString(text)
}
