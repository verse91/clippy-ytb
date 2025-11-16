package controller

import (
	"bufio"
	"fmt"
	"net/url"
	"time"

	"github.com/verse91/ytb-clipy/backend/internal/repo"
	"github.com/verse91/ytb-clipy/backend/internal/service"
	"github.com/verse91/ytb-clipy/backend/pkg/logger"
	"github.com/verse91/ytb-clipy/backend/pkg/response"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/supabase-community/supabase-go"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type VideoController struct {
	VideoService *service.VideoService
}

type VideoRequest struct {
	URL              string `json:"url" binding:"required"`
	AutoBlockSponsor bool   `json:"autoBlockSponsor"`
	ImportThumbnail  bool   `json:"importThumbnail"`
}

type TimeRangeVideoRequest struct {
	URL       string `json:"url" binding:"required"`
	StartTime int    `json:"start_time" binding:"required"`
	EndTime   int    `json:"end_time" binding:"required"`
}

func NewVideoController(supabaseClient *supabase.Client) *VideoController {
	videoRepo := repo.NewVideoRepo(supabaseClient)
	return &VideoController{
		VideoService: service.NewVideoService(videoRepo),
	}
}

func (vc *VideoController) DownloadHandler(c fiber.Ctx) error {
	var req VideoRequest

	// Debug logging
	// fmt.Printf("Request method: %s\n", c.Method())
	// fmt.Printf("Request path: %s\n", c.Path())
	// fmt.Printf("Request headers: %+v\n", c.GetReqHeaders())
	// fmt.Printf("Raw request body: %s\n", string(c.Body()))

	if err := c.Bind().JSON(&req); err != nil {
		logger.Log.Error("JSON bind error in download request",
			zap.Error(err),
			zap.String("handler", "DownloadHandler"),
		)
		return response.ErrorResponse(c, response.ErrInvalidRequestBody, fmt.Sprintf("Invalid request body: %v", err))
	}

	logger.Log.Info("Parsed download request",
		zap.String("url", req.URL),
		zap.String("handler", "DownloadHandler"),
	)

	if req.URL == "" {
		logger.Log.Warn("URL is empty after parsing",
			zap.String("handler", "DownloadHandler"),
		)
		return response.ErrorResponse(c, response.ErrURLRequired, "URL is required")
	}

	parsedURL, err := url.ParseRequestURI(req.URL)
	if err != nil {
		logger.Log.Warn("Invalid URL format",
			zap.String("url", req.URL),
			zap.Error(err),
			zap.String("handler", "DownloadHandler"),
		)
		return response.ErrorResponse(c, response.ErrInvalidRequestBody, "Invalid URL format")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		logger.Log.Warn("Unsupported URL scheme",
			zap.String("url", req.URL),
			zap.String("scheme", parsedURL.Scheme),
			zap.String("handler", "DownloadHandler"),
		)
		return response.ErrorResponse(c, response.ErrInvalidRequestBody, "URL must use HTTP or HTTPS scheme")
	}

	downloadID, err := vc.VideoService.DownloadFullVideo(req.URL, req.AutoBlockSponsor, req.ImportThumbnail)
	if err != nil {
		logger.Log.Error("Failed to start download",
			zap.Error(err),
			zap.String("url", req.URL),
			zap.String("handler", "DownloadHandler"),
		)
		return response.ErrorResponse(c, response.ErrDownloadStartFailed, "Failed to start download: "+err.Error())
	}

	data := fiber.Map{
		"download_id": downloadID,
		"status":      "processing",
		"message":     "Download started successfully",
	}

	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return response.ErrorResponse(c, response.ErrSerializeResponse, "Failed to serialize response")
	}

	return c.Status(fiber.StatusOK).Send(prettyJSON)
}

func (vc *VideoController) GetDownloadStatus(c fiber.Ctx) error {
	downloadID := c.Params("id")

	if downloadID == "" {
		return response.ErrorResponse(c, response.ErrDownloadIDRequired, "Download ID is required")
	}

	status, err := vc.VideoService.GetDownloadStatus(downloadID)
	if err != nil {
		return response.ErrorResponse(c, response.ErrDownloadNotFound, "Download not found or failed to get status")
	}

	data := fiber.Map{
		"download_id": downloadID,
		"status":      status,
	}

	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return response.ErrorResponse(c, response.ErrSerializeStatus, "Failed to serialize response")
	}

	return c.Status(fiber.StatusOK).Send(prettyJSON)
}

