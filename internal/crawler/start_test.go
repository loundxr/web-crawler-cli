package crawler_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/loundxr/web-crawler-cli/internal/config"
	"github.com/loundxr/web-crawler-cli/internal/crawler"
	"github.com/loundxr/web-crawler-cli/internal/fetcher"
	"github.com/loundxr/web-crawler-cli/internal/model"
)

type testPage struct {
	title string
	links []string
}

func newTestServer(pages map[string]testPage) *httptest.Server {
	mux := http.NewServeMux()

	for path, page := range pages {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			var linksHTML strings.Builder
			for _, l := range page.links {
				fmt.Fprintf(&linksHTML, `<a href="%s">link</a>`, l)
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, `<html><head><title>%s</title></head><body>%s</body></html>`, page.title, linksHTML.String())
		})
	}

	return httptest.NewServer(mux)
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func mapByResource(nodes []*model.Node) map[string]*model.Node {
	m := make(map[string]*model.Node, len(nodes))
	for _, n := range nodes {
		m[n.Resource] = n
	}
	return m
}

func TestCrawler_Start_MultiLevelTraversal(t *testing.T) {
	server := newTestServer(map[string]testPage{
		"/":   {title: "root", links: []string{"/a", "/b"}},
		"/a":  {title: "A", links: []string{"/a1"}},
		"/b":  {title: "B"},
		"/a1": {title: "A1"},
	})
	defer server.Close()

	f := fetcher.New(10)
	c := crawler.New(
		f,
		config.Config{
			MaxDepth:       2,
			MaxConcurrency: 10,
			OverallTimeout: 1 * time.Minute,
			RequestTimeout: 5 * time.Second,
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

	root := roots[0]
	if root.Title != "root" {
		t.Errorf("root title=%q, want %q", root.Title, "root")
	}

	if len(root.Links) != 2 {
		t.Fatalf("root links len=%d, want 2", len(root.Links))
	}

	children := mapByResource(root.Links)

	a, ok := children[server.URL+"/a"]
	if !ok {
		t.Fatalf("child node %s/a not found among %+v", server.URL, root.Links)
	}
	if a.Title != "A" {
		t.Errorf("/a title=%q, want %q", a.Title, "A")
	}
	if len(a.Links) != 1 {
		t.Fatalf("/a links len=%d, want 1", len(a.Links))
	}

	a1 := a.Links[0]
	// why localhost:8888/a1, not localhost:8888/a/a1???????
	if a1.Resource != server.URL+"/a1" || a1.Title != "A1" {
		t.Errorf("/a1=%+v, want resource=%s/a1, title=%q", a1, server.URL, "A1")
	}
	if len(a1.Links) != 0 {
		t.Fatalf("/a1 on max depth must be a leaf, not a node, got len(a1.Links)=%d", len(a1.Links))
	}

	b, ok := children[server.URL+"/b"]
	if !ok {
		t.Fatalf("child node %s/b not found among %+v", server.URL, root.Links)
	}
	if b.Title != "B" {
		t.Errorf("/b title=%q, want %q", b.Title, "B")
	}
	if len(b.Links) != 0 {
		t.Fatalf("/b links len=%d, want 0", len(b.Links))
	}
}
