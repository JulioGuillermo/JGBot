package test

import (
	"JGBot/js/exec"
	"JGBot/js/runners"
	"testing"
)

func TestJS(t *testing.T) {
	output, err := runners.RunCode(`// Test JS
		// import { log } from "JGBot/js/log"
console.log("Hello World")
export default { a: 1 }
`,
		exec.TypeModule(),
	)

	if err != nil {
		t.Error(err)
		return
	}

	expectedLog := " - JS LOG >>> Hello World\n"
	if output.Logs != expectedLog {
		t.Errorf("Wrong logs\nEXPECTED:\n%s\nGOT:\n%s\n", expectedLog, output.Logs)
	}

	if output.Result != `{"a":1}` {
		t.Error("Wrong result")
	}
}
