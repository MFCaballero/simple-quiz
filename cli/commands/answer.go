package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MFCaballero/simple-quiz/cli/config"
	"github.com/MFCaballero/simple-quiz/cli/session"
	"github.com/spf13/cobra"
)

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
				log.Fatal("command only allowed for logged users")
			}
			ctx := context.WithValue(cmd.Context(), "userID", session.ID)
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
			userID := cmd.Context().Value("userID").(string)
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
			userID := cmd.Context().Value("userID").(string)
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
			userID := cmd.Context().Value("userID").(string)
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
			userID := cmd.Context().Value("userID").(string)
			scoreData, err := getUserScoreData(config.BackendURL, userID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("**** Your Quiz Results ****")
			fmt.Printf("Your Score: %s\n", scoreData.Score)
			fmt.Printf("Total Questions: %d\n", scoreData.TotalQuestions)
			fmt.Printf("Total Correct Answered: %d\n", scoreData.CorrectAnswers)
			fmt.Printf("You socored %s better than all users that have taken the quiz\n", scoreData.BetterThan)
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
	Score          string `json:"score"`
	TotalQuestions int    `json:"total_questions"`
	CorrectAnswers int    `json:"correct_answers"`
	BetterThan     string `json:"better_than"`
	AnswersDetail  []struct {
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
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func getUserScoreData(url, userID string) (scoreData, error) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%s/score", url, userID))
	if err != nil {
		return scoreData{}, fmt.Errorf("error getting user score data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return scoreData{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data scoreData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return scoreData{}, fmt.Errorf("error decoding response: %v", err)
	}

	return data, nil
}
