package usecase

import (
	"log"

	"github.com/MFCaballero/simple-quiz/internal/domain/model"
)

type Services struct {
	*UserService
	*QuestionService
}

func LoadServices(userRepo model.UserRepository, questionRepo model.QuestionRepository, logger *log.Logger) Services {
	return Services{
		UserService:     NewUserService(userRepo, questionRepo, logger),
		QuestionService: NewQuestionService(questionRepo, logger),
	}
}
