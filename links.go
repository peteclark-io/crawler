package main

import (
	"log"
	"net/url"
	"sync"
)

type Link struct {
	sync.RWMutex
	Parent   *Link            `json:"parent"`
	URL      *url.URL         `json:"url"`
	Children map[string]*Link `json:"children"`
}

type Links struct {
	sync.RWMutex
	Data map[string]*Link
}

func (l *Links) add(link *Link) bool {
	if link.URL == nil || l.Data == nil {
		log.Fatal("Please provide a link containing a URL and/or properly initialise the Links struct")
	}

	if l.linkAlreadyProcessed(link) {
		return true
	}

	l.Lock()
	defer l.Unlock()

	l.Data[link.URL.Path] = link

	if link.Parent != nil {
		link.Parent.addChild(link)
	}

	return false
}

func (l *Links) linkAlreadyProcessed(link *Link) bool {
	l.RLock()
	defer l.RUnlock()

	_, ok := l.Data[link.URL.Path]
	return ok
}

func (l *Link) addChild(child *Link) {
	l.Lock()
	defer l.Unlock()

	l.Children[child.URL.Path] = child
}
