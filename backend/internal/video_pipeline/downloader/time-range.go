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


func TimeRangeFHD(videoURL string, begin, end int, downloadID string) error {
	start := time.Now()
    secondsToHHMMSS := func(sec int) string {
        h := sec / 3600
        m := (sec % 3600) / 60
        s := sec % 60
        return fmt.Sprintf("%02dh%02dm%02ds", h, m, s)
    }
    beginInt := secondsToHHMMSS(begin)
    endInt := secondsToHHMMSS(end)
	// fmt.Println("Begin, end:", begin, end)

	// ../bin/yt-dlp.exe --no-playlist -f 'bv*[height<=1080][vcodec~=avc1]+ba*[ext=m4a]/bv*[height<=1080]+ba*[ext=m4a]/bv*+ba*/best[height<=1080]/best'  -S 'res:1080,+codec:avc1,+br' --download-sections '*30-90' -o 'outputDir/%(title)s (%(height)sp, h264).%(ext)s' 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
	cmd_1080p := exec.Command(
		ytDlpPath,
		"--no-playlist",
		"-f", `bv*[height<=1080][vcodec~=avc1]+ba*[ext=m4a]/bv*[height<=1080]+ba*[ext=m4a]/bv*+ba*/best[height<=1080]/best`,
		"-S", "res:1080,+codec:avc1,+br",
		"--download-section", fmt.Sprintf("*%d-%d", begin, end),
		"-o", filepath.Join(outputDir, fmt.Sprintf("%%(title)s (%s-%s,%%(height)sp, h264).%%(ext)s", beginInt, endInt)),
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
    return nil
}
