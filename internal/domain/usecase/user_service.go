package usecase

import (
	"log"
	"net/http"

	"github.com/MFCaballero/simple-quiz/internal/domain/model"
)

type UserService struct {
	repository model.UserRepository
	logger     *log.Logger
}

func NewUserService(repo model.UserRepository, logger *log.Logger) *UserService {
	return &UserService{
		repository: repo,
		logger:     logger,
	}
}

func (us *UserService) AnswerQuestion(w http.ResponseWriter, r *http.Request) {

}

func (us *UserService) GetAnswered(w http.ResponseWriter, r *http.Request) {

}

func (us *UserService) PostAnswers(w http.ResponseWriter, r *http.Request) {

}
