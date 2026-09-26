package crawler

import (
	"testing"
)

func TestSameDomain(t *testing.T) {
	tests := []struct {
		name     string
		link     string
		rootHost string
		want     bool
	}{
		{
			name:     "same host",
			link:     "https://example.com/about",
			rootHost: "example.com",
			want:     true,
		},
		{
			name:     "other host",
			link:     "https://other.com/another",
			rootHost: "example.com",
			want:     false,
		},
		{
			name:     "ignore case",
			link:     "https://EXAMPLE.com/news",
			rootHost: "example.com",
			want:     true,
		},
		{
			name:     "ignore port",
			link:     "https://example.com:6767/members",
			rootHost: "example.com",
			want:     true,
		},
		{
			name:     "subdomain",
			link:     "https://cool.example.com/map",
			rootHost: "example.com",
			want:     false,
		},
		{
			name:     "broken URL",
			link:     "https://example.com/%ww",
			rootHost: "example.com",
			want:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := sameDomain(tc.link, tc.rootHost)
			if got != tc.want {
				t.Errorf("sameDomain(%q, %q) = %v, want %v", tc.link, tc.rootHost, got, tc.want)
			}
		})
	}
}

func TestHostOf(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		want    string
		wantErr bool
	}{
		{
			name: "default URL",
			link: "https://example.com/results",
			want: "example.com",
		},
		{
			name: "with port",
			link: "https://example.com:8888",
			want: "example.com",
		},
		{
			name:    "broken URL",
			link:    "https://example.com/%kk",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := hostOf(tc.link)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("hostOf(%q) didn't return an error when expected", tc.link)
				}
				return
			}

			if err != nil {
				t.Fatalf("hostOf(%q) unexpected error: %v", tc.link, err)
			}
			if got != tc.want {
				t.Errorf("hostOf(%q) = %q, want %q", tc.link, got, tc.want)
			}
		})
	}
}
