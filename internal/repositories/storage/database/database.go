package database

import (
	"todo/internal/repositories/contracts"
	"todo/internal/repositories/storage/result"
	baseDB "todo/internal/shared/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func New(appConfig contracts.AppConfig) contracts.Storage {
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

func (m *dbStorage) Delete(uuid string) error {
	sql := "DELETE FROM todos WHERE uuid= ?"

	return m.ExecuteSQL(m.appConfig, sql, uuid)
}

func (m *dbStorage) AddTodo(msg string) (err error) {
	sql := "INSERT INTO todos (todo, uuid) values (?,?)"
	newUUID := uuid.New().String()
	return m.ExecuteSQL(m.appConfig, sql, msg, newUUID)
}

func (m *dbStorage) ListTodos() ([]result.Todo, error) {
	sql := "SELECT todo, uuid FROM todos"
	results, err := m.RunSelectSQL(m.appConfig, sql)
	if err != nil {
		return nil, err
	}

	var resultList = make([]result.Todo, len(results))
	for i, row := range results {
		resultList[i] = result.Todo{
			Todo: string(row["todo"].([]byte)),
			UUID: string(row["uuid"].([]byte)),
		}
	}

	return resultList, nil
}
