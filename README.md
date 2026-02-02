# :clapper: yt-scr
![Made with Go](https://img.shields.io/badge/Made%20with-Go-00ADD8?style=flat&logo=go)

> **Screenshot** YouTube video at specific timestamp

:warning: **Early Development** – Flags and behavior may change

`yt-scr` is a lightweight command-line utility written in Go that captures a screenshot
from a YouTube video at a given timestamp. It uses [yt-dlp](https://github.com/yt-dlp/yt-dlp) 
to download a short fragment of the video and [ffmpeg](https://ffmpeg.org/) to extract the exact frame. 

## Features
- Accepts both `youtu.be/...` and `youtube.com/watch?...` links
- Supports timestamps in formats like `t=42`, `t=1m23s`, or `t=1h2m35s`
- Extracts an exact frame at the requested time
- Outputs frame as `<videoid>_at_<timestamp>.png` by default
- Configurable options via command‑line flags

## Requirements
- [Go](https://golang.org/) 1.21+
- [yt-dlp](https://github.com/yt-dlp/yt-dlp/releases) installed and available in `PATH`
- [ffmpeg](https://ffmpeg.org/) installed and available in `PATH`

## Installation
Clone the repository and build the binary:
```bash
git clone https://github.com/ernieputnic/yt-scr
cd yt-scr
go build ./cmd/yt-scr
```

## Usage
Basic usage with timestamp in URL:
```bash
./yt-scr "https://youtu.be/dQw4w9WgXcQ?t=42"
```

Override timestamp:
```bash
./yt-scr https://youtu.be/dQw4w9WgXcQ\?t\=42 -t 43 
```

Set timestamp, set custom filename, keep fragment — flags can be combined:
```bash
./yt-scr --url https://youtu.be/dQw4w9WgXcQ -t 1m23s -o frame.png -k
```

Show help:
```bash
./yt-scr -h
```

## Flags
- `--url` - YouTube video URL (required)
- `-t` - Timestamp (e.g., `42`, `1m23s`, `1h2m3s`) - optional if URL contains `t=`
- `-o` - Output filename - optional, defaults to `<videoid>_at_<timestamp>.png`
- `-k` - Keep video fragment file - optional, defaults to `<videoid>_at_<timestamp>_fragment.mp4`
- `-h` - Show help

## Example Output
```bash
$ ./yt-scr --url "https://youtu.be/dQw4w9WgXcQ?t=42"
Downloading fragment...
[yt-dlp output...]
Extracting frame...
[ffmpeg output...]
Frame captured successfully: dQw4w9WgXcQ_at_42.png
```

## License
MIT License. See [LICENSE](LICENSE) for details.



