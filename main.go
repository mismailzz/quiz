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

type problem struct {
	question string
	answer   string
}

func main() {

	// 1. Read the CSV file
	filename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 5, "the time limit for the quiz in seconds")
	flag.Parse()

	// Read the file records
	fileRecords := readFile(*filename)
	quizProblems := parseFileRecords(fileRecords)

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

func parseFileRecords(fileRecords [][]string) []problem {
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
	return problems
}

func runQuiz(problems []problem, timeLimit int) (correctAnswers int, totalQuestions int) {
	// Quiz logic to be implemented
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the Quiz!")
	fmt.Println("---------------------")

	// Set up the timer
	timeout := time.NewTimer(time.Duration(timeLimit) * time.Second)

	correctAnswers = 0
	totalQuestions = 0
	for _, p := range problems {
		fmt.Printf("Question: %s = ?\n", p.question)
		fmt.Print("-> ")

		answerCh := make(chan string)

		// Start a goroutine to read user input (answers)
		go func() {
			text, _ := reader.ReadString('\n')
			// Remove any trailing newline characters
			text = strings.TrimSpace(text)
			answerCh <- text
		}()

		select {
		case answer := <-answerCh:
			text := answer
			if strings.Compare(p.answer, text) == 0 {
				correctAnswers++
				fmt.Println("Correct!")
			} else {
				fmt.Printf("Wrong! The correct answer is %s\n", p.answer)
			}
			totalQuestions++
		case <-timeout.C:
			fmt.Println("\nTime's up!")
			return correctAnswers, totalQuestions
		}
	}
	return correctAnswers, totalQuestions
}
