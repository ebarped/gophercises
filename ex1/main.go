package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

type Exercise struct {
	question string
	solution string
}

func (e Exercise) CheckSolution(response string) bool {
	if e.solution != response {
		return false
	}
	return true
}

type Quiz struct {
	ex    []Exercise
	score uint
}

// Shuffle randomizes the order of exercises of the quiz
func (q Quiz) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(q.ex), func(i, j int) {
		q.ex[i], q.ex[j] = q.ex[j], q.ex[i]
	})
}

// GetAnswer gets the answer from a question to the user & sends the response to answers channel
func GetAnswer(e Exercise, answers chan string) {
	var response string
	fmt.Scanln(&response)
	answers <- response
}

// LoadQuiz loads a quiz from a .csv file
func LoadQuiz(file string) Quiz {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	var quiz Quiz
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		exercise := Exercise{
			question: record[0],
			solution: record[1],
		}
		quiz.ex = append(quiz.ex, exercise)
	}
	return quiz
}

func main() {
	csvFile := flag.String("csvFile", "problems.csv", "path to csv file")
	timeout := flag.Duration("timeout", time.Duration(30*time.Second), "duration of the exam")
	shuffle := flag.Bool("shuffle", false, "shuffle the questions")
	flag.Parse()

	fmt.Println("Exam file: " + *csvFile)
	fmt.Printf("Duration of the quiz: %s\n", *timeout)
	fmt.Printf("Shuffle quiz?: %t\n", *shuffle)

	quiz := LoadQuiz(*csvFile)

	if *shuffle {
		quiz.Shuffle()
	}

	fmt.Println("Press [Enter] to begin the quiz!")
	fmt.Scanln()

	answers := make(chan string) // channel to receive answers to exercises from user

	timer := time.NewTimer(*timeout) // will publish a value in timer.C when the timeout is passed

	for i, ex := range quiz.ex {
		fmt.Printf("Q%d: %s\n", i+1, ex.question)
		fmt.Printf("A: ")

		go GetAnswer(ex, answers) //ask for an answer in a non-blocking way

		select {
		case <-timer.C:
			fmt.Printf("\nExam finished!\n")
			fmt.Printf("Your score is %d/%d\n", quiz.score, len(quiz.ex))
			return
		case answer := <-answers:
			if ex.CheckSolution(answer) {
				quiz.score++
			}
		}
	}

}
