package crawler

import (
	"context"
	"sync"

	"github.com/loundxr/web-crawler-cli/internal/model"
)

type task struct {
	url      string
	depth    int
	rootHost string
	parent   *model.Node
}

type taskResult struct {
	node   *model.Node
	parent *model.Node
}

func (c *Crawler) Start(ctx context.Context) []*model.Node {
	submit := make(chan task)
	jobs := make(chan task)
	results := make(chan taskResult)

	var tasksWG sync.WaitGroup
	stopDispatch := make(chan struct{})

	go c.dispatch(submit, jobs, stopDispatch)

	var workersWG sync.WaitGroup
	for i := 0; i < c.cfg.MaxConcurrency; i++ {
		workersWG.Add(1)
		go func() {
			defer workersWG.Done()
			c.worker(ctx, jobs, submit, results, &tasksWG)
		}()
	}

	roots := make([]*model.Node, 0)
	collectDone := make(chan struct{})
	go func() {
		collect(results, &roots)
		close(collectDone)
	}()

	for _, u := range c.cfg.StartURLs {
		rootHost, err := hostOf(u)
		if err != nil {
			c.logger.Error(
				"failed to extract host",
				"error", err,
				"url", u,
			)
			continue
		}

		if !c.tryVisit(u) {
			c.logger.Error(
				"URL already visited",
				"url", u,
			)
			continue
		}

		tasksWG.Add(1)
		submit <- task{url: u, depth: 0, rootHost: rootHost, parent: nil}
	}

	tasksWG.Wait()
	close(stopDispatch)

	workersWG.Wait()

	close(results)
	<-collectDone

	return roots
}
