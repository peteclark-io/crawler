package main

import "net/url"

type Link struct {
	URL      *url.URL         `json:"url"`
	Children map[string]*Link `json:"children"`
}

type Links map[string]*Link

func (l *Links) add(link *Link) {

}
