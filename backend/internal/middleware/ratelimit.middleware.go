package middleware

import (
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

var clients = make(map[string]*Client)
var mu sync.Mutex

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
		limiter := rate.NewLimiter(5, 10)
		clients[ip] = &Client{limiter, time.Now()}
		return limiter
	}
	client.lastSeen = time.Now()
	return client.limiter
}

func CleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > time.Minute*3 {
				delete(clients, ip)
			}
		}
		mu.Unlock()
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
