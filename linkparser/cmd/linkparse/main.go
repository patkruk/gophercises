package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/patkruk/gophercises/linkparser"
)

func main() {
	fileName := flag.String("file", "ex1.html", "an html file to be read and parsed")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	links := linkparser.Parse(file)

	fmt.Println(links)
}
