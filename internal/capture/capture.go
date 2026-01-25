package capture

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

const offset int = 1

func ExtractVideoIdAndTs(raw string) (id, ts string, err error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", err
	}

	// Handle youtu.be short links
	if u.Host == "youtu.be" {
		id = strings.TrimPrefix(u.Path, "/")
		ts = u.Query().Get("t")
		return id, ts, nil
	}

	// Handle youtube.com links
	if strings.Contains(u.Host, "youtube.com") {
		switch {
		case strings.HasPrefix(u.Path, "watch"):
			id = u.Query().Get("v")
			ts = u.Query().Get("t")
		case strings.HasPrefix(u.Path, "embed"):
			id = path.Base(u.Path)
			ts = u.Query().Get("start")
		case strings.HasPrefix(u.Path, "shorts"):
			id = path.Base(u.Path)
			ts = u.Query().Get("t")
		}
	}

	ts = strings.TrimSuffix(ts, "s")

	return id, ts, nil
}

func ParseYoutubeTime(raw string) (int, error) {

	raw = strings.TrimSpace(strings.ToLower(raw))

	// if strings.HasSuffix(raw, "s") {
	// 	raw = strings.TrimSuffix(raw, "s")
	// }
	if cut, ok := strings.CutSuffix(raw, "s"); ok {
		raw = cut
	}

	// Plain int seconds
	if secs, err := strconv.Atoi(raw); err == nil {
		return secs, nil
	}

	// Composite time formats
	var total int
	var num string
	for _, ch := range raw {
		if ch >= '0' && ch <= '9' {
			num += string(ch)
			continue
		}
		if num == "" {
			continue
		}
		val, _ := strconv.Atoi(num)
		switch ch {
		case 'h':
			total += val * 3600
		case 'm':
			total += val * 60
		case 's':
			total += val
		}
		num = ""
	}
	if total == 0 {
		return 0, fmt.Errorf("invalid timestamp: %s", raw)
	}
	return total, nil
}

// Converts seconds into HH:MM:SS
func FormatFFmpegTime(secs int) string {
	return fmt.Sprintf("%02d:%02d:%02d", secs/3600, (secs%3600)/60, secs%60)
}

// Calculates start and finish for yt-dlp fragment download
func DeriveRange(ts, offset int) (start, finish int) {
	start = ts - offset
	if start < 0 {
		start = 0
	}
	finish = ts + offset
	return start, finish
}

// Downloads a short clip around the timestamp using yt-dlp
func GetFragment(url string, ts, offset int, out string) error {
	start, finish := DeriveRange(ts, offset)

	// yt-dlp command: download fragment between start and finish
	cmd := exec.Command(
		"yt-dlp", "-f", "bv", "--remux-video", "mp4",
		"--download-sections", fmt.Sprintf("*%d-%d", start, finish),
		"--force-overwrites", "--force-keyframes-at-cuts", "-o", out, url)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlp failed: %v", err)
	}

	return nil
}

// Extracts a single frame at the exact timestamp from the fragment
func GetFrame(url, tsRaw, fragFileName, out string) error {
	// Parse timestamp
	tsSecs, err := ParseYoutubeTime(tsRaw)
	if err != nil {
		return fmt.Errorf("invalid timestamp: %w", err)
	}

	// Download fragment

	if err := GetFragment(url, tsSecs, offset, fragFileName); err != nil {
		return fmt.Errorf("fragment extraction failed: %w", err)
	}

	// Extract frame at original ts
	frame := offset
	frameStr := FormatFFmpegTime(frame)

	cmd := exec.Command("ffmpeg", "-i", fragFileName, "-ss", frameStr, "-frames:v", "1", out)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %v", err)
	}

	return nil

}
