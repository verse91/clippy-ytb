package downloader

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

func FullVideoFHD(videoURL string) error {
	start := time.Now()

	// yt-dlp \
	// -f "bv[height=1080][vcodec^=avc1]+ba[ext=m4a]/bv[height=1080][vcodec^=avc1]" \
	// -S "+vbr,+abr" \
	// -o "<outputDir>/%(title)s (1080p, h264).%(ext)s" \
	// "<videoURL>"
	cmd_1080p := exec.Command(
		ytDlpPath,
		"-f", "bv[height=1080][vcodec^=avc1]+ba[ext=m4a]/bv[height=1080][vcodec^=avc1]",
		"-S", "+vbr,+abr", // sort bitrate
		"-o", filepath.Join(outputDir, "%(title)s (1080p, h264).%(ext)s"),
		videoURL,
	)
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd_1080p.Stdout = &stdoutBuf
	cmd_1080p.Stderr = &stderrBuf

	if err := cmd_1080p.Run(); err != nil {
		fmt.Println("Fail:", err)
		return err
	}

	// Combine both stdout and stderr for scanning, since yt-dlp may print to either
	combined := append(stdoutBuf.Bytes(), stderrBuf.Bytes()...)
	scanner := bufio.NewScanner(bytes.NewReader(combined))
	re := regexp.MustCompile(`\[Merger\] Merging formats into \"(.+)\"`)
	var mergedFile string
	for scanner.Scan() {
		line := scanner.Text()
		if matches := re.FindStringSubmatch(line); matches != nil {
			mergedFile = matches[1]
			break
		}
	}

	if mergedFile != "" {
		base := filepath.Base(mergedFile)
		// // Example: "Rick Astley - Never Gonna Give You Up (Official Video) (4K Remaster) (1080p, h264).mp4"
		// // Cut .mp4
		// ext := filepath.Ext(base)
		// title := base[:len(base)-len(ext)]
		// // Remove " (1080p, h264)"
		// reTitle := regexp.MustCompile(`^(.*) \(1080p, h264\)$`)
		// if m := reTitle.FindStringSubmatch(title); m != nil {
		//     title = m[1]
		// }
		fmt.Println("✅ Download sucessfully:", base)
		// fmt.Println("🎵 Video title:", title)
	} else {
		fmt.Println("♻️ Video is already downloaded.")
	}
	fmt.Println("Took:", time.Since(start))
	return nil
}

// func HD(videoURL string) {
// 	start := time.Now()

// 	// yt-dlp \
// 	// -f "bv[height=720][vcodec^=avc1]+ba[ext=m4a]/bv[height=720][vcodec^=avc1]" \
// 	// -S "+vbr,+abr" \
// 	// -o "<outputDir>/%(title)s (720p, h264).%(ext)s" \
// 	// "<videoURL>"
// 	cmd_720p := exec.Command(
// 		ytDlpPath,
// 		"-f", "bv[height=720][vcodec^=avc1]+ba[ext=m4a]/bv[height=720][vcodec^=avc1]",
// 		"-S", "+vbr,+abr", // sort bitrate
// 		"-o", filepath.Join(outputDir, "%(title)s (720p, h264).%(ext)s"),
// 		videoURL,
// 	)
// 	var stdoutBuf bytes.Buffer
// 	var stderrBuf bytes.Buffer
// 	cmd_720p.Stdout = &stdoutBuf
// 	cmd_720p.Stderr = &stderrBuf

// 	if err := cmd_720p.Run(); err != nil {
// 		fmt.Println("Fail:", err)
// 		return
// 	}

// 	// Combine both stdout and stderr for scanning, since yt-dlp may print to either
// 	combined := append(stdoutBuf.Bytes(), stderrBuf.Bytes()...)
// 	scanner := bufio.NewScanner(bytes.NewReader(combined))
// 	re := regexp.MustCompile(`\[Merger\] Merging formats into \"(.+)\"`)
// 	var mergedFile string
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if matches := re.FindStringSubmatch(line); matches != nil {
// 			mergedFile = matches[1]
// 			break
// 		}
// 	}

// 	if mergedFile != "" {
// 		base := filepath.Base(mergedFile)
// 		// // Example: "Rick Astley - Never Gonna Give You Up (Official Video) (4K Remaster) (720p, h264).mp4"
// 		// // Cut .mp4
// 		// ext := filepath.Ext(base)
// 		// title := base[:len(base)-len(ext)]
// 		// // Remove " (720p, h264)"
// 		// reTitle := regexp.MustCompile(`^(.*) \(720p, h264\)$`)
// 		// if m := reTitle.FindStringSubmatch(title); m != nil {
// 		//     title = m[1]
// 		// }
// 		fmt.Println("✅ Download sucessfully:", base)
// 		// fmt.Println("🎵 Video title:", title)
// 	} else {
// 		fmt.Println("♻️ Video is already downloaded.")
// 	}
// 	fmt.Println("Took:", time.Since(start))
// }
