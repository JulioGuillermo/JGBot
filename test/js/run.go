package main

import (
	"JGBot/files"
	"JGBot/js/exec"
	"JGBot/js/jsaddons/httpaddon"
	"JGBot/js/jsaddons/virtualfilesaddon"
	"JGBot/js/runners"
	"fmt"
	"os"

	"github.com/fastschema/qjs"
)

func main() {
	fmt.Println(os.Getwd())

	vfRoot, err := files.GetVirtualRoot("test")
	if err != nil {
		panic(err)
	}

	output, err := runners.RunModule(
		"/init.js",
		"test/js",
		httpaddon.WithHttp(),
		virtualfilesaddon.WithVirtualFile("VF", vfRoot),
		exec.TypeModule(),
		exec.FlagAsync(),
		exec.WithFunc("GetArgs", func(ctx *qjs.This) (*qjs.Value, error) {
			return qjs.ToJsValue(ctx.Context(), map[string]any{
				"arg": "hi",
			})
		}),
		exec.WithMainCall(),
		exec.WithAwait(),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(output.String())
}
