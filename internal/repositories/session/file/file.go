package file

import (
	"encoding/json"
	"fmt"
	"os"
	"todo/internal/repositories/contracts"
)

const (
	sessionFolderName = "./session/"
)

func New() contracts.Session {
	return &session{}
}

type session struct {
}

// Delete implements contracts.Session.
func (s *session) Delete(sessionID string) error {
	fileName := s.getFileName(sessionID)
	err := os.Remove(fileName)
	if err != nil && os.IsNotExist(err) {
		return nil
	}

	return err
}

// Get implements contracts.Session.
func (s *session) Get(sessionID string, key string) (string, error) {
	sessionData, err := s.GetAll(sessionID)
	if err != nil {
		return "", err
	}

	result, ok := sessionData[key]
	if !ok {
		return "", fmt.Errorf("session variable cannot found")
	}

	return result, nil
}

// GetAll implements contracts.Session.
func (s *session) GetAll(sessionID string) (map[string]string, error) {
	result := map[string]string{}
	fileName := s.getFileName(sessionID)
	content, err := os.ReadFile(fileName)
	if err != nil && os.IsNotExist(err) {
		return result, nil
	}

	if err != nil {
		return nil, err
	}

	return s.fromSessionData(string(content))
}

// Set implements contracts.Session.
func (s *session) Set(sessionID string, key string, value string) error {
	sessionData, err := s.GetAll(sessionID)
	if err != nil {
		return err
	}

	sessionData[key] = value

	sessionAsString, err := s.toSessionData(sessionData)
	if err != nil {
		return err
	}

	fileName := s.getFileName(sessionID)

	return os.WriteFile(fileName, []byte(sessionAsString), 0644)
}

func (s *session) fromSessionData(data string) (map[string]string, error) {
	var result map[string]string
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *session) toSessionData(data map[string]string) (string, error) {
	result, err := json.Marshal(data)
	if err != nil {
		return "", nil
	}

	return string(result), nil
}

func (s *session) getFileName(sessionID string) string {
	return sessionFolderName + sessionID + ".json"
}
