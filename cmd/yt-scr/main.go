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
	url := flag.String("url", "", "YouTube video URL")
	ts := flag.String("t", "", "Timestamp e.g. 42, 1h23m57s (optional, derived from link)")
	out := flag.String("o", "", "Output filename PNG/JPG (optional)")
	kf := flag.Bool("k", false, "Keep video fragment file")
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
	}

	fmt.Printf("Frame captured successfully: %s\n", cfg.OutFile)

	if !cfg.KeepFrag {
		err := cfg.CleanupFragment()
		if err != nil {
			log.Printf("Warning: %v", err)
		}
	} else {
		fmt.Printf("Video fragment kept: %s\n", cfg.FragFile)
	}

}
