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

// Tạo bản ghi download mới
func (vr *VideoRepo) CreateDownloadRequest(id, videoURL string) error {
	data := map[string]interface{}{
		"id":     id,
		"url":    videoURL,
		"status": "processing",
	}
	_, _, err := vr.client.From("downloads").Insert(data, false, "", "", "").Execute()
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("download already exists")
		}
		return fmt.Errorf("insert error: %w", err)
	}
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
