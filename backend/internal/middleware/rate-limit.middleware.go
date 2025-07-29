package middleware

import (
	"context"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/verse91/ytb-clipy/backend/pkg/logger"
	"github.com/verse91/ytb-clipy/backend/pkg/response"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*Client)
	mu      sync.Mutex

	// Configurable rate limiter parameters
	requestsPerSecond = getEnvAsFloat("RATE_LIMIT_RPS", 5)
	burst             = getEnvAsInt("RATE_LIMIT_BURST", 10)
	cleanupInterval   = getEnvAsInt("RATE_LIMIT_CLEANUP_INTERVAL", 60) // seconds
	clientTTL         = getEnvAsInt("RATE_LIMIT_CLIENT_TTL", 180)      // seconds
)

func getEnvAsInt(name string, defaultVal int) int {
	if valStr := os.Getenv(name); valStr != "" {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}

func getEnvAsFloat(name string, defaultVal float64) float64 {
	if valStr := os.Getenv(name); valStr != "" {
		if val, err := strconv.ParseFloat(valStr, 64); err == nil {
			return val
		}
	}
	return defaultVal
}

func getClientIP(c fiber.Ctx) string {
	ip := c.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.IP()
	} else {
		ip = strings.Split(ip, ",")[0]
	}
	return strings.TrimSpace(ip)
}

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	client, exists := clients[ip]

	if !exists {
		limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burst)
		clients[ip] = &Client{limiter, time.Now()}
		return limiter
	}
	client.lastSeen = time.Now()
	return client.limiter
}

func CleanupClients(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(cleanupInterval) * time.Second):
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > time.Duration(clientTTL)*time.Second {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}
}

// ab -n 20 -c 1 http://localhost:8080/api/v1/download
// https://httpd.apache.org/docs/2.4/programs/ab.html
// windows user: https://www.apachelounge.com/download/
func RateLimitMiddleware(c fiber.Ctx) error {
	ip := getClientIP(c)
	limiter := getLimiter(ip)
	if !limiter.Allow() {
		logger.Log.Warn("Too many requests",
			zap.String("ip", ip),
			zap.Time("time", time.Now()),
		)
		return response.ErrorResponse(c, response.ErrTooManyRequests, "Too many requests. Please try again later.")
	}
	return c.Next()
}
