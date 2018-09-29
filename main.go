package main

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := cli.App("crawler", "Command line tool which crawls a given web domain")
	rootURL := app.String(cli.StringOpt{
		Name:   "root-url",
		Value:  "https://monzo.com/",
		Desc:   "Root URL from which we begin crawling",
		EnvVar: "ROOT_URL",
	})

	log.SetLevel(log.InfoLevel)

	app.Action = func() {
		start := time.Now()
		log.WithField("rootURL", *rootURL).Info("Starting web crawler")

		client := &http.Client{}
		root, err := url.Parse(*rootURL)
		if err != nil {
			log.WithError(err).Fatal("Failed to parse provided url")
		}

		rootLink := &Link{URL: root, Children: make(map[string]*Link)}
		c := newCrawler(client)
		c.crawlRoot(rootLink)

		log.WithField("total", len(c.links.Data)).WithField("duration", time.Now().Sub(start).String()).Info("Crawling complete")
	}

	log.SetLevel(log.InfoLevel)
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("App could not start, error=[%s]\n", err)
		return
	}
}
