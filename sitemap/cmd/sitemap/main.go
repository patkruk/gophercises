package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/patkruk/gophercises/linkparser"
)

// Links contains a map of all links found on a site
var Links map[string]Link

// Link represents an individual link found on a site
type Link struct {
	URL     string
	Visited bool
}

func main() {
	links := make(map[string]Link)

	domain := flag.String("domain", "http://agajewelrydesigns.com", "a domain of the website we are creating a sitemap for")
	flag.Parse()

	// call the specified domain
	response, err := makeGetCall(*domain)
	if err != nil {
		log.Fatal(err)
	}

	// parse all links on the site
	urls := parseLinks(response)

	// add new links to the map of links
	for _, link := range urls {
		// skip links to the homepage
		if link.Href == "/" || link.Href == *domain || link.Href == "" {
			continue
		}

		// skip mailto links

		// prepend domain if missing
		href := link.Href
		if string(href[0]) == "/" {
			href = *domain + href
		}

		// skip links to external sites

		// check if link already exists before setting
		// we don't want to overwrite

		links[href] = Link{
			URL:     href,
			Visited: false,
		}
	}

	fmt.Println(links)
}

func parseLinks(r io.Reader) []linkparser.Link {
	return linkparser.Parse(r)
}

func makeGetCall(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// defer resp.Body.Close()

	return resp.Body, nil
}
