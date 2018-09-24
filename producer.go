package main

import (
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func crawler(client *http.Client, root *url.URL) ([]*url.URL, error) {
	log.WithField("url", root.String()).Info("crawling new link")

	req, err := http.NewRequest("GET", root.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	hrefs, err := parser(resp.Body)
	if err != nil {
		return nil, err
	}

	links := make([]*url.URL, 0)

	for _, href := range hrefs {
		child, err := parseURL(root, href)
		if err != nil {
			continue
		}

		l, err := crawler(client, child)
		if err != nil {
			continue
		}

		links = append(links, l...)
	}

	return links, nil
}
