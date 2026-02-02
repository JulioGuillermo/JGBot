package exec

import "github.com/fastschema/qjs"

func TypeModule() Option {
	return func(e *Executor) error {
		e.AddOption(qjs.TypeModule())
		return nil
	}
}

func FlagAsync() Option {
	return func(e *Executor) error {
		e.AddOption(qjs.FlagAsync())
		return nil
	}
}
