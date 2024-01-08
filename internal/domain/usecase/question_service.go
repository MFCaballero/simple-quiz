package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MFCaballero/simple-quiz/internal/domain/model"
	"github.com/go-chi/chi/v5"
)

type QuestionService struct {
	repository model.QuestionRepository
	logger     *log.Logger
}

func NewQuestionService(repo model.QuestionRepository, logger *log.Logger) *QuestionService {
	return &QuestionService{
		repository: repo,
		logger:     logger,
	}
}

func (qs *QuestionService) GetAllQuestions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMessage := "An error occured getting all questions"

	questions, err := qs.repository.GetAllQuestions(ctx)
	if err != nil {
		http.Error(w, errMessage, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(qs.toQuestionsDTO(questions)); err != nil {
		qs.logger.Printf("error encoding questions to json: %v", err)
		http.Error(w, errMessage, http.StatusInternalServerError)
		return
	}
}

func (qs *QuestionService) GetQuestion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "question")
	errMessage := fmt.Sprintf("An error occured getting question with id %s", id)

	question, err := qs.repository.GetQuestion(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(qs.toQuestionDTO(question)); err != nil {
		qs.logger.Printf("error encoding question to json: %v", err)
		http.Error(w, errMessage, http.StatusInternalServerError)
		return
	}
}

type QuestionDTO struct {
	Label   string      `json:"label"`
	Options []OptionDTO `json:"options"`
}
type OptionDTO struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

func (qs *QuestionService) toQuestionsDTO(questions model.QuestionMap) map[string]QuestionDTO {
	questionsMap := make(map[string]QuestionDTO, len(questions))
	for id, question := range questions {
		questionsMap[id] = qs.toQuestionDTO(&question)
	}
	return questionsMap
}

func (qs *QuestionService) toQuestionDTO(question *model.Question) QuestionDTO {
	optionsDTO := make([]OptionDTO, len(question.Options))
	for i, option := range question.Options {
		optionsDTO[i] = OptionDTO{
			ID:    option.ID,
			Label: option.Label,
		}
	}
	return QuestionDTO{
		Label:   question.Label,
		Options: optionsDTO,
	}
}
