package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/MFCaballero/simple-quiz/internal/domain/model"
)

type QuestionRepository struct {
	mu       *sync.RWMutex
	logger   *log.Logger
	dataPath string
}

func NewQuestionRepository(logger *log.Logger) model.QuestionRepository {
	mu := &sync.RWMutex{}
	dataPath := "./db/questions.json"
	return &QuestionRepository{
		mu:       mu,
		logger:   logger,
		dataPath: dataPath,
	}
}

func (qr *QuestionRepository) GetAllQuestions(ctx context.Context) (model.QuestionMap, error) {
	qr.mu.RLock()
	defer qr.mu.RUnlock()

	questions, err := qr.readQuestionsFromFile()
	if err != nil {
		qr.logger.Printf("error: getting questions: %v", err)
		return nil, err
	}

	return questions, nil
}

func (qr *QuestionRepository) GetQuestion(ctx context.Context, id string) (*model.Question, error) {
	qr.mu.RLock()
	defer qr.mu.RUnlock()

	questions, err := qr.readQuestionsFromFile()
	if err != nil {
		qr.logger.Printf("error: getting question: %v", err)
		return nil, err
	}

	question, exists := questions[id]
	if !exists {
		err = fmt.Errorf("question with id %s not found", id)
		qr.logger.Printf("error: getting question: %v", err)
		return nil, err
	}

	return &question, nil
}

func (qr *QuestionRepository) readQuestionsFromFile() (model.QuestionMap, error) {
	fileContent, err := os.ReadFile(qr.dataPath)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %v", qr.dataPath, err)
	}

	questions := model.QuestionMap{}
	if err := json.Unmarshal(fileContent, &questions); err != nil {
		return nil, fmt.Errorf("decoding json: %v", err)
	}

	return questions, nil
}
