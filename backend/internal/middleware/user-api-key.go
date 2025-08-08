// future use
package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/verse91/ytb-clipy/backend/pkg/response"
	"github.com/verse91/ytb-clipy/backend/pkg/utils"
)

// APIKeyMiddleware validates API key for general API access
func APIKeyMiddleware(c fiber.Ctx) error {
	apiKey := c.Get("X-API-Key")
	if apiKey == "" {
		return response.ErrorResponse(c, 401, "API Key is missing")
	}

	expectedAPIKey := utils.GetEnv("USER_API_KEY", "")
	if expectedAPIKey == "" {
		return response.ErrorResponse(c, 500, "USER_API_KEY environment variable not configured")
	}

	if apiKey != expectedAPIKey {
		return response.ErrorResponse(c, 403, "Invalid API Key")
	}
	return c.Next()
}

// AdminAuthMiddleware validates admin access for credit management
func AdminAuthMiddleware(c fiber.Ctx) error {
	// Check for admin API key
	adminKey := c.Get("X-Admin-Key")
	if adminKey == "" {
		return response.ErrorResponse(c, 401, "Admin authentication required")
	}

	// Get admin key from environment variable
	expectedAdminKey := utils.GetEnv("ADMIN_SECRET_KEY", "")
	if expectedAdminKey == "" {
		return response.ErrorResponse(c, 500, "Admin secret key not configured")
	}

	if adminKey != expectedAdminKey {
		return response.ErrorResponse(c, 403, "Invalid admin credentials")
	}

	return c.Next()
}

// UserAuthMiddleware validates user can only access their own data
func UserAuthMiddleware(c fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return response.ErrorResponse(c, 400, "User ID is required")
	}

	// Get JWT token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return response.ErrorResponse(c, 401, "Authorization header required")
	}

	// Extract token from "Bearer <token>" format
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return response.ErrorResponse(c, 401, "Invalid authorization header format")
	}

	tokenString := tokenParts[1]

	// Get JWT secret from environment
	jwtSecret := utils.GetEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return response.ErrorResponse(c, 500, "JWT secret not configured")
	}

	// Parse and validate JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return response.ErrorResponse(c, 401, "Invalid token")
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract user ID from claims
		authUserID, ok := claims["sub"].(string)
		if !ok {
			return response.ErrorResponse(c, 401, "Invalid token claims")
		}

		// User can only access their own data
		if authUserID != userID {
			return response.ErrorResponse(c, 403, "Access denied: can only access own data")
		}

		return c.Next()
	}

	return response.ErrorResponse(c, 401, "Invalid token")
}
