package pkg

import (
	"io"
	"net/url"
)

// DisplayLinks Displays links in a flat parent -> child style
func DisplayLinks(w io.Writer, links *Links) {
	for _, l := range links.Data {
		if len(l.Children) > 0 {
			displayParent(w, l, "")
			for _, child := range l.Children {
				displayChild(w, child, "  ")
			}
		} else {
			displayChild(w, l, "")
		}
	}
}

// Displays links in a cascading style from the root URL
func displayLinkWithChildren(w io.Writer, link *Link, spacer string) {
	if len(link.Children) > 0 {
		displayParent(w, link, spacer)
		for _, child := range link.Children {
			displayLinkWithChildren(w, child, spacer+"  ")
		}
	} else {
		displayChild(w, link, spacer)
	}
}

func displayParent(w io.Writer, link *Link, spacer string) {
	w.Write([]byte(spacer + "- " + normalize(link.URL) + " ->\r\n"))
}

func displayChild(w io.Writer, link *Link, spacer string) {
	w.Write([]byte(spacer + "- " + normalize(link.URL) + "\r\n"))
}

func normalize(u *url.URL) string {
	return u.Scheme + "://" + u.Host + u.Path
}
