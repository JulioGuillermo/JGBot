package ftools

import "testing"

func TestProtectBold(t *testing.T) {
	input := "**bold** and __bold__"
	got := ProtectBold(input)
	// Expect placeholders replacing keys
	if got == input {
		t.Error("ProtectBold did not change input")
	}
	// Check content preserved
	if len(got) < len("bold") || got == "" {
		t.Error("ProtectBold content lost")
	}
}

func TestRestoreBold(t *testing.T) {
	// Simulate protected string
	input := BoldPlaceholder + "test" + BoldPlaceholder
	got := RestoreBold(input, "*")
	if got != "*test*" {
		t.Errorf("RestoreBold failed: got %q", got)
	}
}

func TestFormatItalic(t *testing.T) {
	input := "*italic* and _italic_"
	got := FormatItalic(input, "_")
	if got != "_italic_ and _italic_" {
		t.Errorf("FormatItalic failed: got %q", got)
	}
}

func TestFormatStrike(t *testing.T) {
	input := "~~strike~~"
	got := FormatStrike(input, "~")
	if got != "~strike~" {
		t.Errorf("FormatStrike failed: got %q", got)
	}
}
