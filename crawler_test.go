package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fakeSiteRoot = "/"
const fakeSiteRootHTML = `<!DOCTYPE html>
<html class="no-js " lang="en" dir="ltr">
  <head></head>
  <body>
    <a href="/virus" class="clickbait">Something Really Cool</a>
    <a href="/docs/secret/">Top Secret!</a>
    <a href="https://apple.com">New iPhones 4 U</a>
    <a href="mailto:///email">Contact Us</a>
    <a href="#:">Invalid URL</a>
    <a href=":#">Invalid URL 2</a>
  </body>
</html>`

func serveTestSite(t *testing.T) *httptest.Server {
	r := mux.NewRouter()
	r.HandleFunc(fakeSiteRoot, func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		assert.Equal(t, userAgent, ua)

		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fakeSiteRootHTML))
	}).Methods("GET")

	handleRobotsTxt(r)

	return httptest.NewServer(r)
}

func TestCrawler(t *testing.T) {
	client := &http.Client{}
	server := serveTestSite(t)
	defer server.Close()

	c := newCrawler(client, nil)

	u, err := url.Parse(server.URL + fakeSiteRoot)
	require.NoError(t, err)

	links := c.crawlRoot(newLink(u))
	assert.Len(t, links.Data, 3)
	assert.NotNil(t, links.Data["/"])
	assert.NotNil(t, links.Data["/virus/"])
	assert.NotNil(t, links.Data["/docs/secret/"])
}

func TestCrawler__RobotsTxtUsed(t *testing.T) {
	client := &http.Client{}
	server := serveTestSite(t)
	defer server.Close()

	robots, err := retrieveRobotsTxt(client, server.URL+robotsTxtPath)
	require.NoError(t, err)

	c := newCrawler(client, robots)

	u, err := url.Parse(server.URL + fakeSiteRoot)
	require.NoError(t, err)

	links := c.crawlRoot(newLink(u))
	require.Len(t, links.Data, 2)
	assert.NotNil(t, links.Data["/"])
	assert.NotNil(t, links.Data["/virus/"])
}
