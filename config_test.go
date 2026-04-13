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
