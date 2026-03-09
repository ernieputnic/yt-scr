package capture

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

type FileExt string

const (
	JPG FileExt = "jpg"
	PNG FileExt = "png"
	MP4 FileExt = "mp4"
)

type internalConfig struct {
	offset             int
	allowedScrFileExts []FileExt
	scrFileExt         FileExt
	fragFileExt        FileExt
	showToolsOutput    bool
}

var defaults = &internalConfig{
	offset:             1,
	allowedScrFileExts: []FileExt{JPG, PNG},
	scrFileExt:         PNG,
	fragFileExt:        MP4,
	showToolsOutput:    false,
}

type CaptureConfig struct {
	URL       string // original YouTube URL
	VideoID   string // parsed from URL
	Timestamp int    // target timestamp in seconds
	Offset    int    // seconds before/after timestamp
	Start     int    // fragment start (clamped)
	Finish    int    // fragment finish
	FragFile  string // temporary fragment filename
	KeepFrag  bool   // whether to keep fragment file
	OutFile   string // final frame filename
}

func NewCaptureConfig(url, tsRaw, out string, kf bool) (*CaptureConfig, error) {

	id, urlTsRaw, err := ExtractVideoIDAndTimestamp(url)
	if err != nil {
		return nil, err
	}

	if tsRaw == "" {
		tsRaw = urlTsRaw
	}

	if tsRaw == "" {
		return nil, fmt.Errorf("no timestamp provided (flag or in URL)")
	}

	tsSecs, err := ParseYoutubeTime(tsRaw)
	if err != nil {
		return nil, err
	}

	start, finish := DeriveRange(tsSecs, defaults.offset)

	// Validate output path
	outExt := strings.TrimPrefix(filepath.Ext(out), ".")
	if out == "" {
		out = fmt.Sprintf("%s_at_%d.%s", id, tsSecs, defaults.scrFileExt)
	} else if outExt != "" {
		if !slices.Contains(defaults.allowedScrFileExts, FileExt(strings.ToLower(outExt))) {
			return nil, fmt.Errorf("unsupported file type used: .%s, allowed: %v", outExt, defaults.allowedScrFileExts)
		}
	} else {
		out = filepath.Join(out, fmt.Sprintf("%s_at_%d.%s", id, tsSecs, defaults.scrFileExt))
	}

	fragDir := filepath.Dir(out)
	fragFile := filepath.Join(fragDir, fmt.Sprintf("%s_at_%d_fragment.%s", id, tsSecs, defaults.fragFileExt))

	return &CaptureConfig{
		URL:       url,
		VideoID:   id,
		Timestamp: tsSecs,
		Offset:    defaults.offset,
		Start:     start,
		Finish:    finish,
		FragFile:  fragFile,
		KeepFrag:  kf,
		OutFile:   out,
	}, nil

}

func ExtractVideoIDAndTimestamp(raw string) (id, ts string, err error) {
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

// Checks provided timestamp format and converts to seconds if needed
func ParseYoutubeTime(raw string) (int, error) {

	raw = strings.TrimSpace(strings.ToLower(raw))

	// Plain integer seconds
	if secs, err := strconv.Atoi(raw); err == nil {
		return secs, nil
	}

	// Simple suffix: "42s"
	if strings.HasSuffix(raw, "s") && !strings.ContainsAny(raw, "hm") {
		if secs, err := strconv.Atoi(strings.TrimSuffix(raw, "s")); err == nil {
			return secs, nil
		}
	}

	// Composite time formats
	var total int
	var num string
	var lastUnit rune
	for _, ch := range raw {
		if ch >= '0' && ch <= '9' {
			num += string(ch)
			continue
		}
		if num == "" {
			// Non-digit without a preceding number
			if ch == 'h' || ch == 'm' || ch == 's' {
				return 0, fmt.Errorf("missing number before unit: %c", ch)
			}
		}
		val, _ := strconv.Atoi(num)
		switch ch {
		case 'h':
			if lastUnit >= 'h' {
				return 0, fmt.Errorf("duplicate or out-of-order unit hours")
			}
			total += val * 3600
			lastUnit = 'h'
		case 'm':
			if lastUnit >= 'm' {
				return 0, fmt.Errorf("duplicate or out-of-order minutes")
			}
			total += val * 60
			lastUnit = 'm'
		case 's':
			if lastUnit >= 's' {
				return 0, fmt.Errorf("duplicate or out-of-order seconds")
			}
			total += val
			lastUnit = 's'
		default:
			return 0, fmt.Errorf("invalid time unit: %c", ch)
		}
		num = ""
	}

	// Check for trailing numbers without unit
	if num != "" {
		return 0, fmt.Errorf("trailing number without unit: %s", num)
	}

	return total, nil
}

// Converts seconds into HH:MM:SS
func FormatFFmpegTime(secs int) string {
	return fmt.Sprintf("%02d:%02d:%02d", secs/3600, (secs%3600)/60, secs%60)
}

// Calculates start and finish for yt-dlp fragment download
func DeriveRange(ts, offset int) (start, finish int) {
	start = max(ts-offset, 0)
	finish = ts + offset
	return start, finish
}

// Downloads a short clip around the timestamp using yt-dlp
func (c CaptureConfig) DownloadFragment() error {

	// yt-dlp command: download fragment between start and finish
	cmd := exec.Command(
		"yt-dlp", "-f", "bv", "--remux-video", string(defaults.fragFileExt),
		"--download-sections", fmt.Sprintf("*%d-%d", c.Start, c.Finish),
		"--force-overwrites", "--force-keyframes-at-cuts", "-o", c.FragFile, c.URL)

	if defaults.showToolsOutput {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlp failed: %w", err)
	}
	return nil
}

// Extracts a single frame at the exact timestamp from the fragment
func (c CaptureConfig) ExtractFrame() error {

	frameSeek := c.Timestamp - c.Start
	frameStr := FormatFFmpegTime(frameSeek)

	cmd := exec.Command("ffmpeg", "-i", c.FragFile, "-ss", frameStr, "-frames:v", "1", c.OutFile)

	if defaults.showToolsOutput {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	return nil

}

// Cleans up fragment file
func (c CaptureConfig) CleanupFragment() error {
	if c.FragFile == "" {
		return fmt.Errorf("fragment filename not set")
	}
	if err := os.Remove(c.FragFile); err != nil {
		return fmt.Errorf("failed to remove fragment %s: %w", c.FragFile, err)
	}
	return nil
}
