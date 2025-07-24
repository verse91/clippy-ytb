// future use
package middleware

import (
    "github.com/gofiber/fiber/v3"
    // "github.com/verse91/ytb-clipy/backend/pkg/response"
)

// please make sure use reponse package, i was just lazy at this moment
func APIKeyMiddleware(c fiber.Ctx) error {
    apiKey := c.Get("X-API-Key")
    if apiKey == "" {
        c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "API Key is missing",
        })
    }

    if apiKey != "USER_API_KEY" {
        c.Status((fiber.StatusForbidden)).JSON(fiber.Map{
            "error": "Invalid API Key",
        })
    }
    return c.Next()
}
