package exec

import "github.com/fastschema/qjs"

func WithVar(name string, val *qjs.Value) Option {
	return func(e *Executor) error {
		e.AddVar(name, val)
		return nil
	}
}
