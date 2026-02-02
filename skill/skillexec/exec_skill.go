package skillexec

import (
	"JGBot/ctxs"
	"JGBot/files"
	"JGBot/js/exec"
	"JGBot/js/jsaddons/httpaddon"
	"JGBot/js/jsaddons/virtualfilesaddon"
	"JGBot/js/runners"
	"JGBot/skill"
	"fmt"
	"path"

	"github.com/fastschema/qjs"
)

type SkillArgs map[string]any

func ExecSkillTool(name string, args SkillArgs, rCtx *ctxs.RespondCtx) (string, error) {
	privateRoot, err := files.GetVirtualRoot(rCtx.SessionConf.Origin, name)
	if err != nil {
		return "", err
	}
	sharedRoot, err := files.GetVirtualRoot(rCtx.SessionConf.Origin, "Shared")
	if err != nil {
		return "", err
	}

	output, err := runners.RunModule(
		skill.SkillToolFile,
		path.Join(skill.SkillDir, name),
		httpaddon.WithHttp(),
		virtualfilesaddon.WithVirtualFile("VFPrivate", privateRoot),
		virtualfilesaddon.WithVirtualFile("VFShared", sharedRoot),
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

	fmt.Println(output.Result)
	return output.String(), nil
}
