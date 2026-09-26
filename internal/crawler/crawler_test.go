package crawler

import (
	"sync"
	"testing"
)

func TestCrawler_tryVisitDefault(t *testing.T) {
	c := &Crawler{visited: make(map[string]struct{})}

	if !c.tryVisit("https://example.com") {
		t.Fatalf("tryVisit for first appearance of an URL must return true")
	}
	if c.tryVisit("https://example.com") {
		t.Fatalf("tryVisit for the same URL must return false")
	}
	if !c.tryVisit("https://other.com") {
		t.Fatalf("tryVisit for a new URL must return true")
	}
}

func TestCrawler_tryVisitConcurrent(t *testing.T) {
	c := &Crawler{visited: make(map[string]struct{})}

	urls := []string{"https://example.com/about", "https://example.com/blog", "https://example.com/categories"}
	const NumOfGoroutines = 100

	var wg sync.WaitGroup
	results := make(chan bool, NumOfGoroutines*len(urls))

	for i := 0; i < NumOfGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, u := range urls {
				results <- c.tryVisit(u)
			}
		}()
	}

	wg.Wait()
	close(results)

	trueCount := 0
	for got := range results {
		if got {
			trueCount++
		}
	}

	wantTrue := len(urls)
	if trueCount != wantTrue {
		t.Errorf("tryVisit true results = %d, want %d", trueCount, wantTrue)
	}
	if len(c.visited) != len(urls) {
		t.Errorf("len(visited) = %d, want %d", len(c.visited), len(urls))
	}
}
