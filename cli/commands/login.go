package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/MFCaballero/simple-quiz/cli/config"
	"github.com/MFCaballero/simple-quiz/cli/session"
	"github.com/spf13/cobra"
)

func LoginCommand(sessionManager *session.SessionManager, config config.Config) []*cobra.Command {
	login := &cobra.Command{
		Use:   "login",
		Short: "Login to the quiz app",
		Run: func(cmd *cobra.Command, args []string) {
			name, err := cmd.Flags().GetString("userName")
			if err != nil {
				log.Fatal(err)
			}
			session, err := sessionManager.GetSession()
			if err != nil {
				log.Fatal(err)
			}
			if session != nil {
				log.Fatalf("Already logged user %s", session.Name)
			}
			userID, err := loginUser(name, config.BackendURL)
			if err != nil {
				log.Fatal(err)
			}

			if err := sessionManager.CreateSession(userID, name); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Welcome: %s!", name)
		},
	}
	login.Flags().StringP("userName", "u", "", "Your user name")
	login.MarkFlagRequired("userName")
	logout := &cobra.Command{
		Use:   "logout",
		Short: "Logout to the quiz app",
		Run: func(cmd *cobra.Command, args []string) {
			if err := sessionManager.DeleteSession(); err != nil {
				log.Fatal(err)
			}
			fmt.Print("Good bye!")
		},
	}
	return []*cobra.Command{login, logout}
}

func loginUser(name, url string) (string, error) {
	body := loginRequest{Name: name}
	bodyBytes, err := json.Marshal(&body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal login request body: %v", err)
	}
	response, err := http.Post(url+"/users/login", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to send login request: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("login failed with status code: %d", response.StatusCode)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read login response body: %v", err)
	}

	var loginResp loginResponse
	if err := json.Unmarshal(responseBody, &loginResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal login response: %v", err)
	}

	return loginResp.UserID, nil
}

type loginRequest struct {
	Name string `json:"name"`
}

type loginResponse struct {
	UserID string `json:"user_id"`
}
