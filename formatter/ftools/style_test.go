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

func TestFormatStyleTagsUnderscoreBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Formats standalone underscore emphasis",
			input:    "say _hello_ now",
			expected: "say <i>hello</i> now",
		},
		{
			name:     "Does not format identifiers",
			input:    "ak_arch ak_arch_frontend",
			expected: "ak_arch ak_arch_frontend",
		},
		{
			name:     "Keeps punctuation outside emphasis",
			input:    "(_hello_),",
			expected: "(<i>hello</i>),",
		},
		{
			name:     "Unmatched underscore stays literal",
			input:    "_hello world",
			expected: "_hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatStyleTags(tt.input, "_", "_", "<i>", "</i>")
			if got != tt.expected {
				t.Errorf("FormatStyleTags underscore = %q, want %q", got, tt.expected)
			}
		})
	}
}
