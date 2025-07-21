package controller

import (
	"github.com/verse91/ytb-clipy/backend/internal/errors"
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

func NewVideoController(supabaseClient *supabase.Client) *VideoController {
	videoRepo := repo.NewVideoRepo(supabaseClient)
	return &VideoController{
		VideoService: service.NewVideoService(videoRepo),
	}
}

func (vc *VideoController) DownloadHandler(c fiber.Ctx) error {
	var req VideoRequest

	if err := c.Bind().JSON(&req); err != nil {
		return response.ErrorResponse(c, errors.ErrInvalidRequestBody, "Invalid request body")
	}

	if req.URL == "" {
		return response.ErrorResponse(c, errors.ErrURLRequired, "URL is required")
	}

	downloadID, err := vc.VideoService.DownloadVideo(req.URL)
	if err != nil {
		return response.ErrorResponse(c, errors.ErrDownloadStartFailed, "Failed to start download: "+err.Error())
	}

	data := fiber.Map{
		"download_id": downloadID,
		"status":      "processing",
		"message":     "Download started successfully",
	}

	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return response.ErrorResponse(c, errors.ErrSerializeResponse, "Failed to serialize response")
	}

	return c.Status(fiber.StatusOK).Send(prettyJSON)
}

func (vc *VideoController) GetDownloadStatus(c fiber.Ctx) error {
	downloadID := c.Params("id")

	if downloadID == "" {
		return response.ErrorResponse(c, errors.ErrDownloadIDRequired, "Download ID is required")
	}

	status, err := vc.VideoService.GetDownloadStatus(downloadID)
	if err != nil {
		return response.ErrorResponse(c, errors.ErrDownloadNotFound, "Download not found or failed to get status")
	}

	data := fiber.Map{
		"download_id": downloadID,
		"status":      status,
	}

	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return response.ErrorResponse(c, errors.ErrSerializeStatus, "Failed to serialize response")
	}

	return c.Status(fiber.StatusOK).Send(prettyJSON)
}
