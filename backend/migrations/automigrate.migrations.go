package migrations

import (
	"bytes"
	"context"
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

	body := []byte(fmt.Sprintf(`{"query": %q}`, sql))

	req, err := http.NewRequestWithContext(context.Background(),
		"POST",
		os.Getenv("SUPABASE_DB_ENDPOINT"), 
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("apikey", os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
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
