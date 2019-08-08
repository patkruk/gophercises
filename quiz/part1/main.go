package main

import (
	"flag"
	"io/ioutil"
	"log"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// see if the filename has been provided
	csv := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")

	// once all flags are declared, call flag.Parse() to execute the command-line parsing
	flag.Parse()

	// read the file
	file, err := ioutil.ReadFile(*csv)
	check(err)

	// read individual lines

	// parse csv

	// ask the question

	// store the answer

	// keep track of right & wrong answers
}
