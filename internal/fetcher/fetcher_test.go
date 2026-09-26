package fetcher_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/loundxr/web-crawler-cli/internal/fetcher"
)

func TestFetch_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html><title>success</title></html>"))
	}))
	defer server.Close()

	f := fetcher.New(10)
	outcome := f.Fetch(context.Background(), server.URL)

	if outcome.Body != nil {
		defer outcome.Body.Close()
	}
	if outcome.SkipReason != fetcher.NoSkip {
		t.Fatalf("skip reason = %q, want NoSkip", outcome.SkipReason)
	}
	if outcome.Err != nil {
		t.Fatalf("failed to fetch: %v", outcome.Err)
	}
	if outcome.StatusCode != http.StatusOK {
		t.Errorf("code = %v, want %v", outcome.StatusCode, http.StatusOK)
	}
	if outcome.Body == nil {
		t.Fatal("expected not empty body")
	}

	data, err := io.ReadAll(outcome.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if !strings.Contains(string(data), "<title>success</title>") {
		t.Errorf("unexpected response body: %v", string(data))
	}
}

func TestFetch_Redirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/test", http.StatusFound)
	}))
	defer server.Close()

	f := fetcher.New(10)
	outcome := f.Fetch(context.Background(), server.URL)

	if outcome.SkipReason != fetcher.SkipRedirect {
		t.Fatalf("got SkipReason = %q, want SkipRedirect", outcome.SkipReason)
	}
	if outcome.StatusCode != http.StatusFound {
		t.Errorf("got StatusCode = %v, want %v", outcome.StatusCode, http.StatusFound)
	}
	if outcome.Body != nil {
		t.Error("body expected to be nil after redirect")
	}

}

func TestFetch_BadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnavailableForLegalReasons)
	}))
	defer server.Close()

	f := fetcher.New(10)
	outcome := f.Fetch(context.Background(), server.URL)

	if outcome.StatusCode != http.StatusUnavailableForLegalReasons {
		t.Errorf("got StatusCode = %v, want %v", outcome.StatusCode, http.StatusUnavailableForLegalReasons)
	}
	if outcome.SkipReason != fetcher.SkipBadStatus {
		t.Fatalf("got SkipReason = %q, want SkipBadStatus", outcome.SkipReason)
	}
}

func TestFetch_NonHTML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"test":true}`))
	}))

	f := fetcher.New(10)
	outcome := f.Fetch(context.Background(), server.URL)

	if outcome.SkipReason != fetcher.SkipNonHTML {
		t.Fatalf("got SkipReason = %q, want SkipNonHTML", outcome.SkipReason)
	}
}

func TestFetch_ConnectionError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	closed := server.URL
	server.Close()

	f := fetcher.New(10)
	outcome := f.Fetch(context.Background(), closed)

	if outcome.SkipReason != fetcher.SkipError {
		t.Fatalf("got SkipReason = %q, want SkipError", outcome.SkipReason)
	}
	if outcome.Err == nil {
		t.Error("err == nil, expected connection error")
	}
}

func TestFetch_ConnectionTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	f := fetcher.New(10)
	outcome := f.Fetch(ctx, server.URL)

	if outcome.SkipReason != fetcher.SkipError {
		t.Fatalf("got SkipReason = %q, want SkipError", outcome.SkipReason)
	}
	if outcome.Err == nil {
		t.Error("err == nil, expected connection timeout err")
	}
}
