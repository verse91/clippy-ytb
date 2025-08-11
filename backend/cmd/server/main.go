package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"

	// "log"
	// "fmt"
	// "github.com/goccy/go-json"
	"github.com/verse91/ytb-clipy/backend/db"
	"github.com/verse91/ytb-clipy/backend/db/migrations"
	"github.com/verse91/ytb-clipy/backend/internal/config"
	"github.com/verse91/ytb-clipy/backend/internal/middleware"
	router "github.com/verse91/ytb-clipy/backend/internal/routes"
	"github.com/verse91/ytb-clipy/backend/pkg/logger"
	"github.com/verse91/ytb-clipy/backend/pkg/utils"
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

	// // Debug logging
	// log.Printf("Supabase URL: %s", supabaseURL)
	// log.Printf("Supabase Key length: %d", len(supabaseKey))

	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_DB_ENDPOINT and SUPABASE_SERVICE_ROLE_KEY environment variables are required")
	}

	supaClient, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		log.Fatalf("Failed to create Supabase client: %v", err)
	}
	logger.InitLogger()

	// Load configuration
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("Failed to load configuration")
	}

	// Initialize database and run migrations
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run database migrations
	if err := migrations.RunDatabaseMigrations(); err != nil {
		log.Printf("Warning: Database migrations failed: %v", err)
	}

	app := fiber.New()

	// Add CORS middleware
	allowedOrigins := utils.GetEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://127.0.0.1:3000")

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(allowedOrigins, ","),
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "X-User-ID", "X-Admin-Key"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	app.Use(middleware.RateLimitMiddleware)
	v1 := app.Group("/api/v1")

	// Setup all routes
	router.SetupRoutes(v1, supaClient, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go middleware.CleanupClients(ctx)

	if err := app.Listen(":" + utils.GetEnv("BACKEND_PORT", "8080")); err != nil {
		cancel()
		panic(err)
	}
	// Optionally wait for cleanup goroutine to finish if needed
}
