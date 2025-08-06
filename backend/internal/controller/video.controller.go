package controller

import (
	"fmt"

	"github.com/verse91/ytb-clipy/backend/internal/repo"
	"github.com/verse91/ytb-clipy/backend/internal/service"
	"github.com/verse91/ytb-clipy/backend/pkg/response"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/supabase-community/supabase-go"
)

type VideoController struct {
	VideoService *service.VideoService
}

type VideoRequest struct {
	URL string `json:"url" binding:"required"`
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
		fmt.Printf("JSON bind error: %v\n", err)
		return response.ErrorResponse(c, response.ErrInvalidRequestBody, fmt.Sprintf("Invalid request body: %v", err))
	}

	fmt.Printf("Parsed request: %+v\n", req)

	if req.URL == "" {
		fmt.Printf("URL is empty after parsing\n")
		return response.ErrorResponse(c, response.ErrURLRequired, "URL is required")
	}

	downloadID, err := vc.VideoService.DownloadFullVideo(req.URL)
	if err != nil {
		fmt.Printf("DownloadHandler error: %v\n", err)
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

func (vc *VideoController) DownloadTimeRangeHandler(c fiber.Ctx) error {
	var req TimeRangeVideoRequest

	if err := c.Bind().JSON(&req); err != nil {
		fmt.Printf("JSON bind error: %v\n", err)
		return response.ErrorResponse(c, response.ErrInvalidRequestBody, fmt.Sprintf("Invalid request body: %v", err))
	}

	fmt.Printf("Parsed time range request: %+v\n", req)

	if req.URL == "" {
		return response.ErrorResponse(c, response.ErrURLRequired, "URL is required")
	}

	if req.StartTime < 0 || req.EndTime <= req.StartTime {
		return response.ErrorResponse(c, response.ErrInvalidRequestBody, "Invalid time range: start_time must be >= 0 and end_time must be > start_time")
	}

	downloadID, err := vc.VideoService.DownloadVideoTimeRange(req.URL, req.StartTime, req.EndTime)
	if err != nil {
		fmt.Printf("DownloadTimeRangeHandler error: %v\n", err)
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
