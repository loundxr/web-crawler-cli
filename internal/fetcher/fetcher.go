package fetcher

import (
	"context"
	"net/http"
	"strings"
)

type Fetcher struct {
	client *http.Client
}

func New(maxIdleConnsPerHost int) *Fetcher {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = maxIdleConnsPerHost
	return &Fetcher{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: transport,
		},
	}
}

func (f *Fetcher) Fetch(ctx context.Context, url string) FetchOutcome {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return FetchOutcome{SkipReason: SkipError, Err: err}
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return FetchOutcome{SkipReason: SkipError, Err: err}
	}

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		resp.Body.Close()
		return FetchOutcome{SkipReason: SkipRedirect, StatusCode: resp.StatusCode}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		resp.Body.Close()
		return FetchOutcome{SkipReason: SkipBadStatus, StatusCode: resp.StatusCode}
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		resp.Body.Close()
		return FetchOutcome{SkipReason: SkipNonHTML, StatusCode: resp.StatusCode}
	}

	return FetchOutcome{
		Body:       resp.Body,
		StatusCode: resp.StatusCode,
		SkipReason: NoSkip,
	}
}
