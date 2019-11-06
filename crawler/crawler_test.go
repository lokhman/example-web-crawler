package crawler

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCrawler(t *testing.T) {
	const maxPages uint = 23
	var logger = log.New(os.Stderr, "", log.LstdFlags)

	crawler := NewCrawler(maxPages, logger)
	assert.Equal(t, maxPages, uint(crawler.maxPages))
	assert.Equal(t, logger, crawler.logger)
	assert.NotNil(t, crawler.siteMap)
}

func TestCrawler_Stream(t *testing.T) {
	loggerOut := new(strings.Builder)
	logger := log.New(loggerOut, "", 0)

	crawler := NewCrawler(1000, logger)
	stream := new(strings.Builder)

	t.Run("unknown server", func(t *testing.T) {
		_ = crawler.Stream(stream, "http://unknown")
		assert.Equal(t, "Get http://unknown/: dial tcp: lookup unknown: no such host\n", loggerOut.String())
	})

	t.Run("invalid URL", func(t *testing.T) {
		err := crawler.Stream(stream, "http://unknown/\x03")
		assert.NotNil(t, err)
		assert.Equal(t, "parse http://unknown/\x03: net/url: invalid control character in URL", err.Error())
	})

	t.Run("empty host", func(t *testing.T) {
		for _, url := range []string{"http://", "nothing", ""} {
			err := crawler.Stream(stream, url)
			assert.Equal(t, ErrEmptyHost, err)
		}
	})

	pages := map[string]string{
		"page1.html": "Page 1",
		"page2.html": "Page 2",
		"page3.html": "Page 3",
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "<html>")
		for url, text := range pages {
			fmt.Fprintf(w, `<a href="%s">%s</a>\n`, url, text)
		}
		fmt.Fprintln(w, "</html>")
	})
	server := httptest.NewServer(handler)

	t.Run("invalid stream writer", func(t *testing.T) {
		err := crawler.Stream(nil, server.URL)
		assert.Equal(t, ErrInvalidStreamWriter, err)
	})

	t.Run("max limit reached", func(t *testing.T) {
		for _, maxPages := range []uint{0, 1, 2} {
			crawler := NewCrawler(maxPages, logger)
			err := crawler.Stream(stream, server.URL)
			assert.Equal(t, ErrMaxPageLimit, err)
		}
	})

	t.Run("success", func(t *testing.T) {
		err := crawler.Stream(stream, server.URL)
		assert.Nil(t, err)

		expSiteMap := map[string]struct{}{
			server.URL + "/": {},
		}

		for url := range pages {
			expSiteMap[server.URL+"/"+url] = struct{}{}
		}
		assert.Equal(t, expSiteMap, crawler.siteMap)
		assert.Equal(t, len(expSiteMap), crawler.Count())
	})
}
