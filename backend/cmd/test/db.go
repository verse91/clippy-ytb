package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/verse91/ytb-clipy/backend/db"
	"github.com/verse91/ytb-clipy/backend/pkg/logger"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize logger
	logger.InitLogger()

	// Initialize database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Test database connection
	log.Println("Testing database connection...")

	// Check if tables exist
	var tableCount int64
	db.DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableCount)
	log.Printf("Found %d tables in public schema", tableCount)

	// List all tables
	var tables []string
	db.DB.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name").Scan(&tables)

	log.Println("Tables found:")
	for _, table := range tables {
		log.Printf("  - %s", table)
	}

	// Check specific tables
	expectedTables := []string{"downloads", "time_range_downloads", "profiles"}
	for _, expectedTable := range expectedTables {
		var exists bool
		db.DB.Raw("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = ?)", expectedTable).Scan(&exists)
		if exists {
			log.Printf("✅ Table '%s' exists", expectedTable)
		} else {
			log.Printf("❌ Table '%s' does not exist", expectedTable)
		}
	}

	log.Println("Database test completed successfully!")
}
