package quiz

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Question struct {
	Question string
	Answer   string
}

const (
	AnswerCorrect   int32 = 0
	AnswerIncorrect int32 = 1
	RanOutOfTime    int32 = 2
)

type AnswerProvider func(chan string)

func ReadQuestions(filename string) []Question {
	var questions []Question

	file, err := os.Open(filename)

	if err != nil {
		log.Fatalf("Error while opening %s file: %s\n", filename, err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		log.Fatalf("Error while reading from %s file: %s\n", filename, err)
	}

	for _, record := range records {
		if len(record) != 2 {
			log.Fatal("Records should be in 'question, answer' format!")
		}

		question := strings.TrimSpace(record[0])
		answer := strings.TrimSpace(record[1])

		questions = append(questions, Question{
			Question: question,
			Answer:   answer,
		})
	}

	return questions
}

func Shuffle(questions []Question, seed int64) []Question {
	newQuestions := append(make([]Question, 0, len(questions)), questions...)
	r := rand.New(rand.NewSource(seed))
	r.Shuffle(len(questions), func(i, j int) {
		newQuestions[i], newQuestions[j] = newQuestions[j], newQuestions[i]
	})
	return newQuestions
}

func PlayQuiz(questions []Question, timeLimit int, answerProvider AnswerProvider) int {
	score := 0

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	fmt.Println("----")
	for _, question := range questions {

		result := askQuestion(question, timer, answerProvider)

		if result == RanOutOfTime {
			fmt.Printf("\nRan out of time. Your score %d\n", score)
			return score
		} else if result == AnswerCorrect {
			score++
		}
	}

	fmt.Printf("Congratulations you got %d out of %d questions right\n", score, len(questions))

	return score
}

func askQuestion(question Question, timer *time.Timer, answerProvider AnswerProvider) int32 {
	fmt.Printf("Q: %s\n$ ", question.Question)

	answerChannel := make(chan string)

	go answerProvider(answerChannel)

	for {
		select {
		case <-timer.C:
			return RanOutOfTime
		case userAnswer := <-answerChannel:
			isCorrect := userAnswer == question.Answer

			fmt.Println("userAnswer=", userAnswer)

			if isCorrect {
				fmt.Printf("A: %s [v]\n----\n", question.Answer)
				return AnswerCorrect
			} else {
				fmt.Printf("A: %s [x]\n----\n", question.Answer)
				return AnswerIncorrect
			}
		}
	}

}

func GetAnswerFromStdin(answerChannel chan string) {
	reader := bufio.NewReader(os.Stdin)
	result, _ := reader.ReadString('\n')
	result = strings.TrimRight(result, "\n")
	result = strings.TrimSpace(result)
	answerChannel <- result
}
