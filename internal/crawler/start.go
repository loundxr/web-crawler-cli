package crawler

import (
	"bytes"
	"context"
	"io"

	"github.com/loundxr/web-crawler-cli/internal/fetcher"
	"github.com/loundxr/web-crawler-cli/internal/model"
	"github.com/loundxr/web-crawler-cli/internal/parser"
)

func (c *Crawler) Start(ctx context.Context) []*model.Node {
	results := make(chan *model.Node, len(c.cfg.StartURLs))
	launched := 0

	for _, url := range c.cfg.StartURLs {
		rootHost, err := hostOf(url)
		if err != nil {
			c.logger.Error("failed to extract host", "error", err, "url", url)
			continue
		}

		if !c.tryVisit(url) {
			c.logger.Warn("URL already visited", "url", url)
			continue
		}

		launched++
		go func(url, rootHost string) {
			results <- c.crawlNode(ctx, url, 0, rootHost)
		}(url, rootHost)
	}

	nodes := make([]*model.Node, 0, launched)
	for i := 0; i < launched; i++ {
		if node := <-results; node != nil {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (c *Crawler) crawlNode(ctx context.Context, url string, depth int, rootHost string) *model.Node {
	select {
	case c.sem <- struct{}{}:
	case <-ctx.Done():
		return nil
	}

	reqCtx, cancel := context.WithTimeout(ctx, c.cfg.RequestTimeout)
	outcome := c.fetcher.Fetch(reqCtx, url)

	if outcome.SkipReason != fetcher.NoSkip || outcome.Err != nil {
		cancel()
		<-c.sem
		c.logger.Warn("fetch failed", "reason", outcome.SkipReason, "error", outcome.Err, "url", url)
		return nil
	}

	data, err := io.ReadAll(outcome.Body)
	outcome.Body.Close()

	cancel()
	<-c.sem

	if err != nil {
		c.logger.Warn("read failed", "error", err, "url", url)
		return nil
	}

	page, err := parser.ParseResponseBody(bytes.NewReader(data), url)
	if err != nil {
		c.logger.Warn("parse failed", "error", err, "url", url)
		return nil
	}

	node := model.NewNode(url, page.Title)

	if depth >= c.cfg.MaxDepth || ctx.Err() != nil {
		return node
	}

	candidates := make([]string, 0, len(page.RawLinks))
	for _, link := range page.RawLinks {
		if !sameDomain(link, rootHost) {
			continue
		}
		if !c.tryVisit(link) {
			continue
		}
		candidates = append(candidates, link)
	}

	if len(candidates) == 0 {
		return node
	}

	childResults := make(chan *model.Node, len(candidates))
	for _, link := range candidates {
		go func(link string) {
			childResults <- c.crawlNode(ctx, link, depth+1, rootHost)
		}(link)
	}

	for i := 0; i < len(candidates); i++ {
		if childNode := <-childResults; childNode != nil {
			node.Links = append(node.Links, childNode)
		}
	}

	return node
}
