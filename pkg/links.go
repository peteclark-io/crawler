package pkg

import (
	"log"
	"net/url"
	"sync"
)

type Link struct {
	sync.RWMutex
	Parents  map[string]*Link `json:"parents"`
	URL      *url.URL         `json:"url"`
	Children map[string]*Link `json:"children"`
}

type Links struct {
	sync.RWMutex
	Data map[string]*Link
}

func NewLinks() *Links {
	return &Links{Data: make(map[string]*Link)}
}

func NewLink(uri *url.URL) *Link {
	return &Link{URL: uri, Parents: make(map[string]*Link), Children: make(map[string]*Link)}
}

func (l *Links) add(parent *Link, link *Link) bool {
	if link.URL == nil || l.Data == nil {
		log.Fatal("Please provide a link containing a URL and/or properly initialise the Links struct")
	}

	if parent == nil {
		l.addLink(link) // short circuit early, as this must be the root url
		return false
	}

	if parent.URL.Path == link.URL.Path {
		return true
	}

	parent.addChild(link)
	link.addParent(parent)

	if l.linkAlreadyProcessed(link) {
		return true
	}

	l.addLink(link)
	return false
}

func (l *Links) addLink(link *Link) {
	l.Lock()
	defer l.Unlock()

	l.Data[link.URL.Path] = link
}

func (l *Links) linkAlreadyProcessed(link *Link) bool {
	l.RLock()
	defer l.RUnlock()

	_, ok := l.Data[link.URL.Path]
	return ok
}

func (l *Link) addParent(parent *Link) {
	l.Lock()
	defer l.Unlock()

	l.Parents[parent.URL.Path] = parent
}

func (l *Link) addChild(child *Link) {
	l.Lock()
	defer l.Unlock()

	l.Children[child.URL.Path] = child
}
