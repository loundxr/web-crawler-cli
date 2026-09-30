package crawler

import "github.com/loundxr/web-crawler-cli/internal/model"

func collect(results <-chan taskResult, roots *[]*model.Node) {
	for r := range results {
		if r.node == nil {
			continue
		}
		if r.parent == nil {
			*roots = append(*roots, r.node)
			continue
		}
		r.parent.Links = append(r.parent.Links, r.node)
	}
}
