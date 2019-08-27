package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Links contains a map of all links found on a site
var Links map[string]Link

// Link represents an individual link found on a site
type Link struct {
	URL     string
	Visited bool
}

func main() {
	url := flag.String("url", "http://agajewelrydesigns.com", "a domain of the website we are creating a sitemap for")
	flag.Parse()

	response, err := makeGetCall(*url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(response))
}

func makeGetCall(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
