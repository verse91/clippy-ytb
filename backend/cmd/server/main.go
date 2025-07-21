package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"

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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Supabase client
	supabaseURL := os.Getenv("SUPABASE_DB_ENDPOINT")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_DB_ENDPOINT and SUPABASE_SERVICE_ROLE_KEY environment variables are required")
	}

	supaClient, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		log.Fatalf("Failed to create Supabase client: %v", err)
	}

	app := fiber.New()

	v1 := app.Group("/api/v1")

	// Setup all routes
	router.SetupRoutes(v1, supaClient)

	if err := app.Listen(":" + os.Getenv("BACKEND_PORT")); err != nil {
		panic(err)
	}
}
