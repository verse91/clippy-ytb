package service

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/verse91/ytb-clipy/backend/internal/video_pipeline/downloader"
)

// Status constants for download operations
const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

// Validation constants
const (
	MaxClipDurationSeconds = 3600 // 1 hour maximum clip duration
)

// VideoRepository interface defines the contract for video repository operations
type VideoRepository interface {
	CreateDownloadRequest(id, url string) (string, error)
	UpdateDownloadStatus(id, status, errorMsg string) error
	GetStatus(id string) (string, error)
	CreateTimeRangeDownloadRequest(id, url string, startSec, endSec int) error
	UpdateTimeRangeDownloadStatus(id, status, errorMsg, outputPath string) error
	GetTimeRangeDownloadStatus(id string) (map[string]interface{}, error)
}

type VideoService struct {
	VideoRepo VideoRepository
}

func NewVideoService(videoRepo VideoRepository) *VideoService {
	if videoRepo == nil {
		log.Fatal("VideoRepository cannot be nil")
	}
	return &VideoService{
		VideoRepo: videoRepo,
	}
}

func (vs *VideoService) validateURL(videoURL string) (string, error) {
	if videoURL == "" {
		return "", fmt.Errorf("video URL cannot be empty")
	}

	// Trim whitespace from input
	videoURL = strings.TrimSpace(videoURL)

	parsedURL, err := url.Parse(videoURL)
	if err != nil {
		return "", fmt.Errorf("invalid video URL format: %w", err)
	}

	// Check if URL has a scheme
	if parsedURL.Scheme == "" {
		// No scheme provided, prepend https://
		videoURL = "https://" + videoURL
		parsedURL, err = url.Parse(videoURL)
		if err != nil {
			return "", fmt.Errorf("invalid video URL format after adding https: %w", err)
		}
	} else if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		// Reject URLs with unsupported schemes
		return "", fmt.Errorf("unsupported URL scheme: %s. Only http and https are supported", parsedURL.Scheme)
	}

	// Ensure the URL has a non-empty host
	if parsedURL.Host == "" {
		return "", fmt.Errorf("invalid URL: missing host")
	}

	return parsedURL.String(), nil
}

func (vs *VideoService) DownloadFullVideo(videoURL string) (string, error) {
	validatedURL, err := vs.validateURL(videoURL)
	if err != nil {
		return "", fmt.Errorf("invalid video URL: %w", err)
	}

	// 🧩 Supabase tự sinh ID, ta nhận lại
	downloadID, err := vs.VideoRepo.CreateDownloadRequest("", validatedURL)
	if err != nil {
		return "", fmt.Errorf("failed to create download request: %w", err)
	}

	go func() {
		if err := downloader.FullVideoFHD(validatedURL); err != nil {
			vs.VideoRepo.UpdateDownloadStatus(downloadID, StatusFailed, err.Error())
			return
		}
		vs.VideoRepo.UpdateDownloadStatus(downloadID, StatusCompleted, "")
	}()

	return downloadID, nil
}

func (vs *VideoService) GetDownloadStatus(downloadID string) (string, error) {
	if downloadID == "" {
		return "", fmt.Errorf("download ID cannot be empty")
	}

	status, err := vs.VideoRepo.GetStatus(downloadID)
	if err != nil {
		return "", fmt.Errorf("failed to get download status: %w", err)
	}

	return status, nil
}

// Time Range Download Methods
func (vs *VideoService) DownloadVideoTimeRange(videoURL string, startSec, endSec int) (string, error) {
	validatedURL, err := vs.validateURL(videoURL)
	if err != nil {
		return "", fmt.Errorf("invalid video URL: %w", err)
	}

	// Validate time range with reasonable bounds
	if startSec < 0 {
		return "", fmt.Errorf("invalid time range: startSec must be >= 0")
	}
	if endSec <= startSec {
		return "", fmt.Errorf("invalid time range: endSec must be > startSec")
	}
	if endSec-startSec > MaxClipDurationSeconds {
		return "", fmt.Errorf("invalid time range: clip duration cannot exceed %d seconds", MaxClipDurationSeconds)
	}

	// Generate a temporary ID for tracking
	tempID := uuid.New().String()

	// Store the download request in repository
	if err := vs.VideoRepo.CreateTimeRangeDownloadRequest(tempID, validatedURL, startSec, endSec); err != nil {
		log.Printf("DownloadVideoTimeRange - CreateTimeRangeDownloadRequest error: %v", err)
		return "", fmt.Errorf("failed to create time range download request: %w", err)
	}

	// Start async download and processing
	go func() {
		if err := downloader.TimeRangeFHD(validatedURL, startSec, endSec, tempID); err != nil {
			// Update status in repository with error logging
			if updateErr := vs.VideoRepo.UpdateTimeRangeDownloadStatus(tempID, StatusFailed, err.Error(), ""); updateErr != nil {
				log.Printf("DownloadVideoTimeRange - UpdateTimeRangeDownloadStatus error: %v", updateErr)
			}
			return
		}
		if updateErr := vs.VideoRepo.UpdateTimeRangeDownloadStatus(tempID, StatusCompleted, "", ""); updateErr != nil {
			log.Printf("DownloadVideoTimeRange - UpdateTimeRangeDownloadStatus error: %v", updateErr)
		}
	}()

	return tempID, nil
}

func (vs *VideoService) GetTimeRangeDownloadStatus(downloadID string) (map[string]interface{}, error) {
	if downloadID == "" {
		return nil, fmt.Errorf("download ID cannot be empty")
	}

	status, err := vs.VideoRepo.GetTimeRangeDownloadStatus(downloadID)
	if err != nil {
		return nil, fmt.Errorf("failed to get time range download status: %w", err)
	}

	return status, nil
}
