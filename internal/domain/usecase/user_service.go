package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MFCaballero/simple-quiz/internal/domain/model"
	"github.com/go-chi/chi/v5"
)

type UserService struct {
	userRepo     model.UserRepository
	questionRepo model.QuestionRepository
	logger       *log.Logger
}

func NewUserService(userRepo model.UserRepository, questionRepo model.QuestionRepository, logger *log.Logger) *UserService {
	return &UserService{
		userRepo:     userRepo,
		questionRepo: questionRepo,
		logger:       logger,
	}
}

func (us *UserService) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	errMessage := "An error occured logging user"
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	newUser, err := us.userRepo.CreateUser(r.Context(), model.User{
		Name: loginRequest.Name,
	})
	if err != nil {
		http.Error(w, errMessage, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"user_id": newUser.ID}); err != nil {
		us.logger.Printf("error encoding to json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (us *UserService) GetAnswered(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user")
	user, err := us.userRepo.GetUser(r.Context(), userID)
	errMessage := "An error occured getting user's answers"

	if err != nil {
		http.Error(w, errMessage, http.StatusNotFound)
		return
	}
	questions, err := us.questionRepo.GetAllQuestions(r.Context())
	if err != nil {
		http.Error(w, errMessage, http.StatusInternalServerError)
		return
	}
	response := make([]Answer, len(user.Answers))
	for _, answer := range user.Answers {
		i, err := strconv.Atoi(answer.QuestionID)
		if err != nil {
			http.Error(w, errMessage, http.StatusInternalServerError)
			return
		}
		response[i-1] = Answer{
			Question:   questions[answer.QuestionID].Label,
			QuestionID: answer.QuestionID,
			Option:     answer.Option.Label,
			OptionID:   answer.Option.ID,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		us.logger.Printf("error encoding user %s answers to json: %v", userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (us *UserService) PostAnswers(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user")
	user, err := us.userRepo.GetUser(r.Context(), userID)
	errMessage := "An error occured posting user's answers"

	if err != nil {
		http.Error(w, errMessage, http.StatusNotFound)
		return
	}

	questions, err := us.questionRepo.GetAllQuestions(r.Context())
	if err != nil {
		http.Error(w, errMessage, http.StatusInternalServerError)
		return
	}
	if len(user.Answers) != len(questions) {
		http.Error(w, "Missing questions to answer before finishing", http.StatusForbidden)
		return
	}
	user.FinishedQuiz = true
	totalQuestions := len(questions)
	totalCorrectAnswers := 0
	for _, answer := range user.Answers {
		for _, option := range questions[answer.QuestionID].Options {
			if option.ID == answer.Option.ID && option.IsCorrect {
				totalCorrectAnswers++
				break
			}
		}
	}

	user.Score = float32(totalCorrectAnswers) / float32(totalQuestions)

	if err := us.userRepo.UpdateUser(r.Context(), user); err != nil {
		http.Error(w, errMessage, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Quiz completed successfully!"))
}

func (us *UserService) AnswerQuestion(w http.ResponseWriter, r *http.Request) {
	var answerRequest AnswerRequest
	errMessage := "An error occured answering question"

	if err := json.NewDecoder(r.Body).Decode(&answerRequest); err != nil {
		http.Error(w, errMessage, http.StatusBadRequest)
		return
	}

	userID := chi.URLParam(r, "user")
	user, err := us.userRepo.GetUser(r.Context(), userID)

	if err != nil {
		http.Error(w, errMessage, http.StatusNotFound)
		return
	}

	if user.FinishedQuiz {
		http.Error(w, "User has already finished the quiz", http.StatusForbidden)
		return
	}

	questions, err := us.questionRepo.GetAllQuestions(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve questions", http.StatusInternalServerError)
		return
	}
	question, ok := questions[answerRequest.QuestionID]
	if !ok {
		http.Error(w, errMessage, http.StatusBadRequest)
		return
	}
	answers := []model.Answer{}
	for _, answer := range user.Answers {
		if answer.QuestionID != answerRequest.QuestionID {
			answers = append(answers, answer)
		}
	}
	var isOptionValid bool
	for _, option := range question.Options {
		if option.ID == answerRequest.OptionID {
			answers = append(answers, model.Answer{
				QuestionID: answerRequest.QuestionID,
				Option: model.Option{
					ID:        answerRequest.OptionID,
					Label:     option.Label,
					IsCorrect: option.IsCorrect,
				}})
			isOptionValid = true
		}
	}
	if !isOptionValid {
		http.Error(w, errMessage, http.StatusBadRequest)
	}
	user.Answers = answers
	if err := us.userRepo.UpdateUser(r.Context(), user); err != nil {
		http.Error(w, errMessage, http.StatusInternalServerError)
		return
	}
}

func (us *UserService) GetScoreData(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user")
	errMessage := "An error occured getting user's score data"

	users, err := us.userRepo.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, errMessage, http.StatusInternalServerError)
		return
	}

	user, ok := users[userID]
	if !ok {
		http.Error(w, errMessage, http.StatusNotFound)
		return
	}
	if !user.FinishedQuiz {
		http.Error(w, "user has not finished quiz", http.StatusForbidden)
		return
	}
	totalUsers := len(users)
	betterThanCount := 0
	for _, otherUser := range users {
		if otherUser.ID != userID && otherUser.Score < user.Score {
			betterThanCount++
		}
	}
	betterThan := float32(betterThanCount) / float32(totalUsers)

	questions, err := us.questionRepo.GetAllQuestions(r.Context())
	if err != nil {
		http.Error(w, errMessage, http.StatusInternalServerError)
		return
	}
	scoreData := ScoreData{
		Score:          fmt.Sprintf("%.0f%%", user.Score*100),
		TotalQuestions: len(questions),
		CorrectAnswers: int(user.Score * float32(len(questions))),
		BetterThan:     fmt.Sprintf("%.0f%%", betterThan*100),
	}
	for _, answer := range user.Answers {
		scoreData.AnswersDetail = append(scoreData.AnswersDetail, AnswersDetail{
			Question:  questions[answer.QuestionID].Label,
			Answer:    answer.Option.Label,
			IsCorrect: answer.Option.IsCorrect,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(scoreData); err != nil {
		us.logger.Printf("error encoding score data to json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type LoginRequest struct {
	Name string `json:"name"`
}
type AnswerRequest struct {
	QuestionID string `json:"question_id"`
	OptionID   string `json:"option_id"`
}
type ScoreData struct {
	Score          string          `json:"score"`
	TotalQuestions int             `json:"total_questions"`
	CorrectAnswers int             `json:"correct_answers"`
	BetterThan     string          `json:"better_than"`
	AnswersDetail  []AnswersDetail `json:"answers_detail"`
}
type AnswersDetail struct {
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	IsCorrect bool   `json:"is_correct"`
}

type Answer struct {
	Question   string `json:"question"`
	QuestionID string `json:"question_id"`
	Option     string `json:"option"`
	OptionID   string `json:"option_id"`
}
