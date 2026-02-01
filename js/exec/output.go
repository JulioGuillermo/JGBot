package exec

import "fmt"

type Output struct {
	Logs   string
	Result string
}

func (o *Output) String() string {
	return fmt.Sprintf("LOGS:\n\n%s\n\n\nRESULT:\n\n%s", o.Logs, o.Result)
}
