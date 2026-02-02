package runners

import (
	"JGBot/js/exec"
	"JGBot/js/result"
)

func RunCode(code string, options ...exec.Option) (*result.Output, error) {
	exec, err := exec.NewExecutor()
	if err != nil {
		return nil, err
	}
	defer exec.Close()

	for _, option := range options {
		err = option(exec)
		if err != nil {
			return nil, err
		}
	}

	err = exec.Run("main.js", code)
	if err != nil {
		return nil, err
	}

	jsStr, err := exec.RunProcessors()

	output := &result.Output{
		Logs:   exec.GetLogs(),
		Result: jsStr,
	}

	return output, nil
}
