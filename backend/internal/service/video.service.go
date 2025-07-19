package service

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/verse91/ytb-clipy/backend/internal/repo"
	"github.com/verse91/ytb-clipy/backend/internal/video_pipeline/downloader"
)

type VideoService struct {
	VideoRepo *repo.VideoRepo
}

func NewVideoService(videoRepo *repo.VideoRepo) *VideoService {
	return &VideoService{
		VideoRepo: videoRepo,
	}
}
func (vs *VideoService) validateURL(videoURL string) (string, error) {
	if videoURL == "" {
		return "", fmt.Errorf("video URL cannot be empty")
	}

	parsedURL, err := url.Parse(videoURL)
	if err != nil {
		return "", fmt.Errorf("invalid video URL format: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		// If no scheme, prepend https:// and re-parse
		videoURL = "https://" + videoURL
		parsedURL, err = url.Parse(videoURL)
		if err != nil {
			return "", fmt.Errorf("invalid video URL format after adding https: %w", err)
		}
	}
	return parsedURL.String(), nil
}

func (vs *VideoService) DownloadVideo(videoURL string) (string, error) {
	validatedURL, err := vs.validateURL(videoURL)
	if err != nil {
		return "", fmt.Errorf("invalid video URL: %w", err)
	}

	id := uuid.New().String()

	// Store the download request in repository
	if err := vs.VideoRepo.CreateDownloadRequest(id, validatedURL); err != nil {
		return "", fmt.Errorf("failed to create download request: %w", err)
	}

	// Start async download
	go func() {
		if err := downloader.FHD(validatedURL); err != nil {
			// Update status in repository
			vs.VideoRepo.UpdateDownloadStatus(id, "failed", err.Error())
			return
		}
		vs.VideoRepo.UpdateDownloadStatus(id, "completed", "")
	}()

	return id, nil
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
