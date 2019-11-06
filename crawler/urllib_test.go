package crawler

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		expURL *url.URL
	}{
		{
			name:   "normalised path",
			url:    "http://example.com",
			expURL: &url.URL{Scheme: "http", Host: "example.com", Path: "/", Fragment: ""},
		},
		{
			name:   "empty fragment",
			url:    "http://example.com/path#fragment",
			expURL: &url.URL{Scheme: "http", Host: "example.com", Path: "/path", Fragment: ""},
		},
		{
			name: "invalid control character",
			url:  "http://example.com/\x7f",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := parseURL(tt.url)
			if err != nil {
				assert.Nil(t, tt.expURL)
			} else {
				assert.Equal(t, tt.expURL.Scheme, u.Scheme)
				assert.Equal(t, tt.expURL.Host, u.Host)
				assert.Equal(t, tt.expURL.Path, u.Path)
				assert.Equal(t, tt.expURL.Fragment, u.Fragment)
			}
		})
	}
}
