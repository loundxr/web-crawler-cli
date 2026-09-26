package parser_test

import (
	"strings"
	"testing"

	"github.com/loundxr/web-crawler-cli/internal/parser"
)

const sampleHTML = `<html>
<head><title>Page Example</title></head>
<body>
	<a href="/about">About</a>
	<a href="https://other.com">External</a>
	<a href="mailto:test@example.com">Mail To</a>
	<a>No href</a>
</body>
</html>`

func TestParseResponseBody_ExtractsTitleAndLinks(t *testing.T) {
	page, err := parser.ParseResponseBody(strings.NewReader(sampleHTML), "https://example.com/blog")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if page.Title != "Page Example" {
		t.Errorf("title = %q, want %q", page.Title, "Page Example")
	}

	wantLinks := []string{
		"https://example.com/about",
		"https://other.com",
	}

	if len(wantLinks) != len(page.RawLinks) {
		t.Fatalf("Got %d raw links, want %d", len(page.RawLinks), len(wantLinks))
	}
	for i, raw := range page.RawLinks {
		if wantLinks[i] != raw {
			t.Errorf("Got %q link, want %q", raw, wantLinks[i])
		}
	}
}

func TestParseResponseBody_NoTitle(t *testing.T) {
	html := `<html>
				<head></head>
				<body>
					<p>No title</p>
				</body>
			</html>`

	page, err := parser.ParseResponseBody(strings.NewReader(html), "https://example.com")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if page.Title != "" {
		t.Errorf("title = %q, want empty string", page.Title)
	}
}

func TestParseResponseBody_InvalidBaseURL(t *testing.T) {
	html := `<html><body>some body</body></html>`
	_, err := parser.ParseResponseBody(strings.NewReader(html), "https://example.com/%ww")
	if err == nil {
		t.Fatal("want error due to invalid base url")
	}
}
