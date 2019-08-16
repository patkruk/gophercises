package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/patkruk/gophercises/urlshort"
	bolt "go.etcd.io/bbolt"
)

var yaml = `
- path: /urlshort
  url: http://www.google.com
- path: /urlshort-final
  url: http://www.wp.pl
`

func main() {
	// default fallback
	mux := defaultMux()

	// // Build the MapHandler using the mux as the fallback
	pathsToUrls, err := getMapOfPathUrls()
	if err != nil {
		log.Fatal(err)
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

	// Build the JSONHandler using the yamlHandler as the fallback
	json := `
[
	{
		"path": "/json-godoc",
		"url": "https://golang.org/pkg/encoding/json/"
	}
]
`
	jsonHandler, err := urlshort.JSONHandler([]byte(json), yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
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

func getMapOfPathUrls() (map[string]string, error) {
	// to be returned
	result := make(map[string]string)

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my_bbolt.db", 0600, nil)
	if err != nil {
		return result, err
	}
	defer db.Close()

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	bucketName := "urls"

	// store some data in a bucket (only if it does not exist yet)
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		for key, val := range pathsToUrls {
			value := bucket.Get([]byte(key))
			// only add if it does not exist yet
			if value == nil {
				fmt.Println("Adding new data to the bucket")

				err = bucket.Put([]byte(key), []byte(val))
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return result, err
	}

	// retrieve the data
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("Bucket %s not found", bucketName)
		}

		for key := range pathsToUrls {
			value := bucket.Get([]byte(key))
			if value == nil {
				return fmt.Errorf("Unable to fetch data for %s", key)
			}

			result[key] = string(value)
		}

		return nil
	})
	if err != nil {
		return result, err
	}

	return result, nil
}
