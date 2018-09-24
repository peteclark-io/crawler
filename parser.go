package main

import (
	"errors"
	"io"
	"net/url"

	"golang.org/x/net/html"
)

var errExternalURL = errors.New("provided href points to an external domain")

func parser(r io.Reader) ([]string, error) {
	z := html.NewTokenizer(r)

	hrefs := make([]string, 0)
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			finished := z.Err() == io.EOF
			if finished {
				return hrefs, nil
			} else if err := z.Err(); err != nil {
				return hrefs, z.Err()
			}

		case html.StartTagToken:
			tk := z.Token()

			if tk.Data != "a" {
				continue
			}

			for _, v := range tk.Attr {
				if v.Key == "href" {
					hrefs = append(hrefs, v.Val)
					break
				}
			}
		}
	}
}

func parseURL(root *url.URL, href string) (*url.URL, error) {
	u, err := url.Parse(href)
	if err != nil {
		return nil, err
	}

	if u.Hostname() != "" && u.Hostname() != root.Hostname() {
		return nil, errExternalURL
	}

	u.Host = root.Host
	u.Scheme = root.Scheme
	return u, nil
}
