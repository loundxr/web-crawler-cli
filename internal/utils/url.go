package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func ResolveURL(base *url.URL, href string) (string, bool) {
	href = strings.TrimSpace(href)
	if href == "" {
		return "", false
	}

	parsedURL, err := url.Parse(href)
	if err != nil {
		return "", false
	}

	if parsedURL.Scheme != "" && parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", false
	}

	return base.ResolveReference(parsedURL).String(), true
}

func SplitURLs(urlsStr string) []string {
	parts := strings.Split(urlsStr, ",")
	res := make([]string, 0, len(parts))
	for _, url := range parts {
		url = strings.TrimSpace(url)
		if url != "" {
			res = append(res, url)
		}
	}
	return res
}

func ValidateURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid start URL: %s: %w", raw, err)
	}
	if parsed.Scheme == "" {
		return fmt.Errorf("start URL scheme cannot be empty: %s", raw)
	}
	if parsed.Host == "" {
		return fmt.Errorf("start URL host cannot be empty: %s", raw)
	}
	return nil
}
