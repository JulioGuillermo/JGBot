package ftools

import (
	"fmt"
	"testing"
)

func TestMapCodeBlocks(t *testing.T) {
	input := "Check `code` here"
	mapped, blocks := MapCodeBlocks(input)
	expectedPH := fmt.Sprintf("Check %s0 here", CodePlaceholder)
	if mapped != expectedPH {
		t.Errorf("got %q, want %q", mapped, expectedPH)
	}
	if len(blocks) != 1 || blocks[0] != "`code`" {
		t.Errorf("blocks mismatch: %v", blocks)
	}
}

func TestRestoreCodeBlocks(t *testing.T) {
	input := fmt.Sprintf("Check %s0 here", CodePlaceholder)
	blocks := []string{"`code`"}

	// Test standard restore
	got := RestoreCodeBlocks(input, blocks, nil)
	if got != "Check `code` here" {
		t.Errorf("Restore failed: got %q", got)
	}

	// Test with format func
	gotTyped := RestoreCodeBlocks(input, blocks, func(s string) string {
		return "```" + s[1:len(s)-1] + "```"
	})
	if gotTyped != "Check ```code``` here" {
		t.Errorf("Restore with func failed: got %q", gotTyped)
	}
}
