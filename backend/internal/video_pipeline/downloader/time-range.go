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

var begin, end string

func TimeRangeFHD(videoURL string) {
	start := time.Now()
	fmt.Println("Iput begin: ")
	fmt.Scanf("%s\n", &begin)
	fmt.Println("Iput end: ")
	fmt.Scanf("%s\n", &end)
	fmt.Println("Begin, end:", begin, end)

	// yt-dlp \
	// -f "bv[height=1080][vcodec^=avc1]+ba[ext=m4a]/bv[height=1080][vcodec^=avc1]" \
	// -S "+vbr,+abr" \
	// -o "<outputDir>/%(title)s (1080p, h264).%(ext)s" \
	// "<videoURL>"
	cmd_1080p := exec.Command(
		ytDlpPath,
		"-f", "bv[height=1080][vcodec^=avc1]+ba[ext=m4a]/bv[height=1080][vcodec^=avc1]",
		"-S", "+vbr,+abr", // sort bitrate
		"--download-section", fmt.Sprintf("*%s-%s", begin, end),
		"-o", filepath.Join(outputDir, "%(title)s (1080p, h264).%(ext)s"),
		videoURL,
	)
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd_1080p.Stdout = &stdoutBuf
	cmd_1080p.Stderr = &stderrBuf

	if err := cmd_1080p.Run(); err != nil {
		fmt.Println("Fail:", err)
		return
	}

	// Combine both stdout and stderr for scanning, since yt-dlp may print to either
	combined := append(stdoutBuf.Bytes(), stderrBuf.Bytes()...)
	scanner := bufio.NewScanner(bytes.NewReader(combined))
	re := regexp.MustCompile(`\[download\] Destination: (.+)`)
	var downloaded string
	for scanner.Scan() {
		line := scanner.Text()
		if matches := re.FindStringSubmatch(line); matches != nil {
			downloaded = matches[1]
			break
		}
	}

	if downloaded == "" {
		fmt.Println("Video is already downloaded.")

	} else {
		base := filepath.Base(downloaded) // // Example: "Rick Astley - Never Gonna Give You Up (Official Video) (4K Remaster) (1080p, h264).mp4"
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
	}
	fmt.Println("Took:", time.Since(start))
}
