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
		_ = os.Setenv(key, expected)
		defer func() { _ = os.Unsetenv(key) }()
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
		_ = os.Setenv(key, "123")
		defer func() { _ = os.Unsetenv(key) }()
		val := getEnvInt(key, fallback)
		if val != expected {
			t.Errorf("got %v, want %v", val, expected)
		}
	})

	t.Run("returns fallback on invalid int", func(t *testing.T) {
		_ = os.Setenv(key, "not-an-int")
		defer func() { _ = os.Unsetenv(key) }()
		val := getEnvInt(key, fallback)
		if val != fallback {
			t.Errorf("got %v, want %v", val, fallback)
		}
	})
}
