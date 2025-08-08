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
