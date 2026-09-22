package config

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

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

	urls := fs.String("urls", "", "стартовые URL'ы через запятую (обязательный параметр)")
	depth := fs.Int("depth", 0, "максимальная глубина рекурсивного обхода")
	overallTimeout := fs.Duration("timeout", 1*time.Minute, "общий таймаут выполнения")
	requestTimeout := fs.Duration("request-timeout", 5*time.Second, "таймаут выполнения одного запроса")
	outputPath := fs.String("output", "../../out/result.json", "путь к файлу с результатом (JSON)")
	logPath := fs.String("log", "../../out/crawler.log", "путь к лог-файлу")

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
		MaxConcurrency: 10,
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
