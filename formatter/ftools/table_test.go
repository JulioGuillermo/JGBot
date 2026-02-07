package ftools

import (
	"reflect"
	"testing"
)

func TestExtractTableCells(t *testing.T) {
	input := "| a | b |"
	got := ExtractTableCells(input)
	want := []string{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExtractTableCells failed: got %v, want %v", got, want)
	}
}

func TestProcessTableBlock(t *testing.T) {
	lines := []string{
		"| H1 | H2 |",
		"|---|---|",
		"| V1 | V2 |",
	}
	got := ProcessTableBlock(lines)
	// Expect:
	// 1.
	//   - *H1*: V1
	//   - *H2*: V2

	if len(got) != 3 {
		t.Errorf("ProcessTableBlock length mismatch: got %d lines, want 3", len(got))
	}
	if got[0] != "1. •••" {
		t.Errorf("Row index missing, got %q", got[0])
	}
}
