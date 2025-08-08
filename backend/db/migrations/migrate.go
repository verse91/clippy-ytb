package migrations

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/verse91/ytb-clipy/backend/db"
)

// RunDatabaseMigrations executes SQL schema from file using GORM
func RunDatabaseMigrations() error {
	if db.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Read SQL schema file
	schemaPath := filepath.Join("db", "migrations", "schema.sql")
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Split SQL into individual statements
	statements := splitSQLStatements(string(schemaSQL))

	// Execute each statement using GORM
	for i, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		log.Printf("Executing statement %d/%d", i+1, len(statements))

		// Execute raw SQL using GORM
		err := db.DB.Exec(statement).Error
		if err != nil {
			// Log error but continue with other statements
			log.Printf("Warning: Failed to execute statement %d: %v", i+1, err)
			log.Printf("Statement: %s", statement)
		}
	}

	log.Println("Database migrations completed")
	return nil
}

// splitSQLStatements splits SQL file into individual statements
func splitSQLStatements(sql string) []string {
	var result []string
	var currentStatement strings.Builder
	inDollarQuote := false

	lines := strings.Split(sql, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Skip comments
		if strings.HasPrefix(trimmedLine, "--") {
			continue
		}

		// Handle dollar-quoted strings
		if strings.Contains(line, "$$") {
			if !inDollarQuote {
				// Start of dollar quote
				parts := strings.SplitN(line, "$$", 2)
				if len(parts) >= 2 {
					inDollarQuote = true
					currentStatement.WriteString(line + "\n")
				}
			} else {
				// End of dollar quote
				if strings.Contains(line, "$$") {
					inDollarQuote = false
					currentStatement.WriteString(line + "\n")
				} else {
					currentStatement.WriteString(line + "\n")
				}
			}
			continue
		}

		if inDollarQuote {
			currentStatement.WriteString(line + "\n")
			continue
		}

		// Check for semicolon (end of statement)
		if strings.Contains(line, ";") {
			currentStatement.WriteString(line)
			stmt := strings.TrimSpace(currentStatement.String())
			if stmt != "" {
				result = append(result, stmt)
			}
			currentStatement.Reset()
		} else {
			currentStatement.WriteString(line + "\n")
		}
	}

	// Add any remaining statement
	if currentStatement.Len() > 0 {
		stmt := strings.TrimSpace(currentStatement.String())
		if stmt != "" {
			result = append(result, stmt)
		}
	}

	return result
}
