package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/patkruk/gophercises/urlshort"
)

var yaml = `
- path: /urlshort
  url: http://www.google.com
- path: /urlshort-final
  url: http://www.wp.pl
`

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	yml := []byte(yaml)

	fileName := flag.String("yml", "", "a yml file'")
	flag.Parse()

	// if a yaml file specified, use it
	if isFlagPassed("yml") {
		file, err := os.Open(*fileName)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		yml, err = ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Build the YAMLHandler using the mapHandler as the fallback
	yamlHandler, err := urlshort.YAMLHandler(yml, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)

	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func isFlagPassed(name string) bool {
	found := false

	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})

	return found
}
