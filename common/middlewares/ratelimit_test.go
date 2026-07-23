package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)

	if !rl.allow("192.168.1.1") {
		t.Error("First request should be allowed")
	}

	if !rl.allow("192.168.1.1") {
		t.Error("Second request should be allowed")
	}

	if !rl.allow("192.168.1.1") {
		t.Error("Third request should be allowed")
	}

	if rl.allow("192.168.1.1") {
		t.Error("Fourth request should be blocked")
	}
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)

	if !rl.allow("192.168.1.1") {
		t.Error("First IP should be allowed")
	}

	if rl.allow("192.168.1.1") {
		t.Error("First IP should be blocked after limit")
	}

	if !rl.allow("192.168.1.2") {
		t.Error("Different IP should be allowed")
	}
}

func TestRateLimiter_WindowExpiry(t *testing.T) {
	rl := NewRateLimiter(1, 50*time.Millisecond)

	if !rl.allow("192.168.1.1") {
		t.Error("First request should be allowed")
	}

	if rl.allow("192.168.1.1") {
		t.Error("Second request should be blocked")
	}

	time.Sleep(100 * time.Millisecond)

	if !rl.allow("192.168.1.1") {
		t.Error("Request after window should be allowed")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	handler := RateLimit(2, time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req1 := httptest.NewRequest("POST", "/api/auth/login", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("First request: status = %v, want %v", w1.Code, http.StatusOK)
	}

	req2 := httptest.NewRequest("POST", "/api/auth/login", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Second request: status = %v, want %v", w2.Code, http.StatusOK)
	}

	req3 := httptest.NewRequest("POST", "/api/auth/login", nil)
	req3.RemoteAddr = "192.168.1.1:12345"
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, req3)

	if w3.Code != http.StatusTooManyRequests {
		t.Errorf("Third request: status = %v, want %v", w3.Code, http.StatusTooManyRequests)
	}
}

func TestNewRateLimiter_ZeroWindow(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("NewRateLimiter(1, 0) panicked as expected: %v", r)
			return
		}
	}()

	rl := NewRateLimiter(1, 0)
	if rl == nil {
		t.Fatal("NewRateLimiter returned nil")
	}

	if !rl.allow("127.0.0.1") {
		t.Error("First request should be allowed even with zero window")
	}

	if rl.allow("127.0.0.1") {
		t.Error("Second request should be blocked with limit 1")
	}
}
