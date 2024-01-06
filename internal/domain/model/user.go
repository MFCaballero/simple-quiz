package model

import "context"

type User struct {
	Name    string  `json:"name"`
	Score   float32 `json:"score"`
	Answers []struct {
		QuestionID string `json:"question_id"`
		Option     Option `json:"option"`
	} `json:"answers"`
	FinishedQuiz bool `json:"finished_quiz"`
}

type UserMap map[string]User

type UserRepository interface {
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, id string, user *User) error
	GetUser(ctx context.Context, id string) (*User, error)
	GetAllUsers(ctx context.Context) (UserMap, error)
}
