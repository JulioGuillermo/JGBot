package runners

import (
	"JGBot/js/exec"
	"JGBot/js/loader"
	"JGBot/js/result"
	"fmt"
)

func RunFiles(mainFile string, files []loader.Code, options ...exec.Option) (*result.Output, error) {
	code := loader.GetCode(files, mainFile)
	if code == nil {
		return nil, fmt.Errorf("Fail to find main file: %s", mainFile)
	}

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

	for _, code := range files {
		if code.Key == mainFile {
			continue
		}

		err = exec.LoadModule(code.Key, code.Code)
		if err != nil {
			return nil, err
		}
	}

	err = exec.Run(mainFile, code.Code)
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
