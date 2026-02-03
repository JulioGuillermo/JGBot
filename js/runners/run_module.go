package runners

import (
	"JGBot/js/exec"
	"JGBot/js/loader"
	"JGBot/js/result"
)

func RunModule(mainFile, path string, options ...exec.Option) (*result.Output, error) {
	codes, err := loader.LoadCode(path, mainFile, false)
	if err != nil {
		return nil, err
	}
	// log.Info("Running module...", "entrypoint", mainFile, "path", path)
	return RunFiles(mainFile, codes, options...)
}

func RunModuleFetch(mainFile, path string, options ...exec.Option) (*result.Output, error) {
	codes, err := loader.LoadCode(path, mainFile, true)
	if err != nil {
		return nil, err
	}
	// log.Info("Running module (fetch)...", "entrypoint", mainFile, "path", path)
	// for _, code := range codes {
	// 	log.Info("Loaded module", "key", code.Key)
	// }
	return RunFiles(mainFile, codes, options...)
}
