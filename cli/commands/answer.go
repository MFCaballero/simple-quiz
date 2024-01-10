package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"

	"github.com/MFCaballero/simple-quiz/cli/config"
	"github.com/MFCaballero/simple-quiz/cli/session"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type userIDKey string

const userID userIDKey = "userID"

func AnswerCommand(sessionManager *session.SessionManager, config config.Config) *cobra.Command {
	var userCmd = &cobra.Command{
		Use:   "answer",
		Short: "Interact with quiz",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			session, err := sessionManager.GetSession()
			if err != nil {
				log.Fatal(err)
			}
			if session == nil {
				log.Fatal("Command only allowed for logged users")
			}
			ctx := context.WithValue(cmd.Context(), userID, session.ID)
			cmd.SetContext(ctx)
		},
	}

	userCmd.AddCommand(AnswerQuestionCommand(config))
	userCmd.AddCommand(GetAnsweredCommand(config))
	userCmd.AddCommand(FinishQuizCommand(config))
	userCmd.AddCommand(GetScoreCommand(config))

	return userCmd
}

func AnswerQuestionCommand(config config.Config) *cobra.Command {
	var answerCmd = &cobra.Command{
		Use:   "post",
		Short: "Answer a quiz question",
		Run: func(cmd *cobra.Command, args []string) {
			userID := cmd.Context().Value(userID).(string)
			question, err := cmd.Flags().GetString("question")
			if err != nil {
				log.Fatal(err)
			}
			option, err := cmd.Flags().GetString("option")
			if err != nil {
				log.Fatal(err)
			}
			req := answerRequest{
				QuestionID: question,
				OptionID:   option,
			}
			if err := answerQuestion(req, config.BackendURL, userID); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Question answered")
		},
	}

	answerCmd.Flags().StringP("question", "q", "", "Question number")
	answerCmd.Flags().StringP("option", "o", "", "Option letter")
	answerCmd.MarkFlagRequired("question")
	answerCmd.MarkFlagRequired("option")

	return answerCmd
}

func GetAnsweredCommand(config config.Config) *cobra.Command {
	var getAnsweredCmd = &cobra.Command{
		Use:   "list",
		Short: "Get answered questions",
		Run: func(cmd *cobra.Command, args []string) {
			userID := cmd.Context().Value(userID).(string)
			answered, err := getUserAnswers(config.BackendURL, userID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("**** Your Answers List ****")
			for _, answer := range answered {
				fmt.Printf("%s) %s %s: %s\n", answer.QuestionID, answer.Question, answer.OptionID, answer.Option)
			}
		},
	}

	return getAnsweredCmd
}

func FinishQuizCommand(config config.Config) *cobra.Command {
	var finishCmd = &cobra.Command{
		Use:   "finish",
		Short: "Finish the quiz",
		Run: func(cmd *cobra.Command, args []string) {
			userID := cmd.Context().Value(userID).(string)
			if err := finishQuiz(config.BackendURL, userID); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Quiz finished")
		},
	}

	return finishCmd
}

func GetScoreCommand(config config.Config) *cobra.Command {
	var scoreCmd = &cobra.Command{
		Use:   "score",
		Short: "Get user score",
		Run: func(cmd *cobra.Command, args []string) {
			userID := cmd.Context().Value(userID).(string)
			scoreData, err := getUserScoreData(config.BackendURL, userID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("**** Your Quiz Results ****")
			fmt.Printf("Your Score: %.0f%%\n", scoreData.Score*100)
			fmt.Printf("Total Questions: %d\n", scoreData.TotalQuestions)
			fmt.Printf("Total Correct Answered: %d\n", scoreData.CorrectAnswers)
			fmt.Printf("You scored better than %.0f%% of other quizzers\n", scoreData.BetterThan*100)
			var performance string
			if scoreData.RelativePerformance > 0 {
				performance = fmt.Sprintf("%.2f%% better than", scoreData.RelativePerformance*100)
			} else if scoreData.RelativePerformance == 0 {
				performance = "equal to"
			} else {
				performance = fmt.Sprintf("%.2f%% worse than", math.Abs(float64(scoreData.RelativePerformance))*100)
			}
			fmt.Printf("Your score is %s the average score for other quizzers\n", performance)
			fmt.Println("**** Your Answers Details ****")
			for _, answer := range scoreData.AnswersDetail {
				fmt.Printf("Question: %s\n", answer.Question)
				var isCorrectMsg string
				if answer.IsCorrect {
					isCorrectMsg = "is correct"
				} else {
					isCorrectMsg = "is wrong"
				}
				fmt.Printf("Your Answer: %s %s\n", answer.Answer, isCorrectMsg)
			}
		},
	}

	return scoreCmd
}

type answerRequest struct {
	QuestionID string `json:"question_id"`
	OptionID   string `json:"option_id"`
}

type userAnswer struct {
	Question   string `json:"question"`
	QuestionID string `json:"question_id"`
	Option     string `json:"option"`
	OptionID   string `json:"option_id"`
}

type scoreData struct {
	Score               float32 `json:"score"`
	TotalQuestions      int     `json:"total_questions"`
	CorrectAnswers      int     `json:"correct_answers"`
	BetterThan          float32 `json:"better_than"`
	RelativePerformance float32 `json:"relative_performance"`
	AnswersDetail       []struct {
		Question  string `json:"question"`
		Answer    string `json:"answer"`
		IsCorrect bool   `json:"is_correct"`
	} `json:"answers_detail"`
}

func answerQuestion(body answerRequest, url, userID string) error {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/users/%s/answer", url, userID), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("error posting answer: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return processErrorResponse(resp)
	}

	return nil
}

func getUserAnswers(url, userID string) ([]userAnswer, error) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%s/answered", url, userID))
	if err != nil {
		return nil, fmt.Errorf("error getting user's answered questions: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, processErrorResponse(resp)
	}

	var answered []userAnswer
	if err := json.NewDecoder(resp.Body).Decode(&answered); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return answered, nil
}

func finishQuiz(url, userID string) error {
	resp, err := http.Post(fmt.Sprintf("%s/users/%s/finish", url, userID), "", nil)
	if err != nil {
		return fmt.Errorf("error finishing the quiz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return processErrorResponse(resp)
	}

	return nil
}

func getUserScoreData(url, userID string) (*scoreData, error) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%s/score", url, userID))
	if err != nil {
		return nil, fmt.Errorf("error getting user score data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, processErrorResponse(resp)
	}

	data := &scoreData{}
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return data, nil
}

func processErrorResponse(resp *http.Response) error {
	errorMessage, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return fmt.Errorf("unexpected status code: %d, unable to read response body: %v", resp.StatusCode, readErr)
	}

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("%s, your request is invalid", string(errorMessage))
	default:
		return errors.New(string(errorMessage))
	}
}
