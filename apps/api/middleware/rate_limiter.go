package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type rateLimitEntry struct {
	count     int
	resetAt   time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	entries map[uuid.UUID]*rateLimitEntry
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[uuid.UUID]*rateLimitEntry),
		limit:   limit,
		window:  window,
	}

	// Cleanup goroutine — purge expired entries every 10 minutes
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	for k, v := range rl.entries {
		if now.After(v.resetAt) {
			delete(rl.entries, k)
		}
	}
}

// Middleware returns an echo middleware that rate limits by session ID
func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sessionVal := c.Get(SessionContextKey)
			if sessionVal == nil {
				return next(c)
			}

			sessionID := sessionVal.(uuid.UUID)

			rl.mu.Lock()
			entry, exists := rl.entries[sessionID]
			now := time.Now()

			if !exists || now.After(entry.resetAt) {
				rl.entries[sessionID] = &rateLimitEntry{
					count:   1,
					resetAt: now.Add(rl.window),
				}
				rl.mu.Unlock()
				return next(c)
			}

			if entry.count >= rl.limit {
				remaining := entry.resetAt.Sub(now)
				rl.mu.Unlock()
				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"success": false,
					"error": map[string]interface{}{
						"code":    "RATE_LIMIT_EXCEEDED",
						"message": "Melebihi batas 10 analisis per jam. Coba lagi dalam " + remaining.Round(time.Minute).String(),
					},
				})
			}

			entry.count++
			rl.mu.Unlock()
			return next(c)
		}
	}
}
