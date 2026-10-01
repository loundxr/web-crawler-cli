package crawler

import (
	"bytes"
	"context"
	"io"
	"sync"

	"github.com/loundxr/web-crawler-cli/internal/model"
	"github.com/loundxr/web-crawler-cli/internal/parser"
)

func (c *Crawler) worker(
	ctx context.Context,
	jobs <-chan task,
	submit chan<- task,
	results chan<- taskResult,
	wg *sync.WaitGroup,
) {
	for t := range jobs {
		node := c.processTask(ctx, t, submit, wg)
		results <- taskResult{parent: t.parent, node: node}
		wg.Done()
	}
}

func (c *Crawler) processTask(ctx context.Context, t task, submit chan<- task, wg *sync.WaitGroup) *model.Node {
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	reqCtx, cancel := context.WithTimeout(ctx, c.cfg.RequestTimeout)

	outcome := c.fetcher.Fetch(reqCtx, t.url)

	if outcome.Body == nil {
		cancel()
		c.logger.Warn(
			"fetch failed",
			"status", outcome.StatusCode,
			"reason", outcome.SkipReason,
			"error", outcome.Err,
			"url", t.url,
		)
		return nil
	}

	data, err := io.ReadAll(outcome.Body)
	outcome.Body.Close()
	cancel()

	if err != nil {
		c.logger.Warn(
			"read failed",
			"error", err,
			"url", t.url,
		)
		return nil
	}

	page, err := parser.ParseResponseBody(bytes.NewReader(data), t.url)
	if err != nil {
		c.logger.Warn(
			"parse failed",
			"error", err,
			"url", t.url,
		)
		return nil
	}

	node := model.NewNode(t.url, page.Title)

	if t.depth >= c.cfg.MaxDepth || ctx.Err() != nil {
		return node
	}

	for _, link := range page.RawLinks {
		if ctx.Err() != nil {
			break
		}

		if !sameDomain(link, t.rootHost) {
			continue
		}

		if !c.tryVisit(link) {
			continue
		}

		wg.Add(1)
		select {
		case submit <- task{
			url:      link,
			depth:    t.depth + 1,
			rootHost: t.rootHost,
			parent:   node,
		}:
		case <-ctx.Done():
			wg.Done()
		}
	}

	return node
}
