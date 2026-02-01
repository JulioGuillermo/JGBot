package exec

import (
	"github.com/fastschema/qjs"
)

func Exec(file, code string) (*Output, error) {
	console := Console{}

	rt, err := qjs.New()
	if err != nil {
		return nil, err
	}
	defer rt.Close()

	ctx := rt.Context()

	jsConsole, err := console.GetJSObj(ctx)
	if err != nil {
		return nil, err
	}
	ctx.Global().SetPropertyStr("console", jsConsole)

	result, err := ctx.Eval(file, qjs.Code(code), qjs.FlagAsync(), qjs.TypeModule())
	if err != nil {
		return nil, err
	}

	jsonOuput, err := result.JSONStringify()
	if err != nil {
		return nil, err
	}

	return &Output{
		Logs:   console.String(),
		Result: jsonOuput,
	}, nil
}
