package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ernieputnic/yt-scr/internal/capture"
)

func main() {
	// Define flags
	url := flag.String("url", "", "YouTube video URL - required")
	ts := flag.String("t", "", "Timestamp (e.g. '42', '1m23s', '1h23m57s') - optional if URL contains 't='")
	out := flag.String("o", "", "Output filename PNG/JPG - optional, defaults to '<videoid>_at_<timestamp>.png'")
	kf := flag.Bool("k", false, "Keep video fragment file - optional, defaults to '<videoid>_at_<timestamp>_fragment.mp4'")
	help := flag.Bool("h", false, "Show help")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	// Validate required flags
	if *url == "" {
		fmt.Fprintln(os.Stderr, "Error: --url is required. Use -h for help.")
		os.Exit(1)
	}

	cfg, err := capture.NewCaptureConfig(*url, *ts, *out, *kf)
	if err != nil {
		log.Fatalf("error creating capture config: %v", err)
	}

	// Capture logic

	fmt.Println("Downloading fragment...")
	if err := cfg.DownloadFragment(); err != nil {
		log.Fatalf("yt-scr: failed to download fragment: %v", err)
	}

	fmt.Println("Extracting frame...")
	if err := cfg.ExtractFrame(); err != nil {
		log.Fatalf("ffmpeg: failed to capture frame: %v", err)
	} else if !cfg.KeepFrag {
		err := cfg.CleanupFragment()
		if err != nil {
			log.Printf("Warning: %v", err)
		}
	} else {
		fmt.Printf("Video fragment kept: %s\n", cfg.FragFile)
	}

	fmt.Printf("Frame captured successfully: %s\n", cfg.OutFile)

}
