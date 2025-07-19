package controller

import (
	"github.com/verse91/ytb-clipy/backend/internal/service"
	"github.com/verse91/ytb-clipy/backend/pkg/response"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

type VideoController struct {
    VideoSerice *service.VideoService
}
type VideoRequest struct {
    URL string `json:"url" binding:"required"`
}

func NewVideoController() *VideoController {
    return &VideoController{
        VideoSerice: service.NewVideoService(),
    }
}


func (vc *VideoController) DownloadHandler(c fiber.Ctx) error {
