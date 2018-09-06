package main

import (
	"net/http"
	"os"

	"github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := cli.App("crawler", "Command line tool which crawls a given web domain")
	rootURL := app.String(cli.StringOpt{
		Name:   "root-url",
		Value:  "https://monzo.com",
		Desc:   "Root URL from which we begin crawling",
		EnvVar: "ROOT_URL",
	})

	log.SetLevel(log.InfoLevel)

	app.Action = func() {
		log.WithField("rootURL", *rootURL).Info("Starting web crawler")

		client := &http.Client{}

		robots, err := retrieveRobotsTxt(client, *rootURL+"/robots.txt")
		if err != nil {
			log.WithError(err).WithField("rootURL", *rootURL).Fatal("Failed to retrieve robots.txt, does this domain even exist?")
		}

	}

	log.SetLevel(log.InfoLevel)
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("App could not start, error=[%s]\n", err)
		return
	}
}
