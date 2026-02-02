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
		exec.WithFunc("onResult", func(ctx *qjs.This) (*qjs.Value, error) {
			fmt.Println("Result...", ctx.Args()[0])
			return nil, nil
		}),
		exec.WithMainCall(),
		exec.WithAwait(),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(output.String())
}
