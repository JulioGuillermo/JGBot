package toolargs

import (
	"testing"
)

func TestExtractArg(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		ok       bool
		wantErr  bool
	}{
		// --- Valid JSON with __arg1 ---
		{
			name:     "Valid JSON with __arg1 string",
			input:    `{"__arg1": "content"}`,
			expected: "content",
			ok:       true,
			wantErr:  false,
		},
		{
			name:     "Valid JSON with __arg1 object",
			input:    `{"__arg1": {"key": "value"}}`,
			expected: `{"key":"value"}`,
			ok:       true,
			wantErr:  false,
		},
		{
			name:     "Valid JSON with __arg1 array",
			input:    `{"__arg1": [1, 2, 3]}`,
			expected: `[1,2,3]`,
			ok:       true,
			wantErr:  false,
		},
		{
			name:     "Valid JSON with __arg1 integer",
			input:    `{"__arg1": 123}`,
			expected: `123`,
			ok:       true,
			wantErr:  false,
		},
		{
			name:     "Valid JSON with __arg1 boolean",
			input:    `{"__arg1": true}`,
			expected: `true`,
			ok:       true,
			wantErr:  false,
		},
		{
			name:     "Valid JSON with __arg1 null",
			input:    `{"__arg1": null}`,
			expected: `null`,
			ok:       true,
			wantErr:  false,
		},
		{
			name:     "Valid JSON with __arg1 and other fields",
			input:    `{"other": "data", "__arg1": "content"}`,
			expected: "content",
			ok:       true,
			wantErr:  false,
		},

		// --- Valid JSON without __arg1 (should return as is) ---
		{
			name:     "Valid JSON object without __arg1",
			input:    `{"other": "value"}`,
			expected: `{"other": "value"}`,
			ok:       false,
			wantErr:  false,
		},
		{
			name:     "Valid JSON empty object",
			input:    `{}`,
			expected: `{}`,
			ok:       false,
			wantErr:  false,
		},
		{
			name:     "Valid JSON array",
			input:    `["item1", "item2"]`,
			expected: `["item1", "item2"]`,
			ok:       false,
			wantErr:  false,
		},

		// --- Non-JSON Strings (should return as is) ---
		{
			name:     "Plain string string",
			input:    "simple string",
			expected: "simple string",
			ok:       false,
			wantErr:  false,
		},
		{
			name:     "String with quotes but not JSON",
			input:    `"quoted string"`,
			expected: `"quoted string"`,
			ok:       false,
			wantErr:  false,
		},
		{
			name:     "String starting loosely resembling JSON but no braces",
			input:    `key: value`,
			expected: `key: value`,
			ok:       false,
			wantErr:  false,
		},
		{
			name:     "Multiline string",
			input:    "line1\nline2",
			expected: "line1\nline2",
			ok:       false,
			wantErr:  false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
			ok:       false,
			wantErr:  false,
		},
		{
			name:     "Whitespace only",
			input:    "   ",
			expected: "   ",
			ok:       false,
			wantErr:  false,
		},

		// --- Invalid JSON (starts with { or [ but is broken) ---
		{
			name:     "Invalid JSON object start",
			input:    `{"broken":`,
			expected: "",
			ok:       false,
			wantErr:  true,
		},
		{
			name:     "Invalid JSON array start",
			input:    `[1, 2,`,
			expected: "",
			ok:       false,
			wantErr:  true,
		},
		{
			name:     "Invalid JSON object with whitespace",
			input:    `  { "key": `,
			expected: "",
			ok:       false,
			wantErr:  true,
		},
		{
			name:     "Invalid JSON just brace",
			input:    `{`,
			expected: "",
			ok:       false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := ExtractArg(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if ok != tt.ok {
				t.Errorf("ExtractArg() ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.expected {
				t.Errorf("ExtractArg() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestExtractContent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		// --- Recursive Extraction ---
		{
			name:     "Single level extraction",
			input:    `{"__arg1": "content"}`,
			expected: "content",
			wantErr:  false,
		},
		{
			name:     "Two levels extraction",
			input:    `{"__arg1": {"__arg1": "final"}}`,
			expected: "final",
			wantErr:  false,
		},
		{
			name:     "Three levels extraction",
			input:    `{"__arg1": {"__arg1": {"__arg1": "deep"}}}`,
			expected: "deep",
			wantErr:  false,
		},
		{
			name:     "Recursive with mixed types",
			input:    `{"__arg1": {"__arg1": {"other": "value"}}}`,
			expected: `{"other":"value"}`,
			wantErr:  false,
		},

		// --- No Extraction Cases ---
		{
			name:     "No extraction (plain JSON)",
			input:    `{"foo": "bar"}`,
			expected: `{"foo": "bar"}`,
			wantErr:  false,
		},
		{
			name:     "No extraction (string)",
			input:    "raw string",
			expected: "raw string",
			wantErr:  false,
		},

		// --- Error Cases ---
		{
			name:     "Invalid JSON at first level",
			input:    `{"broken":`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Invalid JSON at second level",
			input:    `{"__arg1": invalid}`,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractContent(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ExtractContent() = %v, want %v", got, tt.expected)
			}
		})
	}
}
