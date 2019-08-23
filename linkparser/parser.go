package linkparser

import "io"

// import (
// 	"golang.org/x/net/html"
// )

// Link represents an html link
type Link struct {
	Href string
	Text string
}

// Parse accepts an io.Reader, parses it and return a slice of Links
func Parse(r io.Reader) []Link {
	var links []Link

	return links
}
