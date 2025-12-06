package main

/*
BUG: When the timeout occurs, and the user transferred to another question
but still, if the user gives the answer the program stuck (unless we give two answers before next timout occurs).
This shows that previous iteration goroutine is still active and waiting for user input.

Problem Analysis:
When the timer expires, the goroutine doing reader.ReadString('\n') is still blocked, waiting for input forever.
That goroutine never dies, and the next question ends up having an extra goroutine from the previous round still listening on stdin.
*/

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
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

	// 1. Read the CSV file
	filename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 5, "the time limit for the quiz in seconds")
	shuffle := flag.Bool("shuffle", false, "shuffle the quiz questions")
	flag.Parse()

	// Read the file records
	fileRecords := readFile(*filename)
	quizProblems := parseFileRecords(fileRecords, *shuffle)

	// 4. Run the quiz
	correctAnswers, totalQuestions := runQuiz(quizProblems, *timeLimit)
	fmt.Printf("\nQuiz Completed! You scored %d out of %d.\n", correctAnswers, totalQuestions)

}

func readFile(filename string) [][]string {

	// 2. Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open the CSV file: %s\n", filename)
	}
	defer file.Close()

	// 3. Read the file contents
	r := csv.NewReader(file)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read the provided CSV file: %s\n", err)
	}

	return records
}

func parseFileRecords(fileRecords [][]string, shuffle bool) []problem {
	problems := make([]problem, len(fileRecords))
	for i, record := range fileRecords {
		if len(record) < 2 { // if csv line doesn't have at least 2 fields or invalid
			continue // or log, depending on design
		}
		problems[i] = problem{
			question: record[0],
			answer:   record[1],
		}
	}

	if shuffle {
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	return problems
}

func runQuiz(problems []problem, timeLimit int) (correctAnswers int, totalQuestions int) {
	fmt.Println("Welcome to the Quiz!")
	fmt.Println("---------------------")

	answerCh := make(chan string, 1)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			answerCh <- strings.TrimSpace(text)
		}
	}()

	for _, p := range problems {
		fmt.Printf("Question: %s = ?\n", p.question)
		fmt.Print("-> ")

		questionTimer := time.NewTimer(time.Duration(timeLimit) * time.Second)

		select {
		case answer := <-answerCh:
			if answer == p.answer {
				correctAnswers++
				fmt.Println("Correct!")
			} else {
				fmt.Printf("Wrong! The correct answer is %s\n", p.answer)
			}
			totalQuestions++

		case <-questionTimer.C:
			fmt.Println("\n Time's up for this question!")
			totalQuestions++
		}
	}

	return correctAnswers, totalQuestions
}
