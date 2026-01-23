# yt-scr
`yt-scr` is a lightweight command-line utility written in Go that captures a single frame
from a YouTube video at a given timestamp. It uses [yt-dlp](https://github.com/yt-dlp/yt-dlp) 
to download a short fragment of the video and [ffmpeg](https://ffmpeg.org/) to extract the exact frame. 

## Features
- Accepts both `youtu.be/...` and `youtube.com/watch?...` links
- Supports timestamps in formats like `t=34`, `t=1m34s`, or `t=1h2m34s`
- Extracts an exact frame at the requested time
- Outputs screenshot as `videoid_timestamp.png` by default
- Flags for custom filename, format, and verbosity

## Requirements
- [Go](https://golang.org/) 1.21+
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) installed and available in `PATH`
- [ffmpeg](https://ffmpeg.org/) installed and available in `PATH`

## Installation
Clone the repository and build the binary:
```
git clone https://github.com/ernieputnic/yt-scr
cd yt-scr
go build ./cmd/yt-scr
```
## Run
After building, you can run `yt-scr` with a YouTube link and timestamp:

```
./yt-scr https://youtu.be/dQw4w9WgXcQ?t=42
```
## License
MIT License. See [LICENSE](LICENSE) for details.



