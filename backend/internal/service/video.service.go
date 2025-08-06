package service

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/verse91/ytb-clipy/backend/internal/repo"
	"github.com/verse91/ytb-clipy/backend/internal/video_pipeline/downloader"
	"github.com/verse91/ytb-clipy/backend/pkg/logger"
	"go.uber.org/zap"
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
		videoURL = "https://" + videoURL
		parsedURL, err = url.Parse(videoURL)
		if err != nil {
			return "", fmt.Errorf("invalid video URL format after adding https: %w", err)
		}
	}
	return parsedURL.String(), nil
}

func (vs *VideoService) DownloadFullVideo(videoURL string) (string, error) {
	validatedURL, err := vs.validateURL(videoURL)
	if err != nil {
		return "", fmt.Errorf("invalid video URL: %w", err)
	}

	tempID := uuid.New().String()

	if err := vs.VideoRepo.CreateDownloadRequest(tempID, validatedURL); err != nil {
		logger.Log.Error("Failed to create download request",
			zap.Error(err),
			zap.String("downloadID", tempID),
			zap.String("videoURL", validatedURL))
		return "", fmt.Errorf("failed to create download request: %w", err)
	}

	logger.Log.Info("Starting full video download",
		zap.String("downloadID", tempID),
		zap.String("videoURL", validatedURL))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Error("Panic in full video download goroutine",
					zap.String("downloadID", tempID),
					zap.Any("panic", r))
				vs.VideoRepo.UpdateDownloadStatus(tempID, "failed", fmt.Sprintf("panic: %v", r))
			}
		}()

		select {
		case <-ctx.Done():
			logger.Log.Warn("Full video download cancelled before start",
				zap.String("downloadID", tempID),
				zap.Error(ctx.Err()))
			vs.VideoRepo.UpdateDownloadStatus(tempID, "failed", "download cancelled before start")
			return
		default:
		}

		resultChan := make(chan error, 1)

		go func() {
			err := downloader.FullVideoFHD(validatedURL)
			resultChan <- err
		}()

		select {
		case err := <-resultChan:
			if err != nil {
				logger.Log.Error("Full video download failed",
					zap.String("downloadID", tempID),
					zap.Error(err))
				vs.VideoRepo.UpdateDownloadStatus(tempID, "failed", err.Error())
				return
			}

			logger.Log.Info("Full video download completed successfully",
				zap.String("downloadID", tempID))
			vs.VideoRepo.UpdateDownloadStatus(tempID, "completed", "")

		case <-ctx.Done():
			logger.Log.Warn("Full video download cancelled due to timeout",
				zap.String("downloadID", tempID),
				zap.Error(ctx.Err()))
			vs.VideoRepo.UpdateDownloadStatus(tempID, "failed", "download timeout")
		}
	}()

	return tempID, nil
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

func (vs *VideoService) DownloadVideoTimeRange(videoURL string, startTime, endTime int) (string, error) {
	validatedURL, err := vs.validateURL(videoURL)
	if err != nil {
		return "", fmt.Errorf("invalid video URL: %w", err)
	}

	if startTime < 0 || endTime <= startTime {
		return "", fmt.Errorf("invalid time range: start_time must be >= 0 and end_time must be > start_time")
	}

	tempID := uuid.New().String()

	if err := vs.VideoRepo.CreateTimeRangeDownloadRequest(tempID, validatedURL, startTime, endTime); err != nil {
		logger.Log.Error("Failed to create time range download request",
			zap.Error(err),
			zap.String("downloadID", tempID),
			zap.String("videoURL", validatedURL),
			zap.Int("startTime", startTime),
			zap.Int("endTime", endTime))
		return "", fmt.Errorf("failed to create time range download request: %w", err)
	}

	logger.Log.Info("Starting time range download",
		zap.String("downloadID", tempID),
		zap.String("videoURL", validatedURL),
		zap.Int("startTime", startTime),
		zap.Int("endTime", endTime))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Error("Panic in time range download goroutine",
					zap.String("downloadID", tempID),
					zap.Any("panic", r))
				vs.VideoRepo.UpdateTimeRangeDownloadStatus(tempID, "failed", fmt.Sprintf("panic: %v", r), "")
			}
		}()

		select {
		case <-ctx.Done():
			logger.Log.Warn("Time range download cancelled before start",
				zap.String("downloadID", tempID),
				zap.Error(ctx.Err()))
			vs.VideoRepo.UpdateTimeRangeDownloadStatus(tempID, "failed", "download cancelled before start", "")
			return
		default:
		}

		resultChan := make(chan struct {
			err        error
			outputFile string
		}, 1)

		go func() {
			err := downloader.TimeRangeFHD(validatedURL, startTime, endTime, tempID)
			if err != nil {
				resultChan <- struct {
					err        error
					outputFile string
				}{err: err, outputFile: ""}
				return
			}

			outputFile := vs.extractOutputFileFromDownload(tempID, validatedURL, startTime, endTime)
			resultChan <- struct {
				err        error
				outputFile string
			}{err: nil, outputFile: outputFile}
		}()

		select {
		case result := <-resultChan:
			if result.err != nil {
				logger.Log.Error("Time range download failed",
					zap.String("downloadID", tempID),
					zap.Error(result.err))
				vs.VideoRepo.UpdateTimeRangeDownloadStatus(tempID, "failed", result.err.Error(), result.outputFile)
				return
			}

			logger.Log.Info("Time range download completed successfully",
				zap.String("downloadID", tempID),
				zap.String("outputFile", result.outputFile))
			vs.VideoRepo.UpdateTimeRangeDownloadStatus(tempID, "completed", "", result.outputFile)

		case <-ctx.Done():
			logger.Log.Warn("Time range download cancelled due to timeout",
				zap.String("downloadID", tempID),
				zap.Error(ctx.Err()))
			vs.VideoRepo.UpdateTimeRangeDownloadStatus(tempID, "failed", "download timeout", "")
		}
	}()

	return tempID, nil
}

func (vs *VideoService) extractOutputFileFromDownload(downloadID, videoURL string, startTime, endTime int) string {
	secondsToHHMMSS := func(sec int) string {
		h := sec / 3600
		m := (sec % 3600) / 60
		s := sec % 60
		return fmt.Sprintf("%02dh%02dm%02ds", h, m, s)
	}

	startStr := secondsToHHMMSS(startTime)
	endStr := secondsToHHMMSS(endTime)

	videoTitle := "video"
	if parsedURL, err := url.Parse(videoURL); err == nil {
		if videoID := vs.extractVideoID(parsedURL.String()); videoID != "" {
			videoTitle = videoID
		}
	}

	expectedFilename := fmt.Sprintf("%s_%s (%s-%s,1080p, h264).mp4", videoTitle, downloadID[:8], startStr, endStr)

	return filepath.Join("outputDir", expectedFilename)
}

func (vs *VideoService) extractVideoID(url string) string {
	patterns := []string{
		`(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/embed/)([a-zA-Z0-9_-]{11})`,
		`(?:youtube\.com/v/|youtube\.com/watch\?.*v=)([a-zA-Z0-9_-]{11})`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(url); len(matches) > 1 {
			return matches[1]
		}
	}

	return fmt.Sprintf("video_%d", len(url))
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
