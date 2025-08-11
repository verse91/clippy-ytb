package downloader

import (
	// "fmt"
	"os"
	"path/filepath"
	"runtime"
)

var (
	ytDlpPath = getExecutablePath("yt-dlp")
	outputDir = getConfigValue("OUTPUT_DIR", "internal/video_pipeline/videos")
)

func getExecutablePath(name string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join("internal", "video_pipeline", "bin", name+".exe")
	}
	return filepath.Join("internal", "video_pipeline", "bin", name)
}

func getConfigValue(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// func download() {
//     var videoURL string

// 	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
//         fmt.Println("Cannot make a directory:", err)
// 	    return
// 	}

//     fmt.Print("Enter the video URL: ")
// 	fmt.Scanln(&videoURL)
//     FHD(videoURL)
// }
