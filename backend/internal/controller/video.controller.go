package controller

import (
	"github.com/verse91/ytb-clipy/backend/internal/service"
	"github.com/verse91/ytb-clipy/backend/pkg/response"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

type VideoController struct {
	VideoService *service.VideoService
}

type VideoRequest struct {
	URL string `json:"url" binding:"required"`
}

func NewVideoController() *VideoController {
	return &VideoController{
		VideoService: service.NewVideoService(nil), // TODO: Pass proper VideoRepo instance
	}
}

func (vc *VideoController) DownloadHandler(c fiber.Ctx) error {
	var req VideoRequest

	if err := c.Bind().JSON(&req); err != nil {
		return response.ErrorResponse(c, 400001, "Invalid request body")
	}

	if req.URL == "" {
		return response.ErrorResponse(c, 400002, "URL is required")
	}

	downloadID, err := vc.VideoService.DownloadVideo(req.URL)
	if err != nil {
		return response.ErrorResponse(c, 500001, "Failed to start download: "+err.Error())
	}

	data := fiber.Map{
		"download_id": downloadID,
		"status":      "processing",
		"message":     "Download started successfully",
	}

	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return response.ErrorResponse(c, 500002, "Failed to serialize response")
	}

	return c.Status(fiber.StatusOK).Send(prettyJSON)
}

func (vc *VideoController) GetDownloadStatus(c fiber.Ctx) error {
	downloadID := c.Params("id")

	if downloadID == "" {
		return response.ErrorResponse(c, 400003, "Download ID is required")
	}

	status, err := vc.VideoService.GetDownloadStatus(downloadID)
	if err != nil {
		return response.ErrorResponse(c, 404001, "Download not found or failed to get status")
	}

	data := fiber.Map{
		"download_id": downloadID,
		"status":      status,
	}

	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return response.ErrorResponse(c, 500003, "Failed to serialize response")
	}

	return c.Status(fiber.StatusOK).Send(prettyJSON)
}
