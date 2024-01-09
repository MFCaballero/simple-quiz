package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MFCaballero/simple-quiz/cli/config"
	"github.com/MFCaballero/simple-quiz/cli/session"
	"github.com/spf13/cobra"
)

func QuestionCommand(sessionManager *session.SessionManager, config config.Config) *cobra.Command {
	var questionCmd = &cobra.Command{
		Use:   "question",
		Short: "Interact with quiz questions",
	}

	questionCmd.AddCommand(ListQuestionsCommand(config))
	questionCmd.AddCommand(GetQuestionCommand(config))

	return questionCmd
}

func ListQuestionsCommand(config config.Config) *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all quiz questions",
		Run: func(cmd *cobra.Command, args []string) {
			questions, err := listQuestions(config.BackendURL)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("List of Quiz Questions:")
			for i := 1; i <= len(questions); i++ {
				fmt.Printf("%d) %v\n", i, questions[strconv.Itoa(i)].Label)
			}
		},
	}

	return listCmd
}

func GetQuestionCommand(config config.Config) *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a quiz question options",
		Run: func(cmd *cobra.Command, args []string) {
			questionNumber, err := cmd.Flags().GetString("questionNumber")
			if err != nil {
				log.Fatal(err)
			}
			question, err := getQuestion(config.BackendURL, questionNumber)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%s) %s\n", questionNumber, question.Label)
			fmt.Println("Options:")
			for _, option := range question.Options {
				fmt.Printf("%s %s\n", option.ID, option.Label)
			}
		},
	}
	getCmd.Flags().StringP("questionNumber", "n", "", "Question number")
	getCmd.MarkFlagRequired("questionNumber")
	return getCmd
}

type question struct {
	Label   string `json:"label"`
	Options []struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	} `json:"options"`
}

func listQuestions(url string) (map[string]question, error) {
	resp, err := http.Get(url + "/questions")
	if err != nil {
		return nil, fmt.Errorf("error getting questions: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	questions := map[string]question{}
	if err := json.NewDecoder(resp.Body).Decode(&questions); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return questions, nil
}

func getQuestion(url, id string) (*question, error) {
	resp, err := http.Get(url + "/questions/" + id)
	if err != nil {
		return nil, fmt.Errorf("error getting question: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("question not found with ID: %s", id)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	question := &question{}
	if err := json.NewDecoder(resp.Body).Decode(question); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return question, nil
}
