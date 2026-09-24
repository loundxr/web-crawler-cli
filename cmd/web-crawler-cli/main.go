package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/loundxr/web-crawler-cli/internal/applog"
	"github.com/loundxr/web-crawler-cli/internal/config"
	"github.com/loundxr/web-crawler-cli/internal/crawler"
	"github.com/loundxr/web-crawler-cli/internal/fetcher"
	"github.com/loundxr/web-crawler-cli/internal/model"
)

func main() {
	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		log.Fatal("parse config error:", err)
	}

	logger, logFile, err := applog.New(cfg.LogPath)
	if err != nil {
		log.Fatal("create logger error:", err)
	}
	defer logFile.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, cfg.OverallTimeout)
	defer cancel()

	f := fetcher.New(cfg.MaxConcurrency)
	c := crawler.New(f, cfg, logger)

	if ctx.Err() != nil {
		log.Println("crawling interrputed or timed out, saving results...")
	}

	nodes := c.Start(ctx)

	if err := writeResultToFile(nodes, cfg.OutputPath); err != nil {
		log.Fatal("write result error:", err)
	}

	log.Println("crawling completed successfully")
}

func writeResultToFile(nodes []*model.Node, path string) error {
	data, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal nodes: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
