package db

import (
	"context"
	"fmt"
	"time"

	"github.com/verse91/ytb-clipy/backend/internal/config"
	log "github.com/verse91/ytb-clipy/backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	dsn := config.LoadConfig().DBUrl
	if dsn == "" {
		return fmt.Errorf("database DSN is empty")
	}

	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}
	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// https://gorm.io/docs/generic_interface.html#Connection-Pool
	sqlDB, err := db.DB()
	if err != nil {
		// Close the GORM DB connection if we can't get the underlying sql.DB
		if closeDB, closeErr := db.DB(); closeErr == nil {
			if closeErr := closeDB.Close(); closeErr != nil {
				log.Log.Warn("Failed to close database connection", zap.Error(closeErr))
			}
		}
		return fmt.Errorf("error getting sql.db: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		if closeErr := sqlDB.Close(); closeErr != nil {
			log.Log.Warn("Failed to close database connection after ping failure", zap.Error(closeErr))
		}
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	log.Log.Info("Database connected successfully")

	return nil
}
