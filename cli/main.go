package main

import (
	"github.com/MFCaballero/simple-quiz/cli/commands"
	"github.com/MFCaballero/simple-quiz/cli/config"
	"github.com/MFCaballero/simple-quiz/cli/session"
	"github.com/spf13/cobra"
)

func main() {
	sessionManager := session.NewSessionManager()
	config := config.LoadConfig()
	rootCmd := &cobra.Command{Use: "quiz"}
	rootCmd.AddCommand(commands.LoginCommand(sessionManager, config)...)
	rootCmd.AddCommand(commands.QuestionCommand(sessionManager, config))
	rootCmd.AddCommand(commands.AnswerCommand(sessionManager, config))

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
