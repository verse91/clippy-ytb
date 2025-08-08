package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/verse91/ytb-clipy/backend/db"
	"github.com/verse91/ytb-clipy/backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	// Current schema version - increment this when making schema changes
	CurrentSchemaVersion = 1
)

func RunDatabaseMigrations() error {
	if db.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Check if migration is needed by comparing schema versions
	if isDatabaseUpToDate() {
		logger.Log.Info("Database is up to date")
		return nil
	}

	// Get current schema version to determine what migrations to run
	currentVersion, err := GetCurrentSchemaVersion()
	if err != nil {
		// If schema_version table doesn't exist, start from version 0
		currentVersion = 0
	}

	logger.Log.Info("Starting database migration",
		zap.Int("current_version", currentVersion),
		zap.Int("target_version", CurrentSchemaVersion))

	// Run migrations incrementally
	err = runMigrations(currentVersion, CurrentSchemaVersion)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Log.Info("Database migrations completed", zap.Int("schema_version", CurrentSchemaVersion))
	return nil
}

// runMigrations runs migrations from currentVersion to targetVersion
func runMigrations(currentVersion, targetVersion int) error {
	// For now, we'll run the full schema since we're at version 1
	// In the future, this can be extended to run specific migration files
	if currentVersion < targetVersion {
		schemaPath := filepath.Join("db", "migrations", "schema.sql")
		schemaSQL, err := os.ReadFile(schemaPath)
		if err != nil {
			return fmt.Errorf("failed to read schema file: %w", err)
		}

		statements := splitSQLStatements(string(schemaSQL))

		// Begin transaction for atomic migration
		tx := db.DB.Begin()
		if tx.Error != nil {
			return fmt.Errorf("failed to begin transaction: %w", tx.Error)
		}

		// Defer rollback in case of error
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		for i, statement := range statements {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue
			}

			err := tx.Exec(statement).Error
			if err != nil {
				// Only log if it's not an "already exists" error
				if !strings.Contains(err.Error(), "already exists") {
					logger.Log.Warn("Warning: Failed to execute statement", zap.Int("statement", i+1), zap.Error(err), zap.String("statement", statement))
				}
			}
		}

		// Update schema version after successful migration
		err = updateSchemaVersionInTx(tx, targetVersion)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update schema version: %w", err)
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit migration transaction: %w", err)
		}
	}

	return nil
}

// isDatabaseUpToDate checks if the database schema version matches the expected version
func isDatabaseUpToDate() bool {
	// First, check if schema_version table exists
	var tableExists bool
	err := db.DB.Raw("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'schema_version')").Scan(&tableExists).Error
	if err != nil {
		logger.Log.Error("Error checking schema_version table", zap.Error(err))
		return false
	}

	if !tableExists {
		logger.Log.Info("Schema version table does not exist, migration needed")
		return false
	}

	// Check current schema version
	var currentVersion int
	err = db.DB.Raw("SELECT version FROM schema_version ORDER BY id DESC LIMIT 1").Scan(&currentVersion).Error
	if err != nil {
		logger.Log.Error("Error reading schema version", zap.Error(err))
		return false
	}

	logger.Log.Info("Schema version check",
		zap.Int("current_version", currentVersion),
		zap.Int("expected_version", CurrentSchemaVersion))

	if currentVersion < CurrentSchemaVersion {
		logger.Log.Info("Schema version mismatch, migration needed",
			zap.Int("current_version", currentVersion),
			zap.Int("expected_version", CurrentSchemaVersion))
		return false
	}

	return true
}

// updateSchemaVersion updates the schema version in the database
func updateSchemaVersion(version int) error {
	// Insert new schema version record
	err := db.DB.Exec("INSERT INTO schema_version (version, applied_at) VALUES (?, NOW())", version).Error
	if err != nil {
		return fmt.Errorf("failed to insert schema version: %w", err)
	}
	return nil
}

