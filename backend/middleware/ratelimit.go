package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// rateLimiter is a small, dependency-free fixed-window limiter keyed by client
// IP. It is intended for light abuse protection.
type rateLimiter struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	limit  int
	window time.Duration
}

// RateLimit returns middleware that allows at most `limit` requests per client
// IP within `window`, responding 429 once the limit is exceeded.
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	rl := &rateLimiter{
		hits:   make(map[string][]time.Time),
		limit:  limit,
		window: window,
	}
	go rl.cleanupLoop()

	return func(c *gin.Context) {
		if !rl.allow(c.ClientIP()) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests, please slow down"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (rl *rateLimiter) allow(key string) bool {
	now := time.Now()
	cutoff := now.Add(-rl.window)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	recent := rl.hits[key][:0]
	for _, t := range rl.hits[key] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}
	if len(recent) >= rl.limit {
		rl.hits[key] = recent
		return false
	}
	rl.hits[key] = append(recent, now)
	return true
}

// cleanupLoop periodically drops stale IP entries so the map does not grow
// unbounded over the lifetime of the process.
func (rl *rateLimiter) cleanupLoop() {
	ticker := time.NewTicker(10 * rl.window)
	defer ticker.Stop()
	for range ticker.C {
		cutoff := time.Now().Add(-rl.window)
		rl.mu.Lock()
		for key, times := range rl.hits {
			hasRecent := false
			for _, t := range times {
				if t.After(cutoff) {
					hasRecent = true
					break
				}
			}
			if !hasRecent {
				delete(rl.hits, key)
			}
		}
		rl.mu.Unlock()
	}
}
