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
	// Check various proxy headers in order of preference
	headers := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"X-Client-IP",
		"CF-Connecting-IP", // Cloudflare
		"True-Client-IP",   // Akamai
	}

	for _, header := range headers {
		if ip := c.Get(header); ip != "" {
			// Handle multiple IPs in X-Forwarded-For (first one is the client)
			ips := strings.Split(ip, ",")
			if len(ips) > 0 {
				clientIP := strings.TrimSpace(ips[0])
				// Basic validation for IP format
				if isValidIP(clientIP) {
					return clientIP
				}
			}
		}
	}

	// Fallback to direct IP
	ip := c.IP()
	if isValidIP(ip) {
		return ip
	}

	// Last resort - return the IP as is
	return ip
}

func isValidIP(ip string) bool {
	// Basic validation for IPv4 and IPv6
	if ip == "" {
		return false
	}

	// Check for IPv4 format
	if strings.Contains(ip, ".") {
		parts := strings.Split(ip, ".")
		if len(parts) != 4 {
			return false
		}
		for _, part := range parts {
			if part == "" {
				return false
			}
		}
		return true
	}

	// Check for IPv6 format (basic check)
	if strings.Contains(ip, ":") {
		return len(ip) > 0
	}

	return false
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
	logger.Log.Info("Starting rate limiter cleanup routine",
		zap.Int("cleanup_interval_seconds", cleanupInterval),
		zap.Int("client_ttl_seconds", clientTTL),
	)

	ticker := time.NewTicker(time.Duration(cleanupInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Rate limiter cleanup routine stopped due to context cancellation")
			return
		case <-ticker.C:
			mu.Lock()
			removedCount := 0
			for ip, client := range clients {
				if time.Since(client.lastSeen) > time.Duration(clientTTL)*time.Second {
					delete(clients, ip)
					removedCount++
				}
			}
			mu.Unlock()

			if removedCount > 0 {
				logger.Log.Info("Cleaned up expired rate limiter clients",
					zap.Int("removed_count", removedCount),
					zap.Int("remaining_clients", len(clients)),
				)
			}
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
