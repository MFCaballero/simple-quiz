package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/MFCaballero/simple-quiz/internal/domain/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	wait          *sync.WaitGroup
	logger        *log.Logger
	services      usecase.Services
	errorChan     chan error
	errorChanDone chan bool
}

func NewApp(logger *log.Logger, wg *sync.WaitGroup, services usecase.Services) App {
	return App{
		wait:          wg,
		logger:        logger,
		services:      services,
		errorChan:     make(chan error),
		errorChanDone: make(chan bool),
	}
}

func (app *App) ListenForErrors() {
	for {
		select {
		case err := <-app.errorChan:
			app.logger.Printf("error: %v", err)
		case <-app.errorChanDone:
			return
		}
	}
}

func (app *App) ListenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *App) shutdown() {

	app.wait.Wait()

	app.errorChanDone <- true

	app.logger.Println("closing channels and shutting down application...")

	close(app.errorChan)
	close(app.errorChanDone)
}

func (app *App) Run() {
	port := 8080 //this should be in config
	addr := fmt.Sprintf(":%d", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: app.routes(),
	}

	app.logger.Printf("starting web server on port %d", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *App) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	mux.Get("/questions", app.services.QuestionService.GetAllQuestions)
	mux.Get("/question/{id}", app.services.QuestionService.GetQuestion)
	return mux
}
