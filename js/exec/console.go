package exec

import (
	"fmt"
	"strings"

	"github.com/fastschema/qjs"
)

type Console struct {
	sb strings.Builder
}

func (c *Console) Log(args ...*qjs.Value) {
	var sb strings.Builder

	sb.WriteString("··· LOG START ···\n")
	for _, arg := range args {
		jstr, err := arg.JSONStringify()
		if err != nil {
			sb.WriteString("Fail to log object: " + err.Error())
		} else {
			sb.WriteString(jstr)
		}
		sb.WriteRune('\n')
	}
	sb.WriteString("··· LOG END ···")

	log := sb.String()

	fmt.Println(log)
	c.sb.WriteString(log)
	c.sb.WriteRune('\n')
}

func (c *Console) Error(args ...*qjs.Value) {
	var sb strings.Builder

	sb.WriteString("··· ERROR LOG START ···\n")
	for _, arg := range args {
		jstr, err := arg.JSONStringify()
		if err != nil {
			sb.WriteString("Fail to log object: " + err.Error())
		} else {
			sb.WriteString(jstr)
		}
		sb.WriteRune('\n')
	}
	sb.WriteString("··· ERROR LOG END ···")

	log := sb.String()

	fmt.Println(log)
	c.sb.WriteString(log)
	c.sb.WriteRune('\n')
}

func (c *Console) Info(args ...*qjs.Value) {
	var sb strings.Builder

	sb.WriteString("··· INFO LOG START ···\n")
	for _, arg := range args {
		jstr, err := arg.JSONStringify()
		if err != nil {
			sb.WriteString("Fail to log object: " + err.Error())
		} else {
			sb.WriteString(jstr)
		}
		sb.WriteRune('\n')
	}
	sb.WriteString("··· INFO LOG END ···")

	log := sb.String()

	fmt.Println(log)
	c.sb.WriteString(log)
	c.sb.WriteRune('\n')
}

func (c *Console) Warn(args ...*qjs.Value) {
	var sb strings.Builder

	sb.WriteString("··· WARN LOG START ···\n")
	for _, arg := range args {
		jstr, err := arg.JSONStringify()
		if err != nil {
			sb.WriteString("Fail to log object: " + err.Error())
		} else {
			sb.WriteString(jstr)
		}
		sb.WriteRune('\n')
	}
	sb.WriteString("··· WARN LOG END ···")

	log := sb.String()

	fmt.Println(log)
	c.sb.WriteString(log)
	c.sb.WriteRune('\n')
}

func (c *Console) Debug(args ...*qjs.Value) {
	var sb strings.Builder

	sb.WriteString("··· DEBUG LOG START ···\n")
	for _, arg := range args {
		jstr, err := arg.JSONStringify()
		if err != nil {
			sb.WriteString("Fail to log object: " + err.Error())
		} else {
			sb.WriteString(jstr)
		}
		sb.WriteRune('\n')
	}
	sb.WriteString("··· DEBUG LOG END ···")

	log := sb.String()

	fmt.Println(log)
	c.sb.WriteString(log)
	c.sb.WriteRune('\n')
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
