package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"

	// "log"
	// "fmt"
	// "github.com/goccy/go-json"
	router "github.com/verse91/ytb-clipy/backend/internal/routes"
	"github.com/verse91/ytb-clipy/backend/migrations"
)

// type RedirectConfig struct {
// 	Params  fiber.Map         // Route parameters
// 	Queries map[string]string // Query map
// }

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = migrations.AutoMigrate()
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	app := fiber.New()

	v1 := app.Group("/api/v1")

	// Setup all routes
	router.SetupRoutes(v1)

	if err := app.Listen(":" + os.Getenv("BACKEND_PORT")); err != nil {
		panic(err)
	}
}
