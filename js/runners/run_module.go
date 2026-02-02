package runners

import (
	"JGBot/js/exec"
	"JGBot/js/loader"
	"JGBot/js/result"
)

func RunModule(mainFile, path string, options ...exec.Option) (*result.Output, error) {
	codes, err := loader.LoadCode(path)
	if err != nil {
		return nil, err
	}
	return RunFiles(mainFile, codes, options...)
}
