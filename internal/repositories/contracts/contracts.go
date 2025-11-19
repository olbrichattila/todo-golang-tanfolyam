package contracts

import (
	"todo/internal/repositories/storage/result"
	"todo/internal/repositories/user/results"
)

type Storage interface {
	AddTodo(string) error
	ListTodos() ([]result.Todo, error)
	Delete(uuid string) error
}

type AppConfig interface {
	StorageType() string
	FrontendType() string
	DBHost() string
	DBPort() string
	DBDatabase() string
	DBUsername() string
	DBPassword() string

	SmtpHost() string
	SmtpPort() string
	SmtpUsername() string
	SmtpPassword() string
	EmailFrom() string
	AppURL() string
	AppPort() string
}

type Session interface {
	GetAll(sessionID string) (map[string]string, error)
	Get(sessionID, key string) (string, error)
	Set(sessionID, key, value string) error
	Delete(sessionID string) error
}

type User interface {
	Activate(email string) error
	UserByEmail(email string) (*results.User, error)
	Authenticate(email, password string) (bool, error)
	AddUser(user *results.User, password string) error
}

type Confirmation interface {
	RegisterConfirmation(email string) (string, error)
	FindByToken(token string) (string, error)
	Delete(email string) error
}
