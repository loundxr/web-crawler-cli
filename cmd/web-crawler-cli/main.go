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

// TODO: logs formatting
// TODO: README
// TODO: unit-tests

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

	start := time.Now()
	nodes := c.Start(ctx)
	fmt.Println("crawl time: ", time.Since(start))

	if err := utils.WriteResults(nodes, cfg.OutputPath); err != nil {
		log.Fatal("write result error:", err)
	}

	log.Println("crawling completed successfully")
}
