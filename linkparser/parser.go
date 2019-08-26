package linkparser

import (
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

// Link represents an html link
type Link struct {
	Href string
	Text string
}

// Parse accepts a Reader, parses it and return a slice of Links
func Parse(r io.Reader) []Link {
	var links []Link

	// parse the Reader and return a parse tree
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var href, text string

			// determine the href
			for _, a := range n.Attr {
				if a.Key == "href" {
					href = a.Val
					break
				}
			}

			// determine the text
			text = getText(n)

			link := Link{
				Href: href,
				Text: text,
			}

			links = append(links, link)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links
}

func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	// if n.Type != html.ElementNode {
	// 	return ""
	// }

	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += getText(c)
	}

	return strings.Join(strings.Fields(ret), " ")
}
