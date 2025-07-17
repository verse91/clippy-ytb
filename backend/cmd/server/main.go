package main

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"

	// "log"
	// "fmt"
	// "github.com/goccy/go-json"
	router "github.com/verse91/ytb-clipy/backend/internal/routes"
)

// type RedirectConfig struct {
// 	Params  fiber.Map         // Route parameters
// 	Queries map[string]string // Query map
// }

func main() {
	godotenv.Load()
	app := fiber.New()

	v1 := app.Group("/api/v1")

	// Setup all routes
	router.SetupRoutes(v1)

	if err := app.Listen(":" + os.Getenv("BACKEND_PORT")); err != nil {
		panic(err)
	}
}
