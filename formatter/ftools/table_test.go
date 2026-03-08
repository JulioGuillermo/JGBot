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

func TestFormatTableMultipleRowsAndNonTableText(t *testing.T) {
	input := "Before\n| Name | Age |\n|---|---|\n| Bob | 30 |\n| Ana | 25 |\nAfter"
	got := FormatTable(input)
	want := "Before\n1. •••\n\t- **Name**: Bob\n\t- **Age**: 30\n2. •••\n\t- **Name**: Ana\n\t- **Age**: 25\nAfter"
	if got != want {
		t.Errorf("FormatTable() = %q, want %q", got, want)
	}
}
