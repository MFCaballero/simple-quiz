package main

import (
	"log"
	"os"
	"sync"

	"github.com/MFCaballero/simple-quiz/internal/domain/usecase"
	"github.com/MFCaballero/simple-quiz/internal/infrastructure/api"
	"github.com/MFCaballero/simple-quiz/internal/infrastructure/repository"
)

func main() {
	logger := log.New(os.Stdout, "[Quiz Logger] ", log.Ldate|log.Ltime)
	wg := &sync.WaitGroup{}
	userRepository := repository.NewUserRepository(logger)
	questionRepository := repository.NewQuestionRepository(logger)
	services := usecase.LoadServices(userRepository, questionRepository, logger)
	app := api.NewApp(logger, wg, services)
	go app.ListenForErrors()
	go app.ListenForShutdown()
	app.Run()
}
