package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// CsvLine represents each line in the csv file
type CsvLine struct {
	Question string
	Answer   string
}

// Results has the correct and wrong counts
type Results struct {
	CorrectCount int
	WrongCount   int
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// see if the filename has been provided
	fileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")

	// once all flags are declared, call flag.Parse() to execute the command-line parsing
	flag.Parse()

	file, err := os.Open(*fileName)
	check(err)

	defer file.Close()

	// read the file and parse it (csv)
	lines, err := csv.NewReader(file).ReadAll()
	check(err)

	results := Results{}

	for _, line := range lines {
		data := CsvLine{
			Question: line[0],
			Answer:   line[1],
		}

		// ask the question
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%v = ?\n", data.Question)

		// check the answer (keep track of right & wrong answers)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSuffix(answer, "\n")

		if answer == data.Answer {
			results.CorrectCount++
		} else {
			results.WrongCount++
		}
	}

	fmt.Printf("\nCorrect answers count: %d. Wrong answers count: %d.", results.CorrectCount, results.WrongCount)
}
