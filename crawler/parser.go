package crawler

import (
	"io"

	"golang.org/x/net/html"
)

// Returns anchor "href" value if found.
func getHref(token html.Token) (href string, ok bool) {
	for _, attr := range token.Attr {
		if attr.Key == "href" {
			return attr.Val, true
		}
	}
	return
}

// Resolves the given URL based on the parsed URL from the Crawler context.
// If the URL is relative, resolves reference from the base URL, otherwise checks if scheme and host are the same.
// Returns normalised URL string or empty string if the URL cannot be resolved.
func (c *Crawler) resolveURL(url string) (newURL string, ok bool) {
	u, err := parseURL(url)
	if err != nil {
		return
	}

	if !u.IsAbs() {
		newURL = c.url.ResolveReference(u).String()
		ok = true
	} else if c.url.Scheme == u.Scheme && c.url.Host == u.Host {
		newURL = u.String()
		ok = true
	}
	return
}

// Parses the HTML data, extracts URLs from anchors and sends them to the channel.
func (c *Crawler) parse(data io.Reader) {
	tokenizer := html.NewTokenizer(data)

	for {
		switch tokenizer.Next() {
		case html.StartTagToken:
			token := tokenizer.Token()
			if token.Data != "a" {
				// ignore non-anchor elements
				continue
			}

			url, ok := getHref(token)
			if !ok {
				// ignore anchor elements with no `href` attribute
				continue
			}

			url, ok = c.resolveURL(url)
			if !ok {
				// ignore all external or broken URLs
				continue
			}
			c.chUrl <- &url
		case html.ErrorToken:
			// end of the document
			return
		}
	}
}
