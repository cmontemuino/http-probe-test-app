package main

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	key := "TEST_ENV_VAR"
	fallback := "fallback"

	t.Run("returns fallback when not set", func(t *testing.T) {
		_ = os.Unsetenv(key)
		val := getEnv(key, fallback)
		if val != fallback {
			t.Errorf("got %v, want %v", val, fallback)
		}
	})

	t.Run("returns value when set", func(t *testing.T) {
		expected := "real-value"
		t.Setenv(key, expected)
		val := getEnv(key, fallback)
		if val != expected {
			t.Errorf("got %v, want %v", val, expected)
		}
	})
}

func TestShutdownTimeoutConfig(t *testing.T) {
	t.Run("default value", func(t *testing.T) {
		_ = os.Unsetenv("SHUTDOWN_TIMEOUT_SECONDS")
		cfg := loadConfig()
		if cfg.ShutdownTimeoutSec != 5 {
			t.Errorf("expected default 5, got %d", cfg.ShutdownTimeoutSec)
		}
	})

	t.Run("custom value", func(t *testing.T) {
		t.Setenv("SHUTDOWN_TIMEOUT_SECONDS", "15")
		cfg := loadConfig()
		if cfg.ShutdownTimeoutSec != 15 {
			t.Errorf("expected 15, got %d", cfg.ShutdownTimeoutSec)
		}
	})
}

func TestGetEnvInt(t *testing.T) {
	key := "TEST_ENV_INT"
	fallback := 42

	t.Run("returns fallback when not set", func(t *testing.T) {
		_ = os.Unsetenv(key)
		val := getEnvInt(key, fallback)
		if val != fallback {
			t.Errorf("got %v, want %v", val, fallback)
		}
	})

	t.Run("returns value when set", func(t *testing.T) {
		expected := 123
		t.Setenv(key, "123")
		val := getEnvInt(key, fallback)
		if val != expected {
			t.Errorf("got %v, want %v", val, expected)
		}
	})

	t.Run("returns fallback on invalid int", func(t *testing.T) {
		t.Setenv(key, "not-an-int")
		val := getEnvInt(key, fallback)
		if val != fallback {
			t.Errorf("got %v, want %v", val, fallback)
		}
	})
}
