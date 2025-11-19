package database

import (
	"fmt"
	"todo/internal/repositories/contracts"
	baseDB "todo/internal/shared/database"

	"github.com/google/uuid"
)

func New(appConfig contracts.AppConfig) contracts.Confirmation {
	return &confirmation{
		appConfig: appConfig,
		BaseDB:    baseDB.BaseDB{},
	}

}

type confirmation struct {
	baseDB.BaseDB
	appConfig contracts.AppConfig
}

func (c *confirmation) Delete(email string) error {
	sql := "DELETE FROM confirmations WHERE email = ?"
	return c.ExecuteSQL(c.appConfig, sql, email)
}

func (c *confirmation) FindByToken(token string) (string, error) {
	sql := "SELECT email FROM confirmations WHERE uuid = ?"
	res, err := c.RunSelectSQL(c.appConfig, sql, token)
	if err != nil {
		return "", err
	}

	if len(res) == 0 {
		return "", nil
	}

	if email, ok := res[0]["email"]; ok {
		if emailAsBytes, ok := email.([]byte); ok {
			return string(emailAsBytes), nil
		}

		return "", fmt.Errorf("cannot resolve token, email is not a byte slice from the database")
	}

	return "", fmt.Errorf("cannot resolve token, email could not retrieved from db")
}

func (c *confirmation) RegisterConfirmation(email string) (string, error) {
	token := uuid.New().String()
	sql := "INSERT INTO confirmations (email, uuid) VALUES (?,?)"
	err := c.ExecuteSQL(c.appConfig, sql, email, token)
	if err != nil {
		return "", err
	}

	return token, nil
}
