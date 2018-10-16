package pkg

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

func TestLinks__AddRootLink(t *testing.T) {
	l := NewLinks()
	u, err := parseURL(monzoURL, "/blog/latest/1/")
	require.NoError(t, err)

	expectedLink := NewLink(u)

	ok := l.add(nil, expectedLink)

	assert.False(t, ok)
	assert.Len(t, l.Data, 1)

	actualLink, ok := l.Data["/blog/latest/1/"]
	require.True(t, ok)
	assert.Equal(t, expectedLink.URL, actualLink.URL)
}

func TestLinks__AddLink(t *testing.T) {
	l := NewLinks()
	u, err := parseURL(monzoURL, "/blog/latest/1/")
	require.NoError(t, err)

	parentLink := NewLink(monzoURL)
	expectedLink := NewLink(u)

	ok := l.add(parentLink, expectedLink)

	assert.False(t, ok)
	assert.Len(t, l.Data, 1)

	actualLink, ok := l.Data["/blog/latest/1/"]
	require.True(t, ok)
	assert.Equal(t, expectedLink.URL, actualLink.URL)

	actualLink, ok = parentLink.Children["/blog/latest/1/"]
	require.True(t, ok, "Tests the parent has been updated with the new child")
	assert.Equal(t, expectedLink.URL, actualLink.URL)

	actualParentLink, ok := expectedLink.Parents[monzoURL.Path]
	require.True(t, ok, "Tests the child has been updated with the new parent")
	assert.Equal(t, parentLink.URL, actualParentLink.URL)
}

func TestLinks__AddLinkChecksParentAndChildMatch(t *testing.T) {
	l := NewLinks()

	parentLink := NewLink(monzoURL)
	matchingLink := NewLink(monzoURL)

	ok := l.add(parentLink, matchingLink)
	assert.True(t, ok)
	assert.Len(t, l.Data, 0)
}

func TestLinks__AddLinkIsIdempotent(t *testing.T) {
	l := NewLinks()

	u, err := parseURL(monzoURL, "/blog/latest/1/")
	require.NoError(t, err)

	parentLink := NewLink(monzoURL)
	expectedLink := NewLink(u)

	ok := l.add(parentLink, expectedLink)
	assert.False(t, ok)
	assert.Len(t, l.Data, 1)

	ok = l.add(parentLink, expectedLink)
	assert.True(t, ok)
	assert.Len(t, l.Data, 1)
}
