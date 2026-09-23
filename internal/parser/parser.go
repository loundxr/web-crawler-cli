package parser

import (
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	Title    string
	RawLinks []string
}

func ParseResponseBody(body io.Reader, baseURL string) (PageData, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return PageData{}, err
	}

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return PageData{}, err
	}

	title := doc.Find("title").Text()
	rawLinks := []string{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok {
			return
		}

		resolved, ok := resolveURL(base, href)
		if !ok {
			return
		}

		rawLinks = append(rawLinks, resolved)
	})

	return PageData{
		Title:    title,
		RawLinks: rawLinks,
	}, nil
}

func resolveURL(base *url.URL, href string) (string, bool) {
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
