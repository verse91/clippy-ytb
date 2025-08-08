package db

import (
	"context"
	"fmt"
	"time"
	log "github.com/verse91/ytb-clipy/backend/pkg/logger"
	"github.com/verse91/ytb-clipy/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	// debug
	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
	db, err := gorm.Open(postgres.Open(config.LoadConfig().DBUrl), cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// https://gorm.io/docs/generic_interface.html#Connection-Pool
	sqlDB, err := db.DB()
    if err != nil {
        return fmt.Errorf("error getting sql.db: %w", err)
    }

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := sqlDB.PingContext(ctx); err != nil {
        sqlDB.Close()
        return fmt.Errorf("failed to ping database: %w", err)
    }
    log.Log.Info("Database connected successfully")

	return nil
}
