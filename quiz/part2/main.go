package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// CsvLine represents each line in the csv file
type CsvLine struct {
	Question string
	Answer   string
}

// Results has the correct and wrong counts
type Results struct {
	TotalCount   int
	CorrectCount int
	WrongCount   int
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")

	flag.Parse()

	file, err := os.Open(*fileName)
	check(err)

	defer file.Close()

	// read the file and parse it (csv)
	lines, err := csv.NewReader(file).ReadAll()
	check(err)

	results := Results{
		TotalCount: len(lines),
	}

	// create a timer
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	for _, line := range lines {
		data := CsvLine{
			Question: line[0],
			Answer:   line[1],
		}

		// ask the question
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%v = ?\n", data.Question)

		// read the answer
		answerCh := make(chan string)
		// trigger a goroutine
		go func() {
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSuffix(answer, "\n")

			// send the answer over the channel
			answerCh <- answer
		}()

		select {
		case <-timer.C: // if the timer times out, stop and print the results
			printResults(&results)
			return
		case answer := <-answerCh: // if you have an answer from the user, check it
			// keep track of right & wrong answers
			if answer == data.Answer {
				results.CorrectCount++
			} else {
				results.WrongCount++
			}
		}
	}

	printResults(&results)
}

func printResults(results *Results) {
	fmt.Printf("\nYou provided %d correct answers out of %d", results.CorrectCount, results.TotalCount)
}
