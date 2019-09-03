package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	urlPackage "net/url"

	"github.com/patkruk/gophercises/linkparser"
)

// Link represents an individual link found on a site
type Link struct {
	URL     string
	Visited bool
}

func main() {
	links := make(map[string]Link)

	domain := flag.String("domain", "http://agajewelrydesigns.com", "a full domain of the website we are creating a sitemap for")
	flag.Parse()

	if !isValidURL(*domain) {
		log.Fatal("Invalid domain provided")
	}

	// call the specified domain
	resp, err := http.Get(*domain)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// parse all links on the site
	urls := parseLinks(resp.Body)

	// add new links to the map of links
	for _, link := range urls {
		if link.Href == "" {
			continue
		}

		href := link.Href

		// skip links to the homepage
		if href == "/" || href == *domain {
			continue
		}

		// skip mailto links
		if len(href) >= 6 && href[0:6] == "mailto" {
			continue
		}

		// prepend domain if missing
		if string(href[0]) == "/" {
			href = *domain + href
		}

		// skip invalid urls
		if !isValidURL(href) {
			continue
		}

		// skip links to external sites
		strippedDomain, err := stripProtocolFromURL(*domain)
		if err != nil {
			log.Fatal(err)
		}
		strippedHref, err := stripProtocolFromURL(href)
		if err != nil {
			log.Fatal(err)
		}
		if len(strippedDomain) > len(strippedHref) || strippedHref[0:len(strippedDomain)] != strippedDomain {
			continue
		}

		// store the link only if it does not exist
		if _, ok := links[href]; !ok {
			links[href] = Link{
				URL:     href,
				Visited: false,
			}
		}
	}

	fmt.Printf("Number of links: %d\n", len(links))
	fmt.Println(links)
	fmt.Println("\nNot finished")
}

func parseLinks(r io.Reader) []linkparser.Link {
	return linkparser.Parse(r)
}

// stripProtocolFromURL returns a "naked" URL without the protocol information
// (e.g. http://www.example.com => example.com)
func stripProtocolFromURL(URL string) (string, error) {
	if len(URL) < 4 {
		return "", fmt.Errorf("Provided URL is less than required 4 characters")
	}

	// remove "http" & "https"
	if len(URL) >= 5 && URL[0:5] == "https" {
		URL = URL[5:]
	} else if URL[0:4] == "http" {
		URL = URL[4:]
	}

	// remove "://"
	if len(URL) >= 3 && URL[0:3] == "://" {
		URL = URL[3:]
	}

	// remove "www."
	if len(URL) >= 4 && URL[0:4] == "www." {
		URL = URL[4:]
	}

	return URL, nil
}

// isValidURL indicates if a given url is valid
func isValidURL(url string) bool {
	_, err := urlPackage.ParseRequestURI(url)
	if err != nil {
		return false
	}

	return true
}
