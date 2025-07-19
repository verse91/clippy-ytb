package downloader

import (
	"fmt"
	"os"
	"path/filepath"

)
var ytDlpPath = filepath.Join("bin", "yt-dlp.exe")
var outputDir = "videos"
func download() {
    var videoURL string

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
        fmt.Println("Cannot make a directory:", err)
	    return
	}

    fmt.Print("Enter the video URL: ")
	fmt.Scanln(&videoURL)
    FHD(videoURL)
}
