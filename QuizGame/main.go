package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type QuizElem struct {
	Question string
	Answer   string
}

func runQuiz(quiz []QuizElem, score *int) {
	totalQuestions := len(quiz)
	for i := 0; i < totalQuestions; i++ {
		fmt.Printf("Question #%v: %v = ", i, quiz[i].Question)
		answer := ""
		fmt.Scan(&answer)
		if answer == quiz[i].Answer {
			*score += 1
		}
	}
}

func main() {
	docName := flag.String("path", "problems.csv", "Path to quiz")
	setTimer := flag.Int("timer", 30, "Number of seconds")

	flag.Parse()

	csvFile, _ := os.Open(*docName)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))

	var quiz []QuizElem
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		quiz = append(quiz, QuizElem{
			Question: line[0],
			Answer:   line[1],
		})
	}

	totalQuestions := len(quiz)
	correct := 0

	fmt.Println("Press <ANY KEY> to start the quiz")
	fmt.Scanln()

	go runQuiz(quiz, &correct)
	timer1 := time.NewTimer(time.Duration(*setTimer) * time.Second)
	<-timer1.C

	fmt.Printf("\nYou scored %v out of %v\n", correct, totalQuestions)
}
