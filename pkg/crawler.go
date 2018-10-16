package pkg

import (
	"net/http"
	"sync"

	"github.com/samclarke/robotstxt"
	log "github.com/sirupsen/logrus"
)

const userAgent = "peteclark-io/crawler"

type crawler struct {
	client    *http.Client
	links     *Links
	wg        *sync.WaitGroup
	robotsTxt *robotstxt.RobotsTxt
}

func NewCrawler(client *http.Client, r *robotstxt.RobotsTxt) *crawler {
	links := NewLinks()
	return &crawler{client: client, links: links, robotsTxt: r, wg: &sync.WaitGroup{}}
}

func (c *crawler) CrawlRoot(root *Link) *Links {
	c.links.add(nil, root)

	err := c.crawlLink(root)
	if err != nil {
		log.WithError(err).Error("Failed to crawl root link")
	}

	c.wg.Wait()
	return c.links
}

func (c *crawler) checkLinkAgainstRobotsTxt(link *Link) bool {
	if c.robotsTxt == nil {
		return true // if no robots txt, go ahead and crawl anyway
	}

	ok, err := c.robotsTxt.IsAllowed(userAgent, link.URL.String())
	if err != nil {
		log.WithError(err).WithField("url", link.URL.String()).Warn("Robots txt failed to check link")
		return true // crawl anyway if robots.txt failed to check the link
	}
	return ok
}

func (c *crawler) crawlLink(link *Link) error {
	log.WithField("url", link.URL.String()).Info("crawling link")

	req, err := http.NewRequest("GET", link.URL.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", userAgent)

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

		l := NewLink(child)

		ok := c.checkLinkAgainstRobotsTxt(l)
		if !ok {
			log.WithField("url", link.URL.String()).Info("not allowed to crawl link")
			continue
		}

		ok = c.links.add(link, l)

		if !ok {
			c.wg.Add(1)
			go func() {
				defer c.wg.Done()
				c.crawlLink(l)
			}()
		}
	}

	return nil
}
