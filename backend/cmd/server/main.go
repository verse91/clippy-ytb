package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"

	// "log"
	// "fmt"
	// "github.com/goccy/go-json"
	"github.com/verse91/ytb-clipy/backend/internal/middleware"
	router "github.com/verse91/ytb-clipy/backend/internal/routes"
	"github.com/verse91/ytb-clipy/backend/pkg/logger"
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

	// Initialize Supabase client
	supabaseURL := os.Getenv("SUPABASE_DB_ENDPOINT")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	// Debug logging
	log.Printf("Supabase URL: %s", supabaseURL)
	log.Printf("Supabase Key length: %d", len(supabaseKey))

	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_DB_ENDPOINT and SUPABASE_SERVICE_ROLE_KEY environment variables are required")
	}

	supaClient, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		log.Fatalf("Failed to create Supabase client: %v", err)
	}
	logger.InitLogger()

	app := fiber.New()
	app.Use(middleware.RateLimitMiddleware)
	v1 := app.Group("/api/v1")

	// Setup all routes
	router.SetupRoutes(v1, supaClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go middleware.CleanupClients(ctx)

	if err := app.Listen(":" + os.Getenv("BACKEND_PORT")); err != nil {
		cancel()
		panic(err)
	}
	// Optionally wait for cleanup goroutine to finish if needed
}
