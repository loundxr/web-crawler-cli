package utils_test

import (
	"net/url"
	"testing"

	"github.com/loundxr/web-crawler-cli/internal/utils"
)

func TestResolveURL(t *testing.T) {
	base, err := url.Parse("https://example.com/something/index.html")
	if err != nil {
		t.Fatalf("failed to parse base URL for testing: %v", err)
	}

	tests := []struct {
		name    string
		href    string
		wantURL string
		wantOK  bool
	}{
		{
			name:    "absolute path",
			href:    "/about",
			wantURL: "https://example.com/about",
			wantOK:  true,
		},
		{
			name:    "relative path without leading slash",
			href:    "about.html",
			wantURL: "https://example.com/something/about.html",
			wantOK:  true,
		},
		{
			name:    "one level up",
			href:    "../blog",
			wantURL: "https://example.com/blog",
			wantOK:  true,
		},
		{
			name:    "protocol-relative url",
			href:    "//other.by/partners",
			wantURL: "https://other.by/partners",
			wantOK:  true,
		},
		{
			name:    "another domain absolute url",
			href:    "https://random.com/policy",
			wantURL: "https://random.com/policy",
			wantOK:  true,
		},
		{
			name:   "mailto - non-HTML",
			href:   "mailto:test@example.com",
			wantOK: false,
		},
		{
			name:   "tel - non-HTML",
			href:   "tel:+88005553535",
			wantOK: false,
		},
		{
			name:   "empty href",
			href:   "",
			wantOK: false,
		},
		{
			name:   "spaces only",
			href:   "    ",
			wantOK: false,
		},
		{
			name:    "anchor link",
			href:    "#contacts",
			wantURL: "https://example.com/something/index.html#contacts",
			wantOK:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := utils.ResolveURL(base, tc.href)
			if ok != tc.wantOK {
				t.Fatalf("utils.ResolveURL(%q) ok = %v, want %v", tc.href, ok, tc.wantOK)
			}

			if ok && got != tc.wantURL {
				t.Errorf("utils.ResolveURL(%q) = %q, want %q", tc.href, got, tc.wantURL)
			}
		})
	}
}
