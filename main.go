package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	limit := flag.Int("limit", 30, "the time limit for the quiz in seconds (default 30)")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	defer file.Close()

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	problems := parseLines(lines)
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})

	score := 0
	ansChan := make(chan string)
	timeout := time.After(time.Second * time.Duration(*limit))

label:
	for i, p := range problems {

		go func() {
			var ans string
			fmt.Printf("Problem #%d: %s = ", i+1, p.question)
			fmt.Scan(&ans)
			ansChan <- ans
		}()

		select {
		case ans := <-ansChan:
			if strings.TrimSpace(ans) == p.answer {
				score++
			}
		case <-timeout:
			break label
		}

	}

	fmt.Println()
	fmt.Printf("You scored %d out of %d", score, len(problems))

}
func parseLines(lines [][]string) []problem {
	p := make([]problem, len(lines))
	for i, l := range lines {
		p[i] = problem{
			question: l[0],
			answer:   l[1],
		}
	}
	return p
}

func askQuestion(i int, p problem) {
	var ans string
	fmt.Printf("Problem #%d: %s = ", i, p.question)
	fmt.Scan(&ans)

}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
