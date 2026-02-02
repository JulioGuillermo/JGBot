package exec

import "github.com/fastschema/qjs"

func WithFunc(name string, fn qjs.Function) Option {
	return func(e *Executor) error {
		e.AddFunc(name, fn)
		return nil
	}
}
