package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type ipLimiter struct {
	clients map[string]*client
	mu      sync.Mutex
	r       rate.Limit
	b       int
}

func NewIPRateLimiter(r rate.Limit, b int) *ipLimiter {
	return &ipLimiter{
		clients: make(map[string]*client),
		r:       r,
		b:       b,
	}
}

func (i *ipLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	c, exists := i.clients[ip]
	if !exists {
		c = &client{
			limiter:  rate.NewLimiter(i.r, i.b),
			lastSeen: time.Now(),
		}
		i.clients[ip] = c
	} else {
		c.lastSeen = time.Now()
	}

	return c.limiter
}

func (i *ipLimiter) CleanupExpired(d time.Duration, stop <-chan struct{}) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			i.mu.Lock()
			for ip, c := range i.clients {
				if time.Since(c.lastSeen) > d {
					delete(i.clients, ip)
				}
			}
			i.mu.Unlock()
		case <-stop:
			return
		}
	}
}


func (i *ipLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "cannot determine IP", http.StatusInternalServerError)
			return
		}

		limiter := i.getLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