// updateSchemaVersionInTx updates the schema version within a transaction
func updateSchemaVersionInTx(tx *gorm.DB, version int) error {
	// Insert new schema version record
	err := tx.Exec("INSERT INTO schema_version (version, applied_at) VALUES (?, NOW())", version).Error
	if err != nil {
		return fmt.Errorf("failed to insert schema version: %w", err)
	}
	return nil
}

// GetCurrentSchemaVersion returns the current schema version from the database
func GetCurrentSchemaVersion() (int, error) {
	var version int
	err := db.DB.Raw("SELECT version FROM schema_version ORDER BY id DESC LIMIT 1").Scan(&version).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get current schema version: %w", err)
	}
	return version, nil
}

func splitSQLStatements(sql string) []string {
	var result []string
	var currentStatement strings.Builder

	// State tracking
	inSingleQuote := false
	inDoubleQuote := false
	inDollarQuote := false
	dollarTag := ""

	// Track position in the string
	pos := 0
	length := len(sql)

	for pos < length {
		char := sql[pos]

		// Handle single quotes
		if char == '\'' && !inDoubleQuote && !inDollarQuote {
			inSingleQuote = !inSingleQuote
			currentStatement.WriteByte(char)
			pos++
			continue
		}

		// Handle double quotes
		if char == '"' && !inSingleQuote && !inDollarQuote {
			inDoubleQuote = !inDoubleQuote
			currentStatement.WriteByte(char)
			pos++
			continue
		}

		// Handle dollar quotes
		if char == '$' && !inSingleQuote && !inDoubleQuote {
			if !inDollarQuote {
				// Start of dollar quote
				inDollarQuote = true
				currentStatement.WriteByte(char)
				pos++

				// Read the tag
				tagStart := pos
				for pos < length && sql[pos] != '$' {
					pos++
				}
				if pos < length {
					dollarTag = sql[tagStart:pos]
					currentStatement.WriteString(dollarTag)
					currentStatement.WriteByte('$')
					pos++
				}
			} else {
				// Check if this is the end of the dollar quote
				if pos+len(dollarTag)+1 < length && sql[pos:pos+len(dollarTag)+1] == dollarTag+"$" {
					currentStatement.WriteString(dollarTag)
					currentStatement.WriteByte('$')
					pos += len(dollarTag) + 1
					inDollarQuote = false
					dollarTag = ""
				} else {
					currentStatement.WriteByte(char)
					pos++
				}
			}
			continue
		}

		// Handle semicolons (statement separators)
		if char == ';' && !inSingleQuote && !inDoubleQuote && !inDollarQuote {
			currentStatement.WriteByte(char)
			stmt := strings.TrimSpace(currentStatement.String())
			if stmt != "" {
				result = append(result, stmt)
			}
			currentStatement.Reset()
			pos++
			continue
		}

		// Handle inline comments (--) when not in quotes
		if char == '-' && !inSingleQuote && !inDoubleQuote && !inDollarQuote {
			if pos+1 < length && sql[pos+1] == '-' {
				// Found inline comment, skip to end of line
				currentStatement.WriteString("--")
				pos += 2
				for pos < length && sql[pos] != '\n' {
					pos++
				}
				if pos < length {
					currentStatement.WriteByte('\n')
					pos++
				}
				continue
			}
		}

		// Handle block comments (/* */) when not in quotes
		if char == '/' && !inSingleQuote && !inDoubleQuote && !inDollarQuote {
			if pos+1 < length && sql[pos+1] == '*' {
				// Found start of block comment
				currentStatement.WriteString("/*")
				pos += 2
				for pos < length {
					if pos+1 < length && sql[pos] == '*' && sql[pos+1] == '/' {
						currentStatement.WriteString("*/")
						pos += 2
						break
					}
					currentStatement.WriteByte(sql[pos])
					pos++
				}
				continue
			}
		}

		// Write the current character
		currentStatement.WriteByte(char)
		pos++
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