func (vc *VideoController) StreamDownloadStatus(c fiber.Ctx) error {
	downloadID := c.Params("id")
	if downloadID == "" {
		return response.ErrorResponse(c, response.ErrDownloadIDRequired, "Download ID is required")
	}

	fasthttpCtx := c.Context()

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	c.Set("X-Accel-Buffering", "no")

	return c.SendStream(fasthttp.NewStreamReader(func(w *bufio.Writer) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-fasthttpCtx.Done():
				logger.Log.Info("Client disconnected from SSE",
					zap.String("download_id", downloadID))
				return

			case <-ticker.C:
				status, err := vc.VideoService.GetDownloadStatus(downloadID)
				if err != nil {
					fmt.Fprintf(w, "data: {\"error\": \"failed to get status\"}\n\n")
					_ = w.Flush()
					continue
				}

				fmt.Fprintf(w, "data: %s\n\n", status)
				if err := w.Flush(); err != nil {
					logger.Log.Warn("Failed to flush SSE data", zap.Error(err))
					return
				}

				if status == "completed" || status == "failed" {
					logger.Log.Info("SSE stream ended",
						zap.String("download_id", downloadID),
						zap.String("status", status))
					return
				}
			}
		}
	}))
}

func (vc *VideoController) DownloadTimeRangeHandler(c fiber.Ctx) error {
	var req TimeRangeVideoRequest

	if err := c.Bind().JSON(&req); err != nil {
		logger.Log.Error("JSON bind error in time range download request",
			zap.Error(err),
			zap.String("handler", "DownloadTimeRangeHandler"),
		)
		return response.ErrorResponse(c, response.ErrInvalidRequestBody, fmt.Sprintf("Invalid request body: %v", err))
	}

	logger.Log.Info("Parsed time range download request",
		zap.String("url", req.URL),
		zap.Int("start_time", req.StartTime),
		zap.Int("end_time", req.EndTime),
		zap.String("handler", "DownloadTimeRangeHandler"),
	)

	if req.URL == "" {
		return response.ErrorResponse(c, response.ErrURLRequired, "URL is required")
	}

	// Validate URL format and scheme
	parsedURL, err := url.ParseRequestURI(req.URL)
	if err != nil {
		logger.Log.Warn("Invalid URL format",
			zap.String("url", req.URL),
			zap.Error(err),
			zap.String("handler", "DownloadTimeRangeHandler"),
		)
		return response.ErrorResponse(c, response.ErrInvalidRequestBody, "Invalid URL format")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		logger.Log.Warn("Unsupported URL scheme",
			zap.String("url", req.URL),
			zap.String("scheme", parsedURL.Scheme),
			zap.String("handler", "DownloadTimeRangeHandler"),
		)
		return response.ErrorResponse(c, response.ErrInvalidRequestBody, "URL must use HTTP or HTTPS scheme")
	}

	if req.StartTime < 0 || req.EndTime <= req.StartTime {
		return response.ErrorResponse(c, response.ErrInvalidRequestBody, "Invalid time range: start_time must be >= 0 and end_time must be > start_time")
	}

	downloadID, err := vc.VideoService.DownloadVideoTimeRange(req.URL, req.StartTime, req.EndTime)
	if err != nil {
		logger.Log.Error("Failed to start time range download",
			zap.Error(err),
			zap.String("url", req.URL),
			zap.Int("start_time", req.StartTime),
			zap.Int("end_time", req.EndTime),
			zap.String("handler", "DownloadTimeRangeHandler"),
		)
		return response.ErrorResponse(c, response.ErrDownloadStartFailed, "Failed to start time range download: "+err.Error())
	}

	data := fiber.Map{
		"download_id": downloadID,
		"status":      "processing",
		"message":     "Time range download started successfully",
		"start_time":  req.StartTime,
		"end_time":    req.EndTime,
	}

	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return response.ErrorResponse(c, response.ErrSerializeResponse, "Failed to serialize response")
	}

	return c.Status(fiber.StatusOK).Send(prettyJSON)
}

func (vc *VideoController) GetTimeRangeDownloadStatusHandler(c fiber.Ctx) error {
	downloadID := c.Params("id")

	if downloadID == "" {
		return response.ErrorResponse(c, response.ErrDownloadIDRequired, "Download ID is required")
	}

	status, err := vc.VideoService.GetTimeRangeDownloadStatus(downloadID)
	if err != nil {
		return response.ErrorResponse(c, response.ErrDownloadNotFound, "Time range download not found or failed to get status")
	}

	prettyJSON, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return response.ErrorResponse(c, response.ErrSerializeStatus, "Failed to serialize response")
	}

	return c.Status(fiber.StatusOK).Send(prettyJSON)
}
