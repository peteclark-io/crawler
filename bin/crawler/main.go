package main

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/jawher/mow.cli"
	"github.com/peteclark-io/crawler/pkg"
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

		rootLink := pkg.NewLink(root)
		c := pkg.NewCrawler(client, nil)
		links := c.CrawlRoot(rootLink)

		log.WithField("total", len(links.Data)).WithField("duration", time.Now().Sub(start).String()).Info("Crawling complete")
		pkg.DisplayLinks(os.Stdout, links)
	}

	log.SetLevel(log.InfoLevel)
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("App could not start, error=[%s]\n", err)
		return
	}
}
