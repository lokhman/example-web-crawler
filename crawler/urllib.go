package crawler

import (
	urllib "net/url"
)

// Parses the given URL and normalises its components.
func parseURL(url string) (u *urllib.URL, err error) {
	u, err = urllib.Parse(url)
	if err != nil {
		return
	}
	if u.Path == "" {
		u.Path = "/"
	}
	u.Fragment = ""
	return
}
