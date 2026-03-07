package test

import (
	"JGBot/js/runners"
	"strings"
	"testing"
)

func TestPrintNumberIsNotTruncated(t *testing.T) {
	output, err := runners.RunCode(`
		print(123456);
		print(123456);
		console.log(123456);
	`)
	if err != nil {
		t.Fatalf("RunCode() error = %v", err)
	}

	if strings.Count(output.Logs, "123456") != 3 {
		t.Fatalf("expected 3 full occurrences of 123456 in logs, got logs:\n%s", output.Logs)
	}
}
