package cli

import "testing"

func TestParseFlagsShowInterfaces(t *testing.T) {
	cfg, err := parseFlags([]string{"-si"})
	if err != nil {
		t.Fatalf("parseFlags returned error: %v", err)
	}

	if !cfg.ShowInterfaces {
		t.Fatal("expected ShowInterfaces to be true")
	}
	if cfg.Device != "" {
		t.Fatalf("expected empty device, got %q", cfg.Device)
	}
}

func TestParseFlagsDeviceAliases(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "short", args: []string{"-i", "en0"}},
		{name: "long", args: []string{"--interface", "en0"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := parseFlags(tt.args)
			if err != nil {
				t.Fatalf("parseFlags returned error: %v", err)
			}
			if cfg.Device != "en0" {
				t.Fatalf("expected device en0, got %q", cfg.Device)
			}
		})
	}
}

func TestParseFlagsRejectsUnknownFlag(t *testing.T) {
	if _, err := parseFlags([]string{"--unknown"}); err == nil {
		t.Fatal("expected unknown flag error")
	}
}
