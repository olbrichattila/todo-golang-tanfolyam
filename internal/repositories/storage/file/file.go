package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"todo/internal/repositories/contracts"
	"todo/internal/repositories/storage/result"

	"github.com/google/uuid"
)

const fileName = "./data/data.txt"

func New() contracts.Storage {
	return &fileStorage{}
}

type fileStorage struct {
}

func (m *fileStorage) Delete(uuid string) error {
	todos, err := m.ListTodos()
	if err != nil {
		return err
	}
	err = os.Remove(fileName)
	if err != nil {
		return err
	}

	for _, todoItem := range todos {
		if todoItem.UUID != uuid {
			m.addTodoWithUUID(todoItem.Todo, todoItem.UUID)
		}
	}

	return nil
}

func (m *fileStorage) AddTodo(msg string) error {
	return m.addTodoWithUUID(msg, uuid.New().String())
}

func (m *fileStorage) addTodoWithUUID(msg, todoUUID string) error {
	fileHandler, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("%w, cannot open file %s", err, fileName)
	}

	defer fileHandler.Close()
	row, err := json.Marshal(result.Todo{
		Todo: msg,
		UUID: todoUUID,
	})

	_, err = fileHandler.WriteString(string(row) + "\n")
	if err != nil {
		return fmt.Errorf("%w, cannot write to the file %s", err, fileName)
	}

	return nil
}

func (m *fileStorage) ListTodos() ([]result.Todo, error) {
	fileHandler, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("%w, cannot open file %s", err, fileName)
	}

	defer fileHandler.Close()

	scanner := bufio.NewScanner(fileHandler)
	resultList := []result.Todo{}

	for scanner.Scan() {
		line := scanner.Text()
		var row result.Todo
		err = json.Unmarshal([]byte(line), &row)
		if err != nil {
			return nil, fmt.Errorf("%w, retrieve row %s", err, fileName)
		}

		resultList = append(resultList, row)
	}

	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("%w, error during read %s", err, fileName)
	}

	return resultList, nil
}
