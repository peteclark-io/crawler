package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var monzoURL *url.URL

func init() {
	monzoURL, _ = url.Parse("https://monzo.com")
}

func TestLinks(t *testing.T) {
	l := newLinks()
	u, err := parseURL(monzoURL, "/blog/latest/1/")
	require.NoError(t, err)

	expectedLink := &Link{URL: u}

	l.add(nil, expectedLink)

	assert.Len(t, l.Data, 1)

	actualLink, ok := l.Data["/blog/latest/1/"]
	require.True(t, ok)
	assert.Equal(t, expectedLink.URL, actualLink.URL)
}
