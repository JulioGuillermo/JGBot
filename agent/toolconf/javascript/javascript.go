package javascript

import (
	"JGBot/agent/tools"
	"JGBot/ctxs"
	"JGBot/js/exec"
	"JGBot/js/runners"
	"context"
	"fmt"
)

type JavaScriptArgs struct {
	Code string `json:"code" description:"The JavaScript code to execute."`
}

type JavaScriptInitializerConf struct{}

func (c *JavaScriptInitializerConf) Name() string {
	return "javascript"
}

func (c *JavaScriptInitializerConf) ToolInitializer(ctx *ctxs.RespondCtx) tools.Tool {
	return &tools.ToolAutoArgs[JavaScriptArgs]{
		ToolName:        c.Name(),
		ToolDescription: "Executes sandboxed JavaScript code (ES2023). Ideal for complex math, data parsing, or logic.",
		ToolFunc: func(ctx context.Context, args JavaScriptArgs) (string, error) {
			output, err := runners.RunCode(args.Code, exec.FlagAsync(), exec.WithAwait())
			if err != nil {
				return "", fmt.Errorf("ERROR: Fail to execute the JavaScript code: %s.", err.Error())
			}

			return fmt.Sprintf("SUCCESS: The JavaScript code executed successfully. The result is:\n\n%s", output.String()), nil
		},
	}
}
