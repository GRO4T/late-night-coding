package quiz

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadQuestions(t *testing.T) {
	// Arrange
	tempFile, _ := os.CreateTemp(t.TempDir(), "problems.csv")
	tempFile.WriteString("\"What country has the highest life expectancy?\", Hong Kong\n")
	tempFile.WriteString("\"Complete the following lyrics - \"\"I should have changed that stupid lock.....\"\"\", I should have made you leave your key")

	// Act
	questions := ReadQuestions(tempFile.Name())

	// Assert
	assert.Equal(t, 2, len(questions), "There should two questions")
	assert.Equal(t, "What country has the highest life expectancy?", questions[0].Question, "Question should match")
	assert.Equal(t, "Hong Kong", questions[0].Answer, "Answer should match")
	assert.Equal(t, "Complete the following lyrics - \"I should have changed that stupid lock.....\"", questions[1].Question, "Question should match")
	assert.Equal(t, "I should have made you leave your key", questions[1].Answer, "Answer should match")
}

func TestShuffle(t *testing.T) {
	questions := []Question{
		{"A", ""},
		{"B", ""},
		{"C", ""},
		{"D", ""},
	}

	newQuestions := Shuffle(questions, 123)

	assert.Equal(t, 4, len(newQuestions), "The length after shuffling should stay the same")
	for i := 0; i < len(questions); i++ {
		assert.NotEqual(t, questions[i].Question, newQuestions[i].Question, "Questions should be shuffled")
	}
}

func TestPlayQuizWhenSingleQuestionAndAnswerIsCorrect(t *testing.T) {
	// Arrange
	timeLimit := 30
	questions := []Question{
		{
			Question: "What country has the highest life expectancy?",
			Answer:   "Hong Kong",
		},
	}

	// Act
	score := PlayQuiz(questions, timeLimit, func(answerChannel chan string) {
		answerChannel <- "Hong Kong"
	})

	// Assert
	assert.Equal(t, score, 1, "Score should be 1")
}

func TestPlayQuizWhenSingleQuestionAndAnswerIsIncorrect(t *testing.T) {
	// Arrange
	timeLimit := 30
	questions := []Question{
		{
			Question: "What country has the highest life expectancy?",
			Answer:   "Hong Kong",
		},
	}

	// Act
	score := PlayQuiz(questions, timeLimit, func(answerChannel chan string) {
		answerChannel <- "Tokyo"
	})

	// Assert
	assert.Equal(t, score, 0, "Score should be 0")
}

func TestPlayQuizWhenAnswerIsCorrectButQuestionTimedOut(t *testing.T) {
	// Arrange
	timeLimit := 0
	questions := []Question{
		{
			Question: "What country has the highest life expectancy?",
			Answer:   "Hong Kong",
		},
	}

	// Act
	score := PlayQuiz(questions, timeLimit, func(answerChannel chan string) {
		answerChannel <- "Hong Kong"
	})

	// Assert
	assert.Equal(t, score, 0, "Score should be 0")
}
