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
	gotMD := RestoreLinks(mappedMD, linksMD, false)
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
	gotRaw := RestoreLinks(mappedRaw, linksRaw, false)
	if gotRaw != inputRaw {
		t.Errorf("RestoreLinks Raw failed: got %q, want %q", gotRaw, inputRaw)
	}

	// Test MD Link with supportMD
	gotMDSupport := RestoreLinks(mappedMD, linksMD, true)
	expectMDSupport := "Check [Google](https://google.com)"
	if gotMDSupport != expectMDSupport {
		t.Errorf("RestoreLinks MD failed: got %q, want %q", gotMDSupport, expectMDSupport)
	}
}

func TestMapLinksMultipleAndUnderscoreContent(t *testing.T) {
	input := "Docs [ak_arch](https://example.com/ak_arch) and https://b.test/x_y"

	mapped, links := MapLinks(input)
	if len(links) != 2 {
		t.Fatalf("MapLinks failed to capture both links: %v", links)
	}

	gotHTML := RestoreLinksHTML(mapped, links)
	wantHTML := `Docs <a href="https://example.com/ak_arch">ak_arch</a> and <a href="https://b.test/x_y">https://b.test/x_y</a>`
	if gotHTML != wantHTML {
		t.Errorf("RestoreLinksHTML failed: got %q, want %q", gotHTML, wantHTML)
	}

	gotPlain := RestoreLinks(mapped, links, false)
	wantPlain := "Docs (ak_arch: https://example.com/ak_arch) and https://b.test/x_y"
	if gotPlain != wantPlain {
		t.Errorf("RestoreLinks failed: got %q, want %q", gotPlain, wantPlain)
	}
}
