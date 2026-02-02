package exec

import (
	"github.com/fastschema/qjs"
)

type Processor func(ctx *qjs.Context, jsResult *qjs.Value) (*qjs.Value, error)

func WithProcessor(processor Processor) Option {
	return func(e *Executor) error {
		e.AddProcessor(processor)
		return nil
	}
}

func WithAwait() Option {
	return WithProcessor(func(ctx *qjs.Context, jsResult *qjs.Value) (*qjs.Value, error) {
		if !jsResult.IsPromise() {
			return jsResult, nil
		}
		return jsResult.Await()
	})
}

func WithMainCall() Option {
	return WithProcessor(func(ctx *qjs.Context, jsResult *qjs.Value) (*qjs.Value, error) {
		if !jsResult.IsFunction() {
			return jsResult, nil
		}
		fn, err := qjs.JsFuncToGo[func() (*qjs.Value, error)](jsResult)
		if err != nil {
			return nil, err
		}
		return fn()
	})
}
