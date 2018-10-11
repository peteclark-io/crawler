package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const robotsTxtPath = "/robots.txt"
const testRobotsTxt = `# robotstxt.org/

User-agent: *
Disallow: /docs/
`

func serveRobotsTxt() *httptest.Server {
	r := mux.NewRouter()
	handleRobotsTxt(r)

	return httptest.NewServer(r)
}

func handleRobotsTxt(r *mux.Router) {
	r.HandleFunc(robotsTxtPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testRobotsTxt))
	}).Methods("GET")
}

func TestRobotsTxt__RetrieveOK(t *testing.T) {
	client := &http.Client{}
	server := serveRobotsTxt()
	defer server.Close()

	robots, err := retrieveRobotsTxt(client, server.URL+robotsTxtPath)
	require.NoError(t, err)

	ok, err := robots.IsAllowed(userAgent, server.URL)
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = robots.IsAllowed(userAgent, server.URL+"/docs/blah")
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestRobotsTxt(t *testing.T) {
	client := &http.Client{}

	robots, err := retrieveRobotsTxt(client, "https://monzo.com/robots.txt")
	require.NoError(t, err)

	ok, err := robots.IsAllowed(userAgent, "https://monzo.com/blog/2015/10/30/we-are-ready/")
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestRobotsTxt__InvalidURL(t *testing.T) {
	client := &http.Client{}

	robots, err := retrieveRobotsTxt(client, "#:")
	assert.EqualError(t, err, `Get #:: unsupported protocol scheme ""`)
	assert.Nil(t, robots)
}

func TestRobotsTxt__RequestFails(t *testing.T) {
	client := &http.Client{}

	robots, err := retrieveRobotsTxt(client, ":#")
	assert.EqualError(t, err, `parse :: missing protocol scheme`)
	assert.Nil(t, robots)
}
