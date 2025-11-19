package database

import (
	"fmt"
	"todo/internal/repositories/contracts"
	"todo/internal/repositories/user/results"
	baseDB "todo/internal/shared/database"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func New(appConfig contracts.AppConfig) contracts.User {
	// TODO if appConfig == nil { return ... }
	return &dbStorage{
		appConfig: appConfig,
		BaseDB:    baseDB.BaseDB{},
	}
}

type dbStorage struct {
	baseDB.BaseDB
	appConfig contracts.AppConfig
}

func (m *dbStorage) Activate(email string) error {
	sql := "UPDATE users SET is_active = 1 WHERE email = ?"
	return m.ExecuteSQL(m.appConfig, sql, email)
}

// AddUser implements contracts.User.
func (m *dbStorage) AddUser(user *results.User, password string) error {
	if user == nil {
		return fmt.Errorf("the AddUser requires user, nil passed")
	}

	sql := "INSERT INTO users (email, name, password) values (?,?,?)"

	hashedPassword, err := m.hashPassword(password)
	if err != nil {
		return err
	}

	return m.ExecuteSQL(m.appConfig, sql, user.Email, user.Name, hashedPassword)
}

// Authenticate implements contracts.User.
func (m *dbStorage) Authenticate(email string, password string) (bool, error) {
	sql := "SELECT password FROM users where email = ? and is_active = 1"

	result, err := m.RunSelectSQL(m.appConfig, sql, email)
	if err != nil {
		return false, err
	}

	if len(result) == 0 {
		return false, nil
	}

	passwordOK := m.checkPasswordHash(password, string(result[0]["password"].([]byte)))

	return passwordOK, nil
}

// UserByEmail implements contracts.User.
func (m *dbStorage) UserByEmail(email string) (*results.User, error) {
	sql := "SELECT email, name FROM users where email = ?"

	result, err := m.RunSelectSQL(m.appConfig, sql, email)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &results.User{
		Email: string(result[0]["email"].([]byte)),
		Name:  string(result[0]["name"].([]byte)),
	}, nil
}

func (m *dbStorage) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (m *dbStorage) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
