// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/domain/model/question.go

// Package mock_model is a generated GoMock package.
package mock_model

import (
	context "context"
	reflect "reflect"

	model "github.com/MFCaballero/simple-quiz/internal/domain/model"
	gomock "github.com/golang/mock/gomock"
)

// MockQuestionRepository is a mock of QuestionRepository interface.
type MockQuestionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockQuestionRepositoryMockRecorder
}

// MockQuestionRepositoryMockRecorder is the mock recorder for MockQuestionRepository.
type MockQuestionRepositoryMockRecorder struct {
	mock *MockQuestionRepository
}

// NewMockQuestionRepository creates a new mock instance.
func NewMockQuestionRepository(ctrl *gomock.Controller) *MockQuestionRepository {
	mock := &MockQuestionRepository{ctrl: ctrl}
	mock.recorder = &MockQuestionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuestionRepository) EXPECT() *MockQuestionRepositoryMockRecorder {
	return m.recorder
}

// GetAllQuestions mocks base method.
func (m *MockQuestionRepository) GetAllQuestions(ctx context.Context) (model.QuestionMap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllQuestions", ctx)
	ret0, _ := ret[0].(model.QuestionMap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllQuestions indicates an expected call of GetAllQuestions.
func (mr *MockQuestionRepositoryMockRecorder) GetAllQuestions(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllQuestions", reflect.TypeOf((*MockQuestionRepository)(nil).GetAllQuestions), ctx)
}

// GetQuestion mocks base method.
func (m *MockQuestionRepository) GetQuestion(ctx context.Context, id string) (*model.Question, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuestion", ctx, id)
	ret0, _ := ret[0].(*model.Question)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuestion indicates an expected call of GetQuestion.
func (mr *MockQuestionRepositoryMockRecorder) GetQuestion(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuestion", reflect.TypeOf((*MockQuestionRepository)(nil).GetQuestion), ctx, id)
}