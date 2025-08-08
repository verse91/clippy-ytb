package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/verse91/ytb-clipy/backend/db"
    "github.com/verse91/ytb-clipy/backend/pkg/logger"
    "go.uber.org/zap"
)

func RunDatabaseMigrations() error {
	if db.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Check if migration is needed by verifying database structure
	if isDatabaseUpToDate() {
		logger.Log.Info("Database is up to date")
		return nil
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

		err := db.DB.Exec(statement).Error
		if err != nil {
			// Only log if it's not an "already exists" error
			if !strings.Contains(err.Error(), "already exists") {
				logger.Log.Warn("Warning: Failed to execute statement", zap.Int("statement", i+1), zap.Error(err), zap.String("statement", statement))
			}
		}
	}

	logger.Log.Info("Database migrations completed")
	return nil
}

// isDatabaseUpToDate checks if all required tables exist
func isDatabaseUpToDate() bool {
	requiredTables := []string{"profiles", "downloads", "time_range_downloads"}

	for _, table := range requiredTables {
		var exists bool
		err := db.DB.Raw("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = ?)", table).Scan(&exists).Error

		if err != nil {
			logger.Log.Error("Error checking table", zap.String("table", table), zap.Error(err))
			return false
		}

		if !exists {
			logger.Log.Warn("Table does not exist, migration needed", zap.String("table", table))
			return false
		}
	}

	return true
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
