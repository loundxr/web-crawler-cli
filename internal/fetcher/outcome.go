package fetcher

import "io"

type SkipReason string

const (
	NoSkip        SkipReason = ""
	SkipRedirect  SkipReason = "redirect"
	SkipNonHTML   SkipReason = "non HTML"
	SkipBadStatus SkipReason = "bad status"
	SkipError     SkipReason = "error"
)

type FetchOutcome struct {
	Body       io.ReadCloser
	StatusCode int
	SkipReason SkipReason
	Err        error
}
