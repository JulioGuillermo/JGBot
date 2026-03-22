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
				Respond: Respond{
					Always: true,
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
				Respond: Respond{
					Always: false,
					Match:  "^!",
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
				Respond: Respond{
					Always: true,
				},
			},
			message: "any message",
			want:    true,
		},
		{
			name: "never respond",
			cfg: &SessionConfiguration{
				Respond: Respond{
					Always: false,
					Match:  "",
				},
			},
			message: "any message",
			want:    false,
		},
		{
			name: "respond on prefix",
			cfg: &SessionConfiguration{
				Respond: Respond{
					Always: false,
					Match:  "^/",
				},
			},
			message: "/help",
			want:    true,
		},
		{
			name: "not respond on prefix",
			cfg: &SessionConfiguration{
				Respond: Respond{
					Always: false,
					Match:  "^/",
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
