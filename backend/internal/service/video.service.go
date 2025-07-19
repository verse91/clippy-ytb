package service

import (
	"github.com/google/uuid"
    "github.com/verse91/ytb-clipy/backend/internal/repo"
	"github.com/verse91/ytb-clipy/backend/internal/video_pipeline/downloader"
)

type VideoService struct {
    VideoRepo *repo.VideoRepo
}

func NewVideoService() *VideoService {
    return &VideoService{
        VideoRepo: repo.NewVideoRepo(),
    }
}

func (vs *VideoService) DownloadVideo(videoURL string) string {
	id := uuid.New().String()
	downloader.FHD(videoURL)
	return id
}
