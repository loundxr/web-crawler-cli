package crawler

import (
	"log/slog"
	"sync"

	"github.com/loundxr/web-crawler-cli/internal/config"
	"github.com/loundxr/web-crawler-cli/internal/fetcher"
)

type Crawler struct {
	fetcher *fetcher.Fetcher
	cfg     config.Config
	logger  *slog.Logger
	sem     chan struct{}
	mu      sync.Mutex
	visited map[string]struct{}
}

func New(f *fetcher.Fetcher, cfg config.Config, l *slog.Logger) *Crawler {
	return &Crawler{
		fetcher: f,
		cfg:     cfg,
		logger:  l,
		sem:     make(chan struct{}, cfg.MaxConcurrency),
		visited: make(map[string]struct{}),
	}
}

func (c *Crawler) tryVisit(url string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.visited[url]; ok {
		return false
	}
	c.visited[url] = struct{}{}
	return true
}
