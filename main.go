package main

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
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the Quiz!")
	fmt.Println("---------------------")

	for _, p := range problems {
		fmt.Printf("Question: %s = ?\n", p.question)
		fmt.Print("-> ")

		// Flush any leftover buffered input before timing
		reader = bufio.NewReader(os.Stdin)

		answerCh := make(chan string, 1)
		questionTimer := time.NewTimer(time.Duration(timeLimit) * time.Second)

		go func() {
			text, _ := reader.ReadString('\n') // This the problem line
			text = strings.TrimSpace(text)
			answerCh <- text
		}()

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
