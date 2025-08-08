package config

import (
	"github.com/verse91/ytb-clipy/backend/pkg/logger"
	"github.com/verse91/ytb-clipy/backend/pkg/utils"
)

type Config struct {
    DBUrl string
}

func LoadConfig() *Config {
    db_url := utils.GetEnv("SUPABASE_DB_URL", "")
    if db_url == "" {
        logger.Log.Error("SUPABASE_DB_URL is not set")
        return nil
    }
    return &Config{
        DBUrl: db_url,
    }
}
