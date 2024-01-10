package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MFCaballero/simple-quiz/internal/domain/model"
	mock_model "github.com/MFCaballero/simple-quiz/internal/domain/model/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_model.NewMockUserRepository(ctrl)

	userService := NewUserService(mockUserRepo, nil, nil)
	ctx := context.Background()

	t.Run("Login Success", func(t *testing.T) {

		mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(&model.User{ID: "1"}, nil)
		requestBody := map[string]string{"name": "John Doe"}
		jsonBody, _ := json.Marshal(requestBody)

		req, err := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(userService.Login).ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
		expectedResponseBody := map[string]string{"user_id": "1"}
		var responseBody map[string]string
		err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponseBody, responseBody)
	})

	t.Run("Login Failure - Bad Request", func(t *testing.T) {
		invalidRequestBody := []byte("invalid_json")

		req, err := http.NewRequest("POST", "/users/login", bytes.NewBuffer(invalidRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		http.HandlerFunc(userService.Login).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Login Failure - Internal Server Error", func(t *testing.T) {
		mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil, errors.New("Internal Server Error"))
		validRequestBody := map[string]string{"name": "John Doe"}
		jsonBody, _ := json.Marshal(validRequestBody)

		req, err := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		http.HandlerFunc(userService.Login).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func setupRouterAndRequest(
	t *testing.T,
	handler http.HandlerFunc,
	method, path, reqURL string,
	body []byte,
) *httptest.ResponseRecorder {
	r := chi.NewRouter()
	if r == nil {
		r = chi.NewRouter()
	}
	r.HandleFunc(path, handler)
	req, err := http.NewRequest(method, reqURL, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	return rr
}

func TestGetAnswered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_model.NewMockUserRepository(ctrl)
	mockQuestionRepo := mock_model.NewMockQuestionRepository(ctrl)

	userService := NewUserService(mockUserRepo, mockQuestionRepo, nil)
	mockUserID := "1"
	mockUser := &model.User{
		ID: mockUserID,
		Answers: []model.Answer{
			{QuestionID: "1", Option: model.Option{ID: "A", Label: "Option A"}},
			{QuestionID: "2", Option: model.Option{ID: "B", Label: "Option B"}},
		},
	}
	mockQuestions := model.QuestionMap{
		"1": {Label: "Question 1"},
		"2": {Label: "Question 2"},
	}
	t.Run("GetAnswered Success", func(t *testing.T) {

		mockUserRepo.EXPECT().GetUser(gomock.Any(), mockUserID).Return(mockUser, nil)
		mockQuestionRepo.EXPECT().GetAllQuestions(gomock.Any()).Return(mockQuestions, nil)

		rr := setupRouterAndRequest(t, userService.GetAnswered, "GET", "/users/{user}/answered", fmt.Sprintf("/users/%s/answered", mockUserID), nil)

		assert.Equal(t, http.StatusOK, rr.Code)

		expectedResponseBody := []Answer{
			{Question: "Question 1", QuestionID: "1", Option: "Option A", OptionID: "A"},
			{Question: "Question 2", QuestionID: "2", Option: "Option B", OptionID: "B"},
		}
		var responseBody []Answer
		err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponseBody, responseBody)
	})

	t.Run("GetAnswered Failure - User Not Found", func(t *testing.T) {
		mockUserID := "nonExistentUserID"
		mockUserRepo.EXPECT().GetUser(gomock.Any(), mockUserID).Return(nil, errors.New("User not found"))

		rr := setupRouterAndRequest(t, userService.GetAnswered, "GET", "/users/{user}/answered", fmt.Sprintf("/users/%s/answered", mockUserID), nil)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("GetAnswered Failure - Internal Server Error", func(t *testing.T) {
		mockUserID := "1"
		mockUserRepo.EXPECT().GetUser(gomock.Any(), mockUserID).Return(mockUser, nil)
		mockQuestionRepo.EXPECT().GetAllQuestions(gomock.Any()).Return(nil, errors.New("Internal Server Error"))

		rr := setupRouterAndRequest(t, userService.GetAnswered, "GET", "/users/{user}/answered", fmt.Sprintf("/users/%s/answered", mockUserID), nil)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestAnswerQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_model.NewMockUserRepository(ctrl)
	mockQuestionRepo := mock_model.NewMockQuestionRepository(ctrl)

	userService := NewUserService(mockUserRepo, mockQuestionRepo, nil)
	mockUserID := "1"
	mockUser := &model.User{
		ID: mockUserID,
		Answers: []model.Answer{
			{QuestionID: "1", Option: model.Option{ID: "A", Label: "Option A"}},
		},
	}
	mockAnswerRequest := AnswerRequest{
		QuestionID: "2",
		OptionID:   "B",
	}
	mockQuestions := model.QuestionMap{
		"2": {
			Label: "Question 2",
			Options: []model.Option{
				{ID: "B", Label: "Option B", IsCorrect: true},
			},
		},
	}

	t.Run("AnswerQuestion Success", func(t *testing.T) {

		mockUserRepo.EXPECT().GetUser(gomock.Any(), mockUserID).Return(mockUser, nil)
		mockQuestionRepo.EXPECT().GetAllQuestions(gomock.Any()).Return(mockQuestions, nil)
		mockUserRepo.EXPECT().UpdateUser(gomock.Any(), mockUser).Return(nil)

		reqBody, err := json.Marshal(mockAnswerRequest)
		assert.NoError(t, err)
		rr := setupRouterAndRequest(t, userService.AnswerQuestion, "POST", "/users/{user}/answer", fmt.Sprintf("/users/%s/answer", mockUserID), reqBody)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("AnswerQuestion Failure - User Not Found", func(t *testing.T) {
		mockUserID := "nonExistentUserID"

		mockUserRepo.EXPECT().GetUser(gomock.Any(), mockUserID).Return(nil, errors.New("User not found"))

		reqBody, err := json.Marshal(mockAnswerRequest)
		assert.NoError(t, err)
		rr := setupRouterAndRequest(t, userService.AnswerQuestion, "POST", "/users/{user}/answer", fmt.Sprintf("/users/%s/answer", mockUserID), reqBody)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("AnswerQuestion Failure - Internal Server Error", func(t *testing.T) {

		mockUserRepo.EXPECT().GetUser(gomock.Any(), mockUserID).Return(mockUser, nil)
		mockQuestionRepo.EXPECT().GetAllQuestions(gomock.Any()).Return(nil, errors.New("Internal Server Error"))

		reqBody, err := json.Marshal(mockAnswerRequest)
		assert.NoError(t, err)

		rr := setupRouterAndRequest(t, userService.AnswerQuestion, "POST", "/users/{user}/answer", fmt.Sprintf("/users/%s/answer", mockUserID), reqBody)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("AnswerQuestion Failure - User Already Finished Quiz", func(t *testing.T) {

		mockUser.FinishedQuiz = true
		mockUserRepo.EXPECT().GetUser(gomock.Any(), mockUserID).Return(mockUser, nil)

		reqBody, err := json.Marshal(mockAnswerRequest)
		assert.NoError(t, err)

		rr := setupRouterAndRequest(t, userService.AnswerQuestion, "POST", "/users/{user}/answer", fmt.Sprintf("/users/%s/answer", mockUserID), reqBody)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("AnswerQuestion Failure - Invalid Option", func(t *testing.T) {
		mockUser.FinishedQuiz = false
		mockAnswerRequest := AnswerRequest{
			QuestionID: "2",
			OptionID:   "C",
		}

		mockUserRepo.EXPECT().GetUser(gomock.Any(), mockUserID).Return(mockUser, nil)
		mockQuestionRepo.EXPECT().GetAllQuestions(gomock.Any()).Return(mockQuestions, nil)
		mockUserRepo.EXPECT().UpdateUser(gomock.Any(), mockUser).Return(nil)

		reqBody, err := json.Marshal(mockAnswerRequest)
		assert.NoError(t, err)

		rr := setupRouterAndRequest(t, userService.AnswerQuestion, "POST", "/users/{user}/answer", fmt.Sprintf("/users/%s/answer", mockUserID), reqBody)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestGetScoreData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_model.NewMockUserRepository(ctrl)
	mockQuestionRepo := mock_model.NewMockQuestionRepository(ctrl)

	userService := NewUserService(mockUserRepo, mockQuestionRepo, nil)
	mockUserID := "1"
	mockUser := &model.User{
		ID:           mockUserID,
		FinishedQuiz: true,
		Score:        0.75,
		Answers: []model.Answer{
			{QuestionID: "1", Option: model.Option{ID: "A", Label: "Option A", IsCorrect: true}},
			{QuestionID: "2", Option: model.Option{ID: "B", Label: "Option B", IsCorrect: false}},
		},
	}
	mockUsers := model.UserMap{
		mockUserID: *mockUser,
		"2":        {ID: "2", FinishedQuiz: true, Score: 0.5},
		"3":        {ID: "3", FinishedQuiz: true, Score: 0.9},
	}
	mockQuestions := model.QuestionMap{
		"1": {Label: "Question 1"},
		"2": {Label: "Question 2"},
	}
	t.Run("GetScoreData Success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetAllUsers(gomock.Any()).Return(mockUsers, nil)
		mockQuestionRepo.EXPECT().GetAllQuestions(gomock.Any()).Return(mockQuestions, nil)

		rr := setupRouterAndRequest(t, userService.GetScoreData, "GET", "/users/{user}/score", fmt.Sprintf("/users/%s/score", mockUserID), nil)

		assert.Equal(t, http.StatusOK, rr.Code)

		expectedResponseBody := ScoreData{
			Score:               0.75,
			TotalQuestions:      2,
			CorrectAnswers:      1,
			BetterThan:          0.5,
			RelativePerformance: 0.07142859,
			AnswersDetail: []AnswersDetail{
				{Question: "Question 1", Answer: "Option A", IsCorrect: true},
				{Question: "Question 2", Answer: "Option B", IsCorrect: false},
			},
		}
		var responseBody ScoreData
		err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponseBody, responseBody)
	})

	t.Run("GetScoreData Failure - User Not Found", func(t *testing.T) {
		mockUserID := "nonExistentUserID"
		mockUsers := model.UserMap{}
		mockUserRepo.EXPECT().GetAllUsers(gomock.Any()).Return(mockUsers, nil)

		rr := setupRouterAndRequest(t, userService.GetScoreData, "GET", "/users/{user}/score", fmt.Sprintf("/users/%s/score", mockUserID), nil)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("GetScoreData Failure - User Not Finished Quiz", func(t *testing.T) {
		mockUsers := model.UserMap{
			mockUserID: {ID: mockUserID, FinishedQuiz: false},
		}
		mockUserRepo.EXPECT().GetAllUsers(gomock.Any()).Return(mockUsers, nil)

		rr := setupRouterAndRequest(t, userService.GetScoreData, "GET", "/users/{user}/score", fmt.Sprintf("/users/%s/score", mockUserID), nil)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("GetScoreData Failure - Internal Server Error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetAllUsers(gomock.Any()).Return(mockUsers, errors.New("Internal Server Error"))

		rr := setupRouterAndRequest(t, userService.GetScoreData, "GET", "/users/{user}/score", fmt.Sprintf("/users/%s/score", mockUserID), nil)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestPostAnswers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_model.NewMockUserRepository(ctrl)
	mockQuestionRepo := mock_model.NewMockQuestionRepository(ctrl)

	userService := NewUserService(mockUserRepo, mockQuestionRepo, nil)

	mockQuestions := model.QuestionMap{
		"1": {Label: "Question 1", Options: []model.Option{{ID: "A", IsCorrect: true}, {ID: "B", IsCorrect: false}}},
		"2": {Label: "Question 2", Options: []model.Option{{ID: "A", IsCorrect: false}, {ID: "B", IsCorrect: true}}},
	}

	t.Run("PostAnswers Success", func(t *testing.T) {
		mockUserID := "1"
		mockUser := &model.User{
			ID: mockUserID,
			Answers: []model.Answer{
				{QuestionID: "1", Option: model.Option{ID: "A", Label: "Option A", IsCorrect: true}},
				{QuestionID: "2", Option: model.Option{ID: "B", Label: "Option B", IsCorrect: false}},
			},
		}
		mockUserRepo.EXPECT().GetUser(gomock.Any(), mockUserID).Return(mockUser, nil)
		mockQuestionRepo.EXPECT().GetAllQuestions(gomock.Any()).Return(mockQuestions, nil)
		mockUserRepo.EXPECT().UpdateUser(gomock.Any(), mockUser).Return(nil)

		rr := setupRouterAndRequest(t, userService.PostAnswers, "POST", "/users/{user}/finish", fmt.Sprintf("/users/%s/finish", mockUserID), nil)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "Quiz completed successfully!", rr.Body.String())
		assert.True(t, mockUser.FinishedQuiz)
		assert.Equal(t, float32(0.5), mockUser.Score) // 1 correct answer out of 2 questions
	})

	t.Run("PostAnswers Failure - User Not Found", func(t *testing.T) {
		mockUserID := "nonExistentUserID"
		mockUserRepo.EXPECT().GetUser(gomock.Any(), mockUserID).Return(nil, errors.New("User not found"))

		rr := setupRouterAndRequest(t, userService.PostAnswers, "POST", "/users/{user}/finish", fmt.Sprintf("/users/%s/finish", mockUserID), nil)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("PostAnswers Failure - Internal Server Error", func(t *testing.T) {
		mockUserID := "2"
		mockUser := &model.User{
			ID: mockUserID,
			Answers: []model.Answer{
				{QuestionID: "1", Option: model.Option{ID: "A", Label: "Option A", IsCorrect: true}},
				{QuestionID: "2", Option: model.Option{ID: "B", Label: "Option B", IsCorrect: false}},
			},
		}
		mockUserRepo.EXPECT().GetUser(gomock.Any(), mockUserID).Return(mockUser, nil)
		mockQuestionRepo.EXPECT().GetAllQuestions(gomock.Any()).Return(nil, errors.New("Internal Server Error"))
		rr := setupRouterAndRequest(t, userService.PostAnswers, "POST", "/users/{user}/finish", fmt.Sprintf("/users/%s/finish", mockUserID), nil)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.False(t, mockUser.FinishedQuiz)
		assert.Equal(t, float32(0), mockUser.Score)
	})

	t.Run("PostAnswers Failure - Missing Questions", func(t *testing.T) {
		mockUserID := "3"
		mockUser := &model.User{
			ID: mockUserID,
			Answers: []model.Answer{
				{QuestionID: "1", Option: model.Option{ID: "A", Label: "Option A", IsCorrect: true}},
			},
		}
		mockUserRepo.EXPECT().GetUser(gomock.Any(), mockUserID).Return(mockUser, nil)
		mockQuestionRepo.EXPECT().GetAllQuestions(gomock.Any()).Return(mockQuestions, nil)

		rr := setupRouterAndRequest(t, userService.PostAnswers, "POST", "/users/{user}/finish", fmt.Sprintf("/users/%s/finish", mockUserID), nil)

		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Equal(t, "Missing questions to answer before finishing\n", rr.Body.String())
		assert.False(t, mockUser.FinishedQuiz)
		assert.Equal(t, float32(0), mockUser.Score)
	})
}
