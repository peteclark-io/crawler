package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	req, err := http.NewRequest("GET", "https://monzo.com", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	hrefs, err := parser(resp.Body)
	require.NoError(t, err)
	assert.NotEmpty(t, hrefs)
}

func BenchmarkParser(b *testing.B) {
	req, err := http.NewRequest("GET", "https://monzo.com", nil)
	if err != nil {
		b.Error(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		b.Error(err)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	r := bytes.NewBuffer(body)

	for n := 0; n < b.N; n++ {
		_, err = parser(r)
		if err != nil {
			b.Error(err)
			continue
		}
	}
}

func TestParseURL(t *testing.T) {
	root, _ := url.Parse("https://monzo.com/blog/")
	u, err := parseURL(root, "latest")
	require.NoError(t, err)

	assert.Equal(t, root.Scheme, u.Scheme)
	assert.Equal(t, root.Host, u.Host)
	assert.Equal(t, root.Path+"latest/", u.Path)
	assert.True(t, strings.HasSuffix(u.Path, "/"), "We require a trailing slash")
}

func TestParseURL__InvalidURL(t *testing.T) {
	root, _ := url.Parse("https://monzo.com/blog/")
	_, err := parseURL(root, ":#")

	assert.EqualError(t, err, "parse :: missing protocol scheme")
}

func TestParseURL__ErrorsForExternalURLs(t *testing.T) {
	root, _ := url.Parse("https://monzo.com/blog/")
	_, err := parseURL(root, "https://github.com/")

	assert.EqualError(t, err, errExternalURL.Error())
}

func TestParseURL__ErrorsForNonHTTPSchemes(t *testing.T) {
	root, _ := url.Parse("https://monzo.com/blog/")
	_, err := parseURL(root, "mailto://git@monzo.com")

	assert.EqualError(t, err, errUnsupportedScheme.Error())
}
