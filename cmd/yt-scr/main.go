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
	ts := flag.String("t", "", "Timestamp (e.g. 42, 1m34s -- optional, derived from link)")
	out := flag.String("o", "", "Output filename. PNG/JPG (optional)")
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

	// Get YouTube video id and timestamp if present
	id, tsRaw, err := capture.ExtractVideoIdAndTs(*url)
	if err != nil {
		log.Fatalf("error parsing URL: %v", err)
	}
	// Set the output filename from Youtube link
	var snapFileExt string = ".png"
	if *out == "" {
		*out = fmt.Sprintf("%s_at_%s%s", id, tsRaw, snapFileExt)
	}

	// Set the video fragment file name from youtube link
	var fragFileExt string = ".mp4"
	fragFileName := fmt.Sprintf("%s_at_%s_fragment%s", id, tsRaw, fragFileExt)

	// use timestamp from video if no flag provided
	if *ts == "" {
		*ts = tsRaw
	}

	// Call capture logic
	if err := capture.GetFrame(*url, *ts, fragFileName, *out); err != nil {
		log.Fatalf("yt-scr: failed to capture frame: %v", err)
	}

	fmt.Printf("Frame captured successfully: %s\n", *out)

}
