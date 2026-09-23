package config

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

// number of goroutines used for fetching pages concurrently
const MaxConcurrency = 10

type Config struct {
	StartURLs      []string
	MaxDepth       int
	OverallTimeout time.Duration
	RequestTimeout time.Duration
	MaxConcurrency int
	OutputPath     string
	LogPath        string
}

func ParseFlags(args []string) (Config, error) {
	fs := flag.NewFlagSet("crawler-cli", flag.ContinueOnError)

	urls := fs.String("urls", "", "comma-separated URLs (required)")
	depth := fs.Int("depth", 0, "maximum depth of recursive traversal")
	overallTimeout := fs.Duration("timeout", 1*time.Minute, "overall timeout for program execution")
	requestTimeout := fs.Duration("request-timeout", 10*time.Second, "timeout per request")
	outputPath := fs.String("output", "../../out/result.json", "path to the result json file")
	logPath := fs.String("log", "../../out/crawler.log", "path to the log file")

	if err := fs.Parse(args); err != nil {
		return Config{}, fmt.Errorf("parse flags: %w", err)
	}

	startURLs := splitURLs(*urls)
	if len(startURLs) == 0 {
		return Config{}, fmt.Errorf("--urls is required")
	}

	if *depth < 0 {
		return Config{}, fmt.Errorf("depth must be non-negative")
	}

	return Config{
		StartURLs:      startURLs,
		MaxDepth:       *depth,
		OverallTimeout: *overallTimeout,
		RequestTimeout: *requestTimeout,
		MaxConcurrency: MaxConcurrency,
		OutputPath:     strings.TrimSpace(*outputPath),
		LogPath:        strings.TrimSpace(*logPath),
	}, nil
}

func splitURLs(urlsStr string) []string {
	parts := strings.Split(urlsStr, ",")
	res := make([]string, 0, len(parts))
	for _, url := range parts {
		url = strings.TrimSpace(url)
		if url != "" {
			res = append(res, url)
		}
	}
	return res
}
