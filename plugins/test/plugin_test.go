package test

import (
	"JGBot/js/exec"
	"JGBot/js/runners"
	"fmt"
	"os"
	"testing"

	"github.com/fastschema/qjs"
)

func TestPlugin(t *testing.T) {
	fmt.Println(os.Getwd())
	output, err := runners.RunModule(
		"/init.js",
		"../test_plugin",
		exec.WithFunc("onResult", func(ctx *qjs.This) (*qjs.Value, error) {
			fmt.Println("Result...", ctx.Args()[0])
			return nil, nil
		}),
		// exec.WithMainCall(),
		exec.WithAwait(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(output.String())
}
