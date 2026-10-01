package crawler_test

import (
	"context"
	"testing"
	"time"

	"github.com/loundxr/web-crawler-cli/internal/config"
	"github.com/loundxr/web-crawler-cli/internal/crawler"
	"github.com/loundxr/web-crawler-cli/internal/fetcher"
	"github.com/loundxr/web-crawler-cli/internal/model"
)

func TestCrawler_Start_Cancellation(t *testing.T) {
	server := newTestServer(map[string]testPage{
		"/":  {title: "root", links: []string{"/a", "/b", "/c"}, delay: 200 * time.Millisecond},
		"/a": {title: "A", delay: 200 * time.Millisecond},
		"/b": {title: "B", delay: 200 * time.Millisecond},
		"/c": {title: "C", delay: 200 * time.Millisecond},
	})
	defer server.Close()

	f := fetcher.New(10)
	c := crawler.New(
		f,
		config.Config{
			MaxDepth:       3,
			OverallTimeout: 5 * time.Second,
			MaxConcurrency: 10,
		},
		newTestLogger(),
	)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan []*model.Node)
	go func() {
		done <- c.Start(ctx)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case <-done:
		if ctx.Err() == nil {
			t.Fatal("ctx.Err() == nil after cancellation")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("c.Start() didn't finish after 3sec after call, possible deadlock")
	}
}

func TestCrawler_Start_OverallTimeout(t *testing.T) {
	server := newTestServer(map[string]testPage{
		"/":  {title: "root", links: []string{"/a"}, delay: 300 * time.Millisecond},
		"/a": {title: "A", delay: 300 * time.Millisecond},
	})
	defer server.Close()

	f := fetcher.New(10)
	c := crawler.New(
		f,
		config.Config{
			MaxDepth:       3,
			OverallTimeout: time.Minute,
			RequestTimeout: 5 * time.Second,
			MaxConcurrency: 10,
			StartURLs:      []string{server.URL},
		},
		newTestLogger(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	initTime := time.Now()
	done := make(chan []*model.Node)
	go func() {
		done <- c.Start(ctx)
	}()

	select {
	case <-done:
		workTime := time.Since(initTime)
		if workTime > 700*time.Millisecond {
			t.Fatalf("c.Start() finished after %v, overallTimeout didn't stopped the app flow", workTime)
		}
		if ctx.Err() == nil {
			t.Fatal("ctx.Err() == nil, want context.DeadlineExceeded")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("c.Start() didn't finish after 3sec after call, possible deadlock")
	}
}
