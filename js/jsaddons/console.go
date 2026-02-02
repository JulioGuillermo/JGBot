package jsaddons

import (
	"fmt"
	"strings"

	"github.com/fastschema/qjs"
)

type Console struct {
	sb strings.Builder
}

func (c *Console) print(text string) {
	fmt.Println(text)
	c.sb.WriteString(text)
	c.sb.WriteRune('\n')
}

func (c *Console) newLog(start, end string, args ...*qjs.Value) {
	var sb strings.Builder

	sb.WriteString(start)
	sb.WriteRune('\n')
	for _, arg := range args {
		jstr, err := arg.JSONStringify()
		if err != nil {
			sb.WriteString("Fail to log object: " + err.Error())
		} else {
			sb.WriteString(jstr)
		}
		sb.WriteRune('\n')
	}
	sb.WriteString(end)

	c.print(sb.String())
}

func (c *Console) Log(args ...*qjs.Value) {
	c.newLog("··· LOG START ···", "··· LOG END ···", args...)
}

func (c *Console) Error(args ...*qjs.Value) {
	c.newLog("··· ERROR LOG START ···", "··· ERROR LOG END ···", args...)
}

func (c *Console) Info(args ...*qjs.Value) {
	c.newLog("··· INFO LOG START ···", "··· INFO LOG END ···", args...)
}

func (c *Console) Warn(args ...*qjs.Value) {
	c.newLog("··· WARN LOG START ···", "··· WARN LOG END ···", args...)
}

func (c *Console) Debug(args ...*qjs.Value) {
	c.newLog("··· DEBUG LOG START ···", "··· DEBUG LOG END ···", args...)
}

func (c *Console) GetJSObj(ctx *qjs.Context) (*qjs.Value, error) {
	logFun, err := qjs.ToJsValue(ctx, c.Log)
	if err != nil {
		return nil, err
	}
	errorFun, err := qjs.ToJsValue(ctx, c.Error)
	if err != nil {
		return nil, err
	}
	infoFun, err := qjs.ToJsValue(ctx, c.Info)
	if err != nil {
		return nil, err
	}
	warnFun, err := qjs.ToJsValue(ctx, c.Warn)
	if err != nil {
		return nil, err
	}
	debugFun, err := qjs.ToJsValue(ctx, c.Debug)
	if err != nil {
		return nil, err
	}

	console, err := qjs.ToJsValue(ctx, c)
	if err != nil {
		return nil, err
	}
	console.SetPropertyStr("log", logFun)
	console.SetPropertyStr("error", errorFun)
	console.SetPropertyStr("info", infoFun)
	console.SetPropertyStr("warn", warnFun)
	console.SetPropertyStr("debug", debugFun)

	return console, nil
}

func (c *Console) Clear() {
	c.sb.Reset()
}

func (c *Console) String() string {
	return c.sb.String()
}
