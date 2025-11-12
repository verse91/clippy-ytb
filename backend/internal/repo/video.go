package repo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/supabase-community/supabase-go"
)

type VideoRepo struct {
	client *supabase.Client
}

func NewVideoRepo(supaClient *supabase.Client) *VideoRepo {
	if supaClient == nil {
		panic("supabase client cannot be nil")
	}
	return &VideoRepo{
		client: supaClient,
	}
}

func (r *VideoRepo) CreateDownloadRequest(_ string, url string) (string, error) {
	data := map[string]interface{}{
		"url":     url,
		"status":  "processing",
		"message": nil,
	}

	var inserted []map[string]interface{}
	_, err := r.client.From("downloads").
		Insert(data, true, "", "", "").ExecuteTo(&inserted)
	if err != nil {
		return "", err
	}

	if len(inserted) == 0 {
		return "", fmt.Errorf("insert succeeded but returned no ID")
	}

	id, ok := inserted[0]["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("failed to parse returned ID")
	}

	return id, nil
}

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

// Time Range Download Methods
func (vr *VideoRepo) CreateTimeRangeDownloadRequest(id, videoURL string, startTime, endTime int) error {
	// Validate time range parameters
	if startTime < 0 || endTime < 0 {
		return fmt.Errorf("start time and end time must be non-negative")
	}
	if startTime >= endTime {
		return fmt.Errorf("start time must be less than end time")
	}

	data := map[string]interface{}{
		"url":        videoURL,
		"start_time": startTime,
		"end_time":   endTime,
		"status":     "processing",
	}

	_, _, err := vr.client.From("time_range_downloads").Insert(data, false, "", "", "").Execute()
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("time range download already exists")
		}
		return fmt.Errorf("insert error: %s", err.Error())
	}

	return nil
}

func (vr *VideoRepo) UpdateTimeRangeDownloadStatus(id, status, message, outputFile string) error {
	data := map[string]interface{}{
		"status":      status,
		"message":     message,
		"output_file": outputFile,
	}
	_, count, err := vr.client.From("time_range_downloads").
		Update(data, "", "").
		Eq("id", id).
		Execute()
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (vr *VideoRepo) GetTimeRangeDownloadStatus(id string) (map[string]interface{}, error) {
	resp, _, err := vr.client.From("time_range_downloads").
		Select("*", "", false).
		Eq("id", id).
		Single().
		Execute()
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}
