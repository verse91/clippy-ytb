package migrations

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/verse91/ytb-clipy/backend/db"
)

func RunDatabaseMigrations() error {
	if db.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	schemaPath := filepath.Join("db", "migrations", "schema.sql")
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	statements := splitSQLStatements(string(schemaSQL))

	for i, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		log.Printf("Executing statement %d/%d", i+1, len(statements))

		err := db.DB.Exec(statement).Error
		if err != nil {
			log.Printf("Warning: Failed to execute statement %d: %v", i+1, err)
			log.Printf("Statement: %s", statement)
		}
	}

	log.Println("Database migrations completed")
	return nil
}

func splitSQLStatements(sql string) []string {
	var result []string
	var currentStatement strings.Builder
	inDollarQuote := false

	lines := strings.Split(sql, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "--") {
			continue
		}

		if strings.Contains(line, "$$") {
			if !inDollarQuote {
				parts := strings.SplitN(line, "$$", 2)
				if len(parts) >= 2 {
					inDollarQuote = true
					currentStatement.WriteString(line + "\n")
				}
			} else {
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

	if currentStatement.Len() > 0 {
		stmt := strings.TrimSpace(currentStatement.String())
		if stmt != "" {
			result = append(result, stmt)
		}
	}

	return result
}
