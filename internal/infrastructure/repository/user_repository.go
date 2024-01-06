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
	dataPath := "./internal/db/users.json"
	return &UserRepository{
		mu:       mu,
		logger:   logger,
		dataPath: dataPath,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user model.User) error {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	// Read existing users data from the file
	users, err := ur.readUsersFromFile()
	if err != nil {
		return err
	}

	// Generate a unique user ID (e.g., length of the map + 1)
	userID := fmt.Sprint(len(users) + 1)
	users[userID] = user

	// Write the updated users data back to the file
	if err := ur.writeUsersToFile(users); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, id string, user *model.User) error {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	// Read existing users data from the file
	users, err := ur.readUsersFromFile()
	if err != nil {
		return err
	}

	// Check if the user with the specified ID exists
	existingUser, exists := users[id]
	if !exists {
		return fmt.Errorf("error: user with id %s not found", id)
	}

	// Update only the fields that need to be modified
	existingUser.Name = user.Name
	existingUser.Score = user.Score
	existingUser.Answers = user.Answers
	existingUser.FinishedQuiz = user.FinishedQuiz

	// Write the updated users data back to the file
	if err := ur.writeUsersToFile(users); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetUser(ctx context.Context, id string) (*model.User, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	// Read existing users data from the file
	users, err := ur.readUsersFromFile()
	if err != nil {
		return nil, err
	}

	// Retrieve the user from the map
	user, exists := users[id]
	if !exists {
		return nil, fmt.Errorf("error: user with id %s not found", id)
	}

	return &user, nil
}

func (ur *UserRepository) GetAllUsers(ctx context.Context) (model.UserMap, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	// Read existing users data from the file
	users, err := ur.readUsersFromFile()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) readUsersFromFile() (model.UserMap, error) {
	content, err := os.ReadFile(ur.dataPath)
	if err != nil {
		if os.IsNotExist(err) || len(content) == 0 {
			return make(model.UserMap), nil
		}
		return nil, err
	}

	var users model.UserMap
	if err := json.Unmarshal(content, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) writeUsersToFile(users model.UserMap) error {
	content, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ur.dataPath, content, 0644)
}
