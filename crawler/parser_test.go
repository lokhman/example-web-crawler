package crawler

import (
	urllib "net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawler_parse(t *testing.T) {
	tests := []struct {
		name string
		data string
		urls []string
	}{
		{
			name: "empty document",
			data: "",
			urls: []string{},
		},
		{
			name: "no anchor elements",
			data: "<html><body><h1>Test</h1><p>Content</p></body></html>",
			urls: []string{},
		},
		{
			name: "anchor with no href attribute",
			data: "<html><a>Anchor</a></html>",
			urls: []string{},
		},
		{
			name: "anchor with many href attributes",
			data: `<html><a href="link1.html" href="link2.html">Anchor</a></html>`,
			urls: []string{"http://example.com/link1.html"},
		},
		{
			name: "anchor with invalid href attribute",
			data: "<html><a href=\"/link\x02.html\">Anchor</a></html>",
			urls: []string{},
		},
		{
			name: "many anchor elements with different href attributes",
			data: `<html>
				<a href="/">Anchor 1</a>
				<a href="link2.html">Anchor 2</a>
				<a href="/link3.html">Anchor 3</a>
				<a href="test/link4.html">Anchor 4</a>
				<a href="http://example.com/link5.html">Anchor 5</a>
				<a href="http://sub.example.com/link6.html">Anchor 6</a>
				<a href="http://external.com/link7.html">Anchor 7</a>
			</html>`,
			urls: []string{
				"http://example.com/",
				"http://example.com/link2.html",
				"http://example.com/link3.html",
				"http://example.com/test/link4.html",
				"http://example.com/link5.html",
			},
		},
	}

	crawler := &Crawler{
		url: &urllib.URL{
			Scheme: "http",
			Host:   "example.com",
			Path:   "/index.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.data)
			crawler.chUrl = make(chan *string)

			go func() {
				defer close(crawler.chUrl)
				crawler.parse(reader)
			}()

			urls := make([]string, 0)
			for url := range crawler.chUrl {
				urls = append(urls, *url)
			}

			assert.Equal(t, tt.urls, urls)
		})
	}
}
