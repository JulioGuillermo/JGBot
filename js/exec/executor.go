package exec

import (
	"JGBot/js/result"

	"github.com/fastschema/qjs"
)

type Executor struct {
	rt         *qjs.Runtime
	ctx        *qjs.Context
	jsResult   *qjs.Value
	console    Console
	processors []Processor
	options    []qjs.EvalOptionFunc
}

func NewExecutor() (*Executor, error) {
	exec := &Executor{
		processors: []Processor{},
		options:    []qjs.EvalOptionFunc{},
	}

	err := exec.initRuntime()
	if err != nil {
		return nil, err
	}

	err = exec.initConsole()
	if err != nil {
		return nil, err
	}

	return exec, nil
}

func (e *Executor) initRuntime() error {
	rt, err := qjs.New()
	if err != nil {
		return err
	}

	ctx := rt.Context()

	e.rt = rt
	e.ctx = ctx
	return nil
}

func (e *Executor) initConsole() error {
	e.console = Console{}
	jsConsole, err := e.console.GetJSObj(e.ctx)
	if err != nil {
		return err
	}
	e.ctx.Global().SetPropertyStr("console", jsConsole)
	e.ctx.SetFunc("print", func(ctx *qjs.This) (*qjs.Value, error) {
		e.console.Log(ctx.Args()...)
		return nil, nil
	})
	return nil
}

func (e *Executor) AddFunc(name string, fn qjs.Function) {
	e.ctx.SetFunc(name, fn)
}

func (e *Executor) AddAddonObj(addon JSAddonObj) {
	jsVal, err := addon.GetJSObj(e.ctx)
	if err != nil {
		panic(err)
	}
	e.ctx.Global().SetPropertyStr(addon.GetName(), jsVal)
}

func (e *Executor) AddVar(name string, val *qjs.Value) {
	e.ctx.Global().SetPropertyStr(name, val)
}

func (e *Executor) AddProcessor(processor Processor) {
	e.processors = append(e.processors, processor)
}

func (e *Executor) AddOption(option qjs.EvalOptionFunc) {
	e.options = append(e.options, option)
}

func (e *Executor) LoadModule(file, code string) error {
	_, err := e.ctx.Load(file, qjs.Code(code))
	return err
}

func (e *Executor) Run(file, code string) (err error) {
	e.jsResult, err = e.ctx.Eval(
		file,
		append(e.options, qjs.Code(code))...,
	)
	return
}

func (e *Executor) RunProcessors() (string, error) {
	result := e.jsResult
	var err error
	for _, processor := range e.processors {
		result, err = processor(e.ctx, result)
		if err != nil {
			return "", err
		}
	}
	return ValStr(result), nil
}

func (e *Executor) GetLogs() string {
	return e.console.String()
}

func (e *Executor) GetResult() (*result.Output, error) {
	return &result.Output{
		Logs:   e.console.String(),
		Result: ValStr(e.jsResult),
	}, nil
}

func (e *Executor) Close() {
	if e.jsResult != nil {
		e.jsResult.Free()
	}
	e.rt.Close()
}
