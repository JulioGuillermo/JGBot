package ftools

import "testing"

func TestFormatHeader(t *testing.T) {
	tests := []struct {
		input string
		want  string
		ok    bool
	}{
		{"# Title", "1️⃣ Title", true},
		{"## Title", "2️⃣ Title", true},
		{"### Title", "3️⃣ Title", true},
		{"#### Title", "4️⃣ Title", true},
		{"##### Title", "5️⃣ Title", true},
		{"###### Title", "6️⃣ Title", true},
		{"Not Header", "Not Header", false},
	}

	for _, tt := range tests {
		got, ok := FormatHeader(tt.input)
		if ok != tt.ok || got != tt.want {
			t.Errorf("FormatHeader(%q) = %q, %v; want %q, %v", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}

func TestFormatHeaderHTML(t *testing.T) {
	tests := []struct {
		input string
		want  string
		ok    bool
	}{
		{"# Title", "<h1>Title</h1>", true},
		{"## Subtitle", "<h2>Subtitle</h2>", true},
		{"###### Tail", "<h6>Tail</h6>", true},
		{" # Not Header", " # Not Header", false},
	}

	for _, tt := range tests {
		got, ok := FormatHeaderHTML(tt.input)
		if ok != tt.ok || got != tt.want {
			t.Errorf("FormatHeaderHTML(%q) = %q, %v; want %q, %v", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}
