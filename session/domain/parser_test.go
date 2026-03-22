package sessiondomain

import (
	"testing"
)

func TestExtractReasoning(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantClean string
		wantThink string
	}{
		{
			name:      "simple case",
			input:     "<think>I'm thinking</think>The answer is 42",
			wantClean: "The answer is 42",
			wantThink: "I'm thinking",
		},
		{
			name:      "multiple thinking blocks",
			input:     "<think>Wait</think>Thinking again<think>Deeply</think>Answer",
			wantClean: "Thinking againAnswer",
			wantThink: "Wait\n\nDeeply",
		},
		{
			name:      "unclosed tag",
			input:     "<think>Streaming content",
			wantClean: "",
			wantThink: "Streaming content",
		},
		{
			name:      "no tags",
			input:     "Plain text message",
			wantClean: "Plain text message",
			wantThink: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClean, gotThink := ExtractReasoning(tt.input)
			if gotClean != tt.wantClean {
				t.Errorf("ExtractReasoning() clean = %v, want %v", gotClean, tt.wantClean)
			}
			if gotThink != tt.wantThink {
				t.Errorf("ExtractReasoning() think = %v, want %v", gotThink, tt.wantThink)
			}
		})
	}
}
