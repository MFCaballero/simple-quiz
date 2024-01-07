package model

import "context"

type Option struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	IsCorrect bool   `json:"is_correct"`
}

type Question struct {
	Label   string   `json:"label"`
	Options []Option `json:"options"`
}

type QuestionMap map[string]Question

type QuestionRepository interface {
	GetAllQuestions(ctx context.Context) (QuestionMap, error)
	GetQuestion(ctx context.Context, id string) (*Question, error)
}
