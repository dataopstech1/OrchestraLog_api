package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/orchestralog/api/pkg/apierror"
	"github.com/orchestralog/api/pkg/response"
)

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
}

type visitor struct {
	count    int
	windowAt time.Time
}

func newRateLimiter() *rateLimiter {
	rl := &rateLimiter{visitors: make(map[string]*visitor)}
	go rl.cleanup()
	return rl
}

func (rl *rateLimiter) allow(ip string, limit int, window time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, ok := rl.visitors[ip]
	if !ok || time.Since(v.windowAt) > window {
		rl.visitors[ip] = &visitor{count: 1, windowAt: time.Now()}
		return true
	}
	v.count++
	return v.count <= limit
}

func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.windowAt) > 10*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

var defaultLimiter = newRateLimiter()

func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if !defaultLimiter.allow(ip, limit, window) {
				response.Error(w, apierror.New(http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Too many requests"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
