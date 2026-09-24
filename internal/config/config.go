package config

import (
	"flag"
	"fmt"
	"net/url"
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
	depth := fs.Int("depth", 2, "maximum depth of recursive traversal")
	overallTimeout := fs.Duration("timeout", 1*time.Minute, "overall timeout for program execution")
	requestTimeout := fs.Duration("request-timeout", 10*time.Second, "timeout per request")
	outputPath := fs.String("output", "out/result.json", "path to the result json file")
	logPath := fs.String("log", "out/crawler.log", "path to the log file")

	if err := fs.Parse(args); err != nil {
		return Config{}, fmt.Errorf("parse flags: %w", err)
	}

	startURLs := splitURLs(*urls)
	if len(startURLs) == 0 {
		return Config{}, fmt.Errorf("--urls is required")
	}

	for _, u := range startURLs {
		if err := validateURL(u); err != nil {
			return Config{}, err
		}
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

func validateURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid start URL: %s: %w", raw, err)
	}
	if parsed.Scheme == "" {
		return fmt.Errorf("start URL scheme cannot be empty: %s", raw)
	}
	if parsed.Host == "" {
		return fmt.Errorf("start URL host cannot be empty: %s", raw)
	}
	return nil
}
