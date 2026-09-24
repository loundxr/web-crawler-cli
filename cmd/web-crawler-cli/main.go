package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/loundxr/web-crawler-cli/internal/applog"
	"github.com/loundxr/web-crawler-cli/internal/config"
	"github.com/loundxr/web-crawler-cli/internal/crawler"
	"github.com/loundxr/web-crawler-cli/internal/fetcher"
	"github.com/loundxr/web-crawler-cli/internal/utils"
)

func main() {
	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		log.Fatal("parse config error:", err)
	}

	logger, cleanup, err := applog.New(cfg.LogPath)
	if err != nil {
		log.Fatal("create logger error:", err)
	}
	defer cleanup()

	f := fetcher.New(cfg.MaxConcurrency)
	c := crawler.New(f, cfg, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, cfg.OverallTimeout)
	defer cancel()

	start := time.Now()
	nodes := c.Start(ctx)
	fmt.Println("crawl time: ", time.Since(start))

	if err := utils.WriteResults(nodes, cfg.OutputPath); err != nil {
		log.Fatal("write result error:", err)
	}

	log.Println("crawling completed successfully")
}
