package crawler

import (
	"net/url"
	"strings"
)

func sameDomain(link, rootHost string) bool {
	linkURL, err := url.Parse(link)
	if err != nil {
		return false
	}
	return strings.EqualFold(linkURL.Hostname(), rootHost)
}

func hostOf(link string) (string, error) {
	linkURL, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	return linkURL.Hostname(), nil
}
