package model

import "context"

type User struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Score        float32  `json:"score"`
	Answers      []Answer `json:"answers"`
	FinishedQuiz bool     `json:"finished_quiz"`
}
type Answer struct {
	QuestionID string `json:"question_id"`
	Option     Option `json:"option"`
}

type UserMap map[string]User

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id string) (*User, error)
	GetAllUsers(ctx context.Context) (UserMap, error)
}
