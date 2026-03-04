package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Model != "flash2" {
		t.Errorf("default model = %s, want flash2", cfg.Model)
	}
	if cfg.Aspect != "16:9" {
		t.Errorf("default aspect = %s, want 16:9", cfg.Aspect)
	}
	if cfg.Size != "2k" {
		t.Errorf("default size = %s, want 2k", cfg.Size)
	}
}

func TestLoadEnvPriority(t *testing.T) {
	// IMGN_API_KEY takes priority over GEMINI_API_KEY
	os.Setenv("IMGN_API_KEY", "imgn-key")
	os.Setenv("GEMINI_API_KEY", "gemini-key")
	defer os.Unsetenv("IMGN_API_KEY")
	defer os.Unsetenv("GEMINI_API_KEY")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.APIKey() != "imgn-key" {
		t.Errorf("got %s, want imgn-key", cfg.APIKey())
	}
}

func TestLoadGeminiFallback(t *testing.T) {
	os.Unsetenv("IMGN_API_KEY")
	os.Setenv("GEMINI_API_KEY", "gemini-key")
	defer os.Unsetenv("GEMINI_API_KEY")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.APIKey() != "gemini-key" {
		t.Errorf("got %s, want gemini-key", cfg.APIKey())
	}
}
