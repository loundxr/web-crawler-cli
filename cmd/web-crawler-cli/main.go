package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/loundxr/web-crawler-cli/internal/config"
	"github.com/loundxr/web-crawler-cli/internal/fetcher"
	"github.com/loundxr/web-crawler-cli/internal/parser"
)

func main() {
	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		log.Fatal("parse config error:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.OverallTimeout)
	defer cancel()

	// fetcher test
	f := fetcher.New(cfg.MaxConcurrency)
	outcome := f.Fetch(ctx, cfg.StartURLs[0])
	data := make([]byte, 0)

	if outcome.Body != nil {
		defer outcome.Body.Close()
		data, err = io.ReadAll(outcome.Body)
		if err != nil {
			log.Fatal("read body error:", err)
		}
		fmt.Printf("body: %s\n", string(data))
	}

	fmt.Printf(
		"code: %d\nskip reason: %s\nerror: %v\n",
		outcome.StatusCode,
		outcome.SkipReason,
		outcome.Err,
	)
	// fetcher test
	//parser test

	page, err := parser.ParseResponseBody(bytes.NewReader(data), cfg.StartURLs[0])
	if err != nil {
		log.Fatal("parse body error:", err)
	}

	fmt.Printf("title: %s\n", page.Title)
	fmt.Printf("raw links: %v\n", page.RawLinks)
	//parser test
}
