package pkg

import (
	"errors"
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

var errExternalURL = errors.New("provided href points to an external domain")
var errUnsupportedScheme = errors.New("provided href uses a non-http scheme, which we don't support")

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
	if !strings.HasSuffix(href, "/") {
		href = href + "/"
	}

	u, err := url.Parse(href)
	if err != nil {
		return nil, err
	}

	if u.Hostname() != "" && u.Hostname() != root.Hostname() {
		return nil, errExternalURL
	}

	if u.IsAbs() && u.Scheme != "https" && u.Scheme != "http" {
		return nil, errUnsupportedScheme
	}

	return root.ResolveReference(u), nil
}
