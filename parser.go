package main

import (
	"io"

	"golang.org/x/net/html"
)

func parser(r io.Reader) ([]string, error) {
	z := html.NewTokenizer(r)
	for {
		tk := z.Next()
		if tk == html.ErrorToken {
			return nil, z.Err()
		}

		switch tk {
		case html.StartTagToken, html.EndTagToken:
			tn, hasAttr := z.TagName()
			if len(tn) == 1 && tn[0] == 'a' && hasAttr {
				for {
					key, val, more := z.TagAttr()
					if string(key) == "href" {
						break
					}

				}
			}
		}

	}
}
