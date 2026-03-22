package sessiondomain

import (
	"testing"
)

func TestSessionConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cfg  *SessionConfiguration
	}{
		{
			name: "full config",
			cfg: &SessionConfiguration{
				Origin:      "telegram",
				Channel:     "general",
				ChatID:      "123:456",
				ChatName:    "Test Chat",
				HistorySize: 50,
				Admin:       "full",
				Allowed:     true,
				ShouldRespond: func(msg string) bool {
					return true
				},
			},
		},
		{
			name: "minimal config",
			cfg: &SessionConfiguration{
				Origin:      "whatsapp",
				Channel:     "chat",
				ChatID:      "999:888",
				ChatName:    "Minimal",
				HistorySize: 10,
				Admin:       "",
				Allowed:     false,
				ShouldRespond: func(msg string) bool {
					return false
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.cfg.Origin == "" {
				t.Error("Origin should not be empty")
			}
			if tt.cfg.ShouldRespond == nil {
				t.Error("ShouldRespond should not be nil")
			}
		})
	}
}

func TestSessionConfigurationShouldRespond(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     *SessionConfiguration
		message string
		want    bool
	}{
		{
			name: "always respond",
			cfg: &SessionConfiguration{
				ShouldRespond: func(msg string) bool {
					return true
				},
			},
			message: "any message",
			want:    true,
		},
		{
			name: "never respond",
			cfg: &SessionConfiguration{
				ShouldRespond: func(msg string) bool {
					return false
				},
			},
			message: "any message",
			want:    false,
		},
		{
			name: "respond on prefix",
			cfg: &SessionConfiguration{
				ShouldRespond: func(msg string) bool {
					return len(msg) > 0 && msg[0] == '/'
				},
			},
			message: "/help",
			want:    true,
		},
		{
			name: "not respond on prefix",
			cfg: &SessionConfiguration{
				ShouldRespond: func(msg string) bool {
					return len(msg) > 0 && msg[0] == '/'
				},
			},
			message: "hello",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.cfg.ShouldRespond(tt.message)
			if got != tt.want {
				t.Errorf("ShouldRespond(%q) = %v, want %v", tt.message, got, tt.want)
			}
		})
	}
}
