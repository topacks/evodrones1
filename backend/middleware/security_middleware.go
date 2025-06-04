package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RateLimiter implementa um limitador de taxa simples
type RateLimiter struct {
	sync.Mutex
	requests map[string]int
	reset    time.Time
	limit    int
	window   time.Duration
}

// NewRateLimiter cria um novo limitador de taxa
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]int),
		limit:    limit,
		window:   window,
		reset:    time.Now().Add(window),
	}
}

// Allow verifica se uma requisição é permitida
func (rl *RateLimiter) Allow(ip string) bool {
	rl.Lock()
	defer rl.Unlock()

	now := time.Now()
	if now.After(rl.reset) {
		rl.requests = make(map[string]int)
		rl.reset = now.Add(rl.window)
	}

	rl.requests[ip]++
	return rl.requests[ip] <= rl.limit
}

var globalRateLimiter = NewRateLimiter(100, time.Minute)

// SecurityHeaders adiciona headers de segurança HTTP
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Content-Security-Policy", "default-src 'self';")
		c.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		return c.Next()
	}
}

// RateLimiterMiddleware implementa um limitador de taxa simples
func RateLimiterMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !globalRateLimiter.Allow(c.IP()) {
			return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
				"error": "limite de requisições excedido",
			})
		}
		return c.Next()
	}
}

// generateRequestID gera um ID único para a requisição
func generateRequestID() string {
	return uuid.New().String()
}

// RequestIDMiddleware adiciona um ID único para cada requisição
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Request-ID", generateRequestID())
		return c.Next()
	}
}

// SecurityMiddleware agrupa todos os middlewares de segurança
func SecurityMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Adiciona headers de segurança
		if err := SecurityHeaders()(c); err != nil {
			return err
		}

		// Adiciona rate limiting
		if err := RateLimiterMiddleware()(c); err != nil {
			return err
		}

		// Adiciona request ID
		if err := RequestIDMiddleware()(c); err != nil {
			return err
		}

		return c.Next()
	}
}
