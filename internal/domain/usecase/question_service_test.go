package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/MFCaballero/simple-quiz/internal/domain/model"
	mock_model "github.com/MFCaballero/simple-quiz/internal/domain/model/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllQuestions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionRepo := mock_model.NewMockQuestionRepository(ctrl)
	questionService := NewQuestionService(mockQuestionRepo, log.Default())

	t.Run("GetAllQuestions Success", func(t *testing.T) {
		mockQuestions := model.QuestionMap{
			"1": {Label: "Question 1", Options: []model.Option{{ID: "A", Label: "Option A"}}},
			"2": {Label: "Question 2", Options: []model.Option{{ID: "B", Label: "Option B"}}},
		}
		mockQuestionRepo.EXPECT().GetAllQuestions(gomock.Any()).Return(mockQuestions, nil)

		rr := setupRouterAndRequest(t, questionService.GetAllQuestions, "GET", "/questions", "/questions", nil)

		assert.Equal(t, http.StatusOK, rr.Code)

		expectedResponseBody := map[string]QuestionDTO{
			"1": {Label: "Question 1", Options: []OptionDTO{{ID: "A", Label: "Option A"}}},
			"2": {Label: "Question 2", Options: []OptionDTO{{ID: "B", Label: "Option B"}}},
		}
		var responseBody map[string]QuestionDTO
		err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponseBody, responseBody)
	})

	t.Run("GetAllQuestions Failure - Internal Server Error", func(t *testing.T) {
		mockQuestionRepo.EXPECT().GetAllQuestions(gomock.Any()).Return(nil, errors.New("Internal Server Error"))

		rr := setupRouterAndRequest(t, questionService.GetAllQuestions, "GET", "/questions", "/questions", nil)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestGetQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionRepo := mock_model.NewMockQuestionRepository(ctrl)
	questionService := NewQuestionService(mockQuestionRepo, log.Default())

	t.Run("GetQuestion Success", func(t *testing.T) {
		mockQuestionID := "1"
		mockQuestion := &model.Question{
			Label:   "Question 1",
			Options: []model.Option{{ID: "A", Label: "Option A"}},
		}
		mockQuestionRepo.EXPECT().GetQuestion(gomock.Any(), mockQuestionID).Return(mockQuestion, nil)

		rr := setupRouterAndRequest(t, questionService.GetQuestion, "GET", "/questions/{question}", fmt.Sprintf("/questions/%s", mockQuestionID), nil)

		assert.Equal(t, http.StatusOK, rr.Code)

		expectedResponseBody := QuestionDTO{
			Label:   "Question 1",
			Options: []OptionDTO{{ID: "A", Label: "Option A"}},
		}
		var responseBody QuestionDTO
		err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponseBody, responseBody)
	})

	t.Run("GetQuestion Failure - Not Found", func(t *testing.T) {
		mockQuestionID := "nonExistentQuestionID"
		mockQuestionRepo.EXPECT().GetQuestion(gomock.Any(), mockQuestionID).Return(nil, errors.New("Question not found"))

		rr := setupRouterAndRequest(t, questionService.GetQuestion, "GET", "/questions/{question}", fmt.Sprintf("/questions/%s", mockQuestionID), nil)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}
