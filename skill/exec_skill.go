package skill

import (
	"JGBot/js/exec"
	"JGBot/js/jsaddons/httpaddon"
	"JGBot/js/runners"
	"fmt"
	"path"

	"github.com/fastschema/qjs"
)

type SkillArgs map[string]string

func ExecSkillTool(name string, args SkillArgs) (string, error) {
	output, err := runners.RunModule(
		SkillToolFile,
		path.Join(SkillDir, name),
		httpaddon.WithHttp(),
		exec.TypeModule(),
		exec.FlagAsync(),
		exec.WithFunc("GetArgs", func(ctx *qjs.This) (*qjs.Value, error) {
			return qjs.ToJsValue(ctx.Context(), args)
		}),
		exec.WithMainCall(),
		exec.WithAwait(),
	)
	if err != nil {
		return "", err
	}

	fmt.Println(output.String())

	return output.Result, nil
}
