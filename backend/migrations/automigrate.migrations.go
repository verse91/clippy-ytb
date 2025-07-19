package migrations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func AutoMigrate() error {
	sql := `
		CREATE TABLE IF NOT EXISTS downloads (
			id UUID PRIMARY KEY,
			url TEXT NOT NULL,
			status TEXT,
			message TEXT,
			created_at TIMESTAMP DEFAULT now()
		);
	`

	// Validate required environment variables
	dbEndpoint := os.Getenv("SUPABASE_DB_ENDPOINT")
	if dbEndpoint == "" {
		return fmt.Errorf("SUPABASE_DB_ENDPOINT environment variable is required")
	}

	serviceRoleKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	if serviceRoleKey == "" {
		return fmt.Errorf("SUPABASE_SERVICE_ROLE_KEY environment variable is required")
	}

	// Use json.Marshal for proper JSON encoding
	requestBody := map[string]string{"query": sql}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(),
		"POST",
		dbEndpoint,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("apikey", serviceRoleKey)
	req.Header.Set("Authorization", "Bearer "+serviceRoleKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("migration failed (%d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}
