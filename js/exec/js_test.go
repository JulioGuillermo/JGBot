package exec_test

import (
	"JGBot/js/exec"
	"testing"
)

func TestJS(t *testing.T) {
	output, err := exec.Exec("test.js", `// Test JS
console.log("Hello World")
export default { a: 1 }
`)

	if err != nil {
		t.Error(err)
	}

	if output.Logs != "··· LOG START ···\n\"Hello World\"\n··· LOG END ···\n" {
		t.Error("Wrong logs")
	}

	if output.Result != `{"a":1}` {
		t.Error("Wrong result")
	}
}
