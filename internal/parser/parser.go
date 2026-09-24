package parser

import (
	"io"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/loundxr/web-crawler-cli/internal/utils"
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

		resolved, ok := utils.ResolveURL(base, href)
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
