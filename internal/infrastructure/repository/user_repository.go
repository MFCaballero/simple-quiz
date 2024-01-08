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

type UserRepository struct {
	mu       *sync.RWMutex
	logger   *log.Logger
	dataPath string
}

func NewUserRepository(logger *log.Logger) model.UserRepository {
	mu := &sync.RWMutex{}
	dataPath := "./db/users.json"
	return &UserRepository{
		mu:       mu,
		logger:   logger,
		dataPath: dataPath,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	users, err := ur.readUsersFromFile()
	if err != nil {
		ur.logger.Printf("error: creating user: %v", err)
		return nil, err
	}

	userID := fmt.Sprint(len(users) + 1)
	user.ID = userID
	users[userID] = user

	if err := ur.writeUsersToFile(users); err != nil {
		ur.logger.Printf("error: creating user: %v", err)
		return nil, err
	}
	return &user, err
}

func (ur *UserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	users, err := ur.readUsersFromFile()
	if err != nil {
		ur.logger.Printf("error: updating user: %v", err)
		return err
	}

	users[user.ID] = *user
	if err := ur.writeUsersToFile(users); err != nil {
		ur.logger.Printf("error: updating user: %v", err)
		return err
	}

	return nil
}

func (ur *UserRepository) GetUser(ctx context.Context, id string) (*model.User, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	users, err := ur.readUsersFromFile()
	if err != nil {
		ur.logger.Printf("error: getting user: %v", err)
		return nil, err
	}

	user, exists := users[id]
	if !exists {
		err = fmt.Errorf("user with id %s not found", id)
		ur.logger.Printf("error: getting user: %v", err)
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetAllUsers(ctx context.Context) (model.UserMap, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	users, err := ur.readUsersFromFile()
	if err != nil {
		ur.logger.Printf("error: getting all users: %v", err)
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) readUsersFromFile() (model.UserMap, error) {
	content, err := os.ReadFile(ur.dataPath)

	if err != nil {
		if os.IsNotExist(err) {
			return make(model.UserMap), nil
		}
		return nil, fmt.Errorf("reading users from file: % v", err)
	}

	if len(content) == 0 {
		return make(model.UserMap), nil
	}

	users := model.UserMap{}
	if err := json.Unmarshal(content, &users); err != nil {
		return nil, fmt.Errorf("decoding content to usermap: % v", err)
	}

	return users, nil
}

func (ur *UserRepository) writeUsersToFile(users model.UserMap) error {
	content, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return fmt.Errorf("writting users to file: %v", err)
	}

	if err := os.WriteFile(ur.dataPath, content, 0644); err != nil {
		return fmt.Errorf("writing users to file: %v", err)
	}
	return nil
}
