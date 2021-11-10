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

func runQuiz(quiz []QuizElem, score *int, questionTimer time.Duration, finished chan<- int) {
	totalQuestions := len(quiz)
	ch := make(chan string)
	for i := 0; i < totalQuestions; i++ {
		fmt.Printf("Question #%v: %v = ", i, quiz[i].Question)
		go func() {
			answ := ""
			fmt.Scan(&answ)
			ch <- answ
		}()
		select {
		case answer := <-ch:
			if answer == quiz[i].Answer {
				*score += 1
			}
		case <-time.After(questionTimer * time.Second):
			fmt.Println()
		}
	}
	close(ch)
	finished <- 1
	close(finished)
}

func main() {
	docName := flag.String("path", "problems.csv", "Path to quiz")
	setQuestionTimer := flag.Int("questionTimer", 5, "Number of seconds")
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

	finished := make(chan int)

	go runQuiz(quiz, &correct, time.Duration(*setQuestionTimer), finished)
	<-finished
	fmt.Printf("\nYou scored %v out of %v\n", correct, totalQuestions)
}
