package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/patkruk/gophercises/cyoa"
)

func main() {
	port := flag.Int("port", 3000, "the port to start the CYOA web application on")
	fileName := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}

	story, err := cyoa.JSONStory(file)
	if err != nil {
		log.Fatal(err)
	}

	h := cyoa.NewHandler(story) // use the defaults
	// tpl := template.Must(template.New("").Parse("Hello!"))
	// h := cyoa.NewHandler(story, cyoa.WithTemplate(tpl))
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
