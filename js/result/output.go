package result

import (
	"fmt"
	"strings"
)

type Output struct {
	Logs   string
	Result string
}

func (o *Output) String() string {
	logs := strings.TrimSpace(o.Logs)
	results := strings.TrimSpace(o.Result)
	if logs != "" && results != "" {
		return fmt.Sprintf("LOGS:\n%s\n\nRESULT:\n%s", o.Logs, o.Result)
	}
	if results != "" {
		return fmt.Sprintf("RESULT:\n%s", o.Result)
	}
	if logs != "" {
		return fmt.Sprintf("LOGS:\n%s\n\nRESULT IS EMPTY", o.Logs)
	}
	return "RESULT IS EMPTY"
}
