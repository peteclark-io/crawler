package main

import "net/url"

type Link struct {
	URL      *url.URL `json:"url"`
	Children []*Link  `json:"children"`
}
