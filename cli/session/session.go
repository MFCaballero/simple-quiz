package session

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type UserSession struct {
	ID   string
	Name string
}

type SessionManager struct {
	mu       *sync.RWMutex
	session  *UserSession
	dataPath string
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		mu:       &sync.RWMutex{},
		dataPath: "./db/session.json",
	}
}

func (sm *SessionManager) CreateSession(id, name string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.session = &UserSession{
		ID:   id,
		Name: name,
	}
	if err := sm.createSessionFile(); err != nil {
		return err
	}

	return nil
}

func (sm *SessionManager) GetSession() (*UserSession, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if err := sm.readSessionFile(); err != nil {
		return nil, err
	}

	return sm.session, nil
}

func (sm *SessionManager) DeleteSession() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.session = nil

	if err := sm.createSessionFile(); err != nil {
		return err
	}

	return nil
}

func (sm *SessionManager) createSessionFile() error {
	content, err := json.MarshalIndent(sm.session, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding content: %v", err)
	}

	if err := os.WriteFile(sm.dataPath, content, 0644); err != nil {
		return fmt.Errorf("writing session to file: %v", err)
	}

	return nil
}

func (sm *SessionManager) readSessionFile() error {
	content, err := os.ReadFile(sm.dataPath)
	if err != nil {
		return fmt.Errorf("reading session from file: %v", err)
	}

	if err := json.Unmarshal(content, &sm.session); err != nil {
		return fmt.Errorf("decoding content: %v", err)
	}

	return nil
}
