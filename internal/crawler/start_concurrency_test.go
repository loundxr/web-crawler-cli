package crawler_test

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/loundxr/web-crawler-cli/internal/config"
	"github.com/loundxr/web-crawler-cli/internal/crawler"
	"github.com/loundxr/web-crawler-cli/internal/fetcher"
)

type concurrencyTracker struct {
	mu      sync.Mutex
	current int
	peak    int
}

func (t *concurrencyTracker) increment() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.current++
	if t.current > t.peak {
		t.peak = t.current
	}
}

func (t *concurrencyTracker) decrement() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.current--
}

func (t *concurrencyTracker) Peak() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.peak
}

func trackConcurrency(t *concurrencyTracker) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.increment()
			defer t.decrement()
			h.ServeHTTP(w, r)
		})
	}
}

func TestCrawler_Start_FanOutConcurrencyLimit(t *testing.T) {
	const (
		fanOut         = 100
		maxConcurrency = 10
	)

	pages := make(map[string]testPage, fanOut+1)
	links := make([]string, 0, fanOut)

	for i := 0; i < fanOut; i++ {
		path := fmt.Sprintf("/page%d", i)
		links = append(links, path)
		pages[path] = testPage{title: fmt.Sprintf("Page %d", i), delay: 30 * time.Millisecond}
	}
	pages["/"] = testPage{title: "root", links: links}

	tracker := &concurrencyTracker{}
	server := newTestServer(pages, trackConcurrency(tracker))
	defer server.Close()

	f := fetcher.New(10)
	c := crawler.New(
		f,
		config.Config{
			MaxDepth:       1,
			RequestTimeout: 5 * time.Second,
			MaxConcurrency: maxConcurrency,
			StartURLs:      []string{server.URL},
		},
		newTestLogger(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	roots := c.Start(ctx)

	if len(roots) != 1 {
		t.Fatalf("len(roots)=%d, want 1", len(roots))
	}

	if len(roots[0].Links) != fanOut {
		t.Fatalf("len(roots[0].Links=%d, want %d", len(roots[0].Links), fanOut)
	}

	peak := tracker.Peak()
	if peak > maxConcurrency {
		t.Errorf("peak of parallel requests = %d, shouldn't be more than maxConcurrency = %d", peak, maxConcurrency)
	}

	if peak < maxConcurrency {
		t.Logf("peak of parallel requests = %d, expected %d", peak, maxConcurrency)
	}
}
