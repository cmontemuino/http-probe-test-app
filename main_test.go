package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func resetGlobalState() {
	readyMu.Lock()
	readyToggle = true
	readyMu.Unlock()
	atomic.StoreUint64(&reqCount, 0)
	startTime = time.Now()
}

func TestHealthzHandler(t *testing.T) {
	resetGlobalState()
	cfg := Config{}
	handler := healthzHandler(cfg)

	t.Run("default success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/healthz", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		if body := rr.Body.String(); body != "OK" {
			t.Errorf("handler returned unexpected body: got %v want %v", body, "OK")
		}
	})

	t.Run("fail override", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/healthz?fail=1", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})

	t.Run("fail after N requests", func(t *testing.T) {
		resetGlobalState()
		cfgN := Config{FailLivenessAfterN: 2}
		h := healthzHandler(cfgN)

		// First request - count 0
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest("GET", "/healthz", nil))
		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}

		// Simulate 2 requests to root
		atomic.StoreUint64(&reqCount, 2)

		rr = httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest("GET", "/healthz", nil))
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected 500 after 2 requests, got %d", rr.Code)
		}
	})
}

func TestReadyzHandler(t *testing.T) {
	resetGlobalState()
	cfg := Config{}
	handler := readyzHandler(cfg)

	t.Run("default ready", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/readyz", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("toggled off", func(t *testing.T) {
		readyMu.Lock()
		readyToggle = false
		readyMu.Unlock()

		req := httptest.NewRequest("GET", "/readyz", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusServiceUnavailable {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusServiceUnavailable)
		}
	})

	t.Run("startup delay", func(t *testing.T) {
		resetGlobalState()
		cfgDelay := Config{ReadyDelaySeconds: 10}
		h := readyzHandler(cfgDelay)

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest("GET", "/readyz", nil))
		if rr.Code != http.StatusServiceUnavailable {
			t.Errorf("expected 503 during delay, got %d", rr.Code)
		}
	})
}

func TestToggleReadyHandler(t *testing.T) {
	resetGlobalState()
	handler := toggleReadyHandler()

	t.Run("only POST allowed", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/toggle-ready", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", rr.Code)
		}
	})

	t.Run("toggle works", func(t *testing.T) {
		readyMu.Lock()
		readyToggle = true
		readyMu.Unlock()

		req := httptest.NewRequest("POST", "/toggle-ready", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		readyMu.RLock()
		isReady := readyToggle
		readyMu.RUnlock()
		if isReady != false {
			t.Error("readyToggle should be false after toggle")
		}
	})
}

func TestInfoHandler(t *testing.T) {
	resetGlobalState()
	Version = "1.2.3"
	GitCommit = "abc1234"
	cfg := Config{PodName: "test-pod", ClusterLabel: "test-cluster"}
	handler := infoHandler(cfg)

	req := httptest.NewRequest("GET", "/info", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp InfoResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal info response: %v", err)
	}

	if resp.Version != "1.2.3" || resp.GitCommit != "abc1234" || resp.PodName != "test-pod" {
		t.Errorf("unexpected info response content: %+v", resp)
	}
}

func TestRootHandler(t *testing.T) {
	resetGlobalState()
	registerMetrics("test")
	cfg := Config{PodName: "test-pod"}
	handler := rootHandler(cfg)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	if count := atomic.LoadUint64(&reqCount); count != 1 {
		t.Errorf("expected reqCount 1, got %d", count)
	}

	if !strings.Contains(rr.Body.String(), "Hello from test-pod") {
		t.Errorf("unexpected body: %s", rr.Body.String())
	}
}

func TestLoadConfig(t *testing.T) {
	_ = os.Setenv("PORT", "9090")
	_ = os.Setenv("PREFIX", "custom")
	defer func() { _ = os.Unsetenv("PORT") }()
	defer func() { _ = os.Setenv("PREFIX", "dummy") }()

	cfg := loadConfig()
	if cfg.Port != 9090 || cfg.Prefix != "custom" {
		t.Errorf("loadConfig failed: %+v", cfg)
	}
}
