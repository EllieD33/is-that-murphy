package middleware_test

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellied33/is-that-murphy/middleware"
)

func TestRateLimiter(t *testing.T) {
	limiter := middleware.NewIPRateLimiter(2, 2)
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := limiter.Middleware(dummyHandler)

	req, _ := http.NewRequest("GET", "verify/", nil)
	req.RemoteAddr = "1.2.3.4:1234"

	rr := httptest.NewRecorder()

	// First 2 should pass
	handler.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK { t.Fatal("expected 200") }
    handler.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK { t.Fatal("expected 200") }

	// Third request should be rate limited
    rr = httptest.NewRecorder()
    handler.ServeHTTP(rr, req)
    if rr.Code != http.StatusTooManyRequests { t.Fatal("expected 429") }
}

func TestIPRateLimiter_SeparateLimitsPerIP(t *testing.T) {
	limiter := middleware.NewIPRateLimiter(1, 1)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := limiter.Middleware(testHandler)

	reqA := httptest.NewRequest("GET", "/verify", nil)
	reqA.RemoteAddr = net.JoinHostPort("10.0.0.1", "1234")

	reqB := httptest.NewRequest("GET", "/verify", nil)
	reqB.RemoteAddr = net.JoinHostPort("10.0.0.2", "5678")

	// First request from A should pass
	rrA1 := httptest.NewRecorder()
	handler.ServeHTTP(rrA1, reqA)
	if rrA1.Code != http.StatusOK {
		t.Fatalf("IP A first request should be OK, got %d", rrA1.Code)
	}

	// First request from B should also pass (different IP)
	rrB1 := httptest.NewRecorder()
	handler.ServeHTTP(rrB1, reqB)
	if rrB1.Code != http.StatusOK {
		t.Fatalf("IP B first request should be OK, got %d", rrB1.Code)
	}

	// Second request from A should rate-limit
	rrA2 := httptest.NewRecorder()
	handler.ServeHTTP(rrA2, reqA)
	if rrA2.Code != http.StatusTooManyRequests {
		t.Fatalf("IP A second request should be rate-limited, got %d", rrA2.Code)
	}
}
