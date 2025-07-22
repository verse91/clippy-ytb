package repo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/supabase-community/supabase-go"
)

type VideoRepo struct {
	client *supabase.Client
}

func NewVideoRepo(supaClient *supabase.Client) *VideoRepo {
	return &VideoRepo{
		client: supaClient,
	}
}

func (vr *VideoRepo) CreateDownloadRequest(id, videoURL string) error {
	data := map[string]interface{}{
		"url":    videoURL,
		"status": "processing",
	}

	// Debug logging
	fmt.Printf("Inserting data: %+v\n", data)

	result, count, err := vr.client.From("downloads").Insert(data, false, "", "", "").Execute()
	if err != nil {
		// Detailed error logging
		// fmt.Printf("CreateDownloadRequest detailed error: %+v\n", err)
		// fmt.Printf("Error type: %T\n", err)
		// fmt.Printf("Error string: %s\n", err.Error())

		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("download already exists")
		}
		return fmt.Errorf("insert error: %s", err.Error())
	}

	fmt.Printf("Insert successful - Result: %s, Count: %d\n", string(result), count)
	return nil
}

// Cập nhật trạng thái (thành công / thất bại)
func (vr *VideoRepo) UpdateDownloadStatus(id, status, message string) error {
	data := map[string]interface{}{
		"status":  status,
		"message": message,
	}
	_, _, err := vr.client.From("downloads").
		Update(data, "", "").
		Eq("id", id).
		Execute()
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}
	return nil
}

func (vr *VideoRepo) GetStatus(id string) (string, error) {
	resp, _, err := vr.client.From("downloads").
		Select("status", "", false).
		Eq("id", id).
		Single().
		Execute()
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", err
	}
	status, ok := result["status"].(string)
	if !ok {
		return "", fmt.Errorf("status not found or not a string")
	}
	return status, nil
}
