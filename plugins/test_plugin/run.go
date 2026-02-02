package main

import (
	"JGBot/js/exec"
	"JGBot/js/jsaddons/httpaddon"
	"JGBot/js/runners"
	"fmt"
	"os"

	"github.com/fastschema/qjs"
)

func main() {
	fmt.Println(os.Getwd())

	output, err := runners.RunModule(
		"/init.js",
		"./plugins/test_plugin",
		httpaddon.WithHttp(),
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
