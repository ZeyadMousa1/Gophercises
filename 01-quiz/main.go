package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const problemsFile = "problems.csv"

var (
	correctAnswers int
	totalQuestions int 
)

func main() {
	var(
		flagProblemsFile = flag.String("p", problemsFile, "the path of the problems CSV file")
		flagTimer = flag.Duration("t", 30 * time.Second, "the max time for the quiz")
	)
	flag.Parse()

	if flagProblemsFile == nil || flagTimer == nil {
		fmt.Println("Missing problems file name and/or timer")
		return;
	}

	fmt.Printf("Hit enter to start quiz from %q in %v\n", *flagProblemsFile, *flagTimer)
	fmt.Scanln()

	// read csv file (problems.csv)
	f, err := os.Open(*flagProblemsFile)
	if  err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return;
	} 
	defer f.Close();

	r := csv.NewReader(f)
	questions, err := r.ReadAll()
	totalQuestions = len(questions)
	if err != nil {
		fmt.Printf("failed to read file: %v\n", err)
		return;
	}

	// start the quizs
	quizDone := startQuiz(questions)

	// define the timer
	quizTimer := time.Tick(*flagTimer)

	// wait for quiz timer or quiz done 
	select{
	case <- quizDone:
	case <- quizTimer:
	}

	// output number of questions (total + correct)
	fmt.Printf("Result: %d/%d\n", correctAnswers, totalQuestions)
}

func startQuiz(questions [][]string) chan bool{
	done := make(chan bool)
	go func (){
		for i, question := range questions {
			// diplay one question at a time
			question, correctAnswer := question[0], question[1]
			fmt.Printf("%d. %s?\n", i+1, question)
	
			// get answer from user, then proceed to next question immediately
			var answer string
			_, err := fmt.Scan(&answer)
			if err != nil {
				fmt.Printf("failed to scan %v\n", err)
				return;
			}
			// clean up answers 
			// - by removingextra white space
			answer = strings.TrimSpace(answer)
			// - lower caseing the answer to avoid capitalization
			answer = strings.ToLower(answer)
			if answer == correctAnswer {
				correctAnswers ++
			}
		}
		done <- true;
	}()
	return done;
}