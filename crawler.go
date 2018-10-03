package main

import (
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

type crawler struct {
	client *http.Client
	links  *Links
	wg     *sync.WaitGroup
}

func newCrawler(client *http.Client) *crawler {
	links := newLinks()
	return &crawler{client: client, links: links, wg: &sync.WaitGroup{}}
}

func (c *crawler) crawlRoot(root *Link) *Links {
	c.links.add(nil, root)

	err := c.crawlLink(root)
	if err != nil {
		log.WithError(err).Error("Failed to crawl root link")
	}

	c.wg.Wait()

	return c.links
}

func (c *crawler) crawlLink(link *Link) error {
	log.WithField("url", link.URL.String()).Info("crawling new link")

	req, err := http.NewRequest("GET", link.URL.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	hrefs, err := parser(resp.Body)
	if err != nil {
		return err
	}

	for _, href := range hrefs {
		child, err := parseURL(link.URL, href)
		if err != nil {
			continue
		}

		l := &Link{URL: child, Parents: make(map[string]*Link), Children: make(map[string]*Link)}
		ok := c.links.add(link, l)

		if !ok {
			c.wg.Add(1)
			go func() {
				c.crawlLink(l)
				c.wg.Done()
			}()
		}
	}

	return nil
}
