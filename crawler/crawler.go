package crawler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	urllib "net/url"
	"sync"
)

var (
	ErrInvalidStreamWriter = errors.New("invalid stream writer")
	ErrEmptyHost           = errors.New("empty host in the given URL")
	ErrMaxPageLimit        = errors.New("reached maximum pages limit")
)

// Crawler context.
type Crawler struct {
	url      *urllib.URL         // parsed base URL
	chUrl    chan *string        // channel for extracted URLs
	chErr    chan error          // channel for routine errors
	siteMap  map[string]struct{} // site map index
	maxPages int                 // maximum pages limit
	stream   io.Writer           // stream writer to output extracted URLs
	logger   *log.Logger         // error logger
	state    int                 // number of running routines
	wg       sync.WaitGroup      // main waiting group
}

// Makes request to a given URL and parses the response body (used as go routine).
func (c *Crawler) process(url string) {
	defer func() {
		// send all recovered panics to the error channel
		if err := recover(); err != nil {
			c.chErr <- err.(error)
		}

		// this will notify the caller that the current routine has finished
		c.chUrl <- nil
		c.wg.Done()
	}()

	// TODO: Implement retry if the server is not responding.
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body := resp.Body
	defer body.Close()

	c.parse(body)
}

// Pushes the given URL to the site map and starts the process routine.
// If the max pages limit is hit, gracefully returns error.
func (c *Crawler) push(url string) error {
	if len(c.siteMap) == c.maxPages {
		return ErrMaxPageLimit
	}

	fmt.Fprintln(c.stream, url)
	c.siteMap[url] = struct{}{}

	c.wg.Add(1)
	go c.process(url)
	c.state++

	return nil
}

// Concurrently creates the site map by the given URL and outputs extracted URLs to the stream.
func (c *Crawler) Stream(stream io.Writer, url string) (err error) {
	if c.stream = stream; c.stream == nil {
		return ErrInvalidStreamWriter
	}

	c.url, err = parseURL(url)
	if err != nil {
		return
	} else if c.url.Host == "" {
		return ErrEmptyHost
	}

	// assemble normalised URL from parsed
	// (see `parseURL` definition)
	url = c.url.String()

	// clear site map to reuse
	// (this code gets optimised by the compiler for Go 1.11+)
	for k := range c.siteMap {
		delete(c.siteMap, k)
	}

	c.chUrl = make(chan *string)
	c.chErr = make(chan error)

	defer func() {
		c.wg.Wait()
		close(c.chUrl)
		close(c.chErr)
	}()

	// push base URL for processing
	if err = c.push(url); err != nil {
		return
	}

	for {
		select {
		case url := <-c.chUrl:
			if url != nil {
				_, found := c.siteMap[*url]
				if !found {
					// push extracted URL for processing
					// or skip if hit the limit
					if err = c.push(*url); err != nil {
						continue
					}
				}
			} else {
				c.state--
				if c.state == 0 {
					// all routines are completed
					// (i.e. all pages were processed)
					return
				}
			}
		case err := <-c.chErr:
			if c.logger != nil {
				c.logger.Println(err)
			}
		}
	}
}

// Returns the size of the site map.
func (c *Crawler) Count() int {
	return len(c.siteMap)
}

// Creates and returns new instance of Crawler.
func NewCrawler(maxPages uint, logger *log.Logger) *Crawler {
	return &Crawler{
		siteMap:  make(map[string]struct{}, maxPages),
		maxPages: int(maxPages),
		logger:   logger,
	}
}
