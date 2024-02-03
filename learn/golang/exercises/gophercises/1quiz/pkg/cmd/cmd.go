package cmd

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/GRO4T/quiz/pkg/quiz"
	"github.com/spf13/cobra"
)

var (
	problemsFilename string
	timeLimit        int
	shuffle          bool

	rootCmd = &cobra.Command{
		Use:     "quiz",
		Short:   "CLI quiz application",
		Version: "0.1.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			filename, _ := cmd.Flags().GetString("file")
			timeLimit, _ := cmd.Flags().GetInt("time")

			questions := quiz.ReadQuestions(filename)

			if shuffle {
				questions = quiz.Shuffle(questions, time.Now().UnixNano())
			}

			fmt.Printf(
				"Welcome to the quiz.\n"+
					"Your task is to answer %d questions.\n"+
					"You will have %d seconds to do that.\n"+
					"Press enter to start the game.\n",
				len(questions),
				timeLimit,
			)

			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')

			quiz.PlayQuiz(questions, timeLimit, quiz.GetAnswerFromStdin)

			return nil
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&problemsFilename, "file", "f", "problems.csv", "CSV file with question-answer pairs")
	rootCmd.PersistentFlags().IntVarP(&timeLimit, "time", "t", 30, "Question's time limit in seconds")
	rootCmd.PersistentFlags().BoolVarP(&shuffle, "shuffle", "s", false, "When set, questions will be randomly shuffled")
}
