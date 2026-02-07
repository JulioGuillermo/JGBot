package ftools

import (
	"testing"
)

func TestMapRestoreLinks(t *testing.T) {
	inputMD := "Check [Google](https://google.com)"
	inputRaw := "Visit https://example.com"

	// Test MD Link
	mappedMD, linksMD := MapLinks(inputMD)
	if len(linksMD) != 1 {
		t.Fatalf("MapLinks MD failed to capture link")
	}
	// Restore with format
	gotMD := RestoreLinks(mappedMD, linksMD)
	expectMD := "Check (Google: https://google.com)"
	if gotMD != expectMD {
		t.Errorf("RestoreLinks MD failed: got %q, want %q", gotMD, expectMD)
	}

	// Test Raw Link
	mappedRaw, linksRaw := MapLinks(inputRaw)
	if len(linksRaw) != 1 {
		t.Fatalf("MapLinks Raw failed to capture link")
	}
	// Restore raw (should be same)
	gotRaw := RestoreLinks(mappedRaw, linksRaw)
	if gotRaw != inputRaw {
		t.Errorf("RestoreLinks Raw failed: got %q, want %q", gotRaw, inputRaw)
	}
}
