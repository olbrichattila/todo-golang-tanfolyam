package memory

import (
	"fmt"
	"sync"
	"todo/internal/repositories/contracts"
	"todo/internal/repositories/storage/result"

	"github.com/google/uuid"
)

func New() contracts.Storage {
	return &memoryStorage{}
}

type memoryStorage struct {
	mu    sync.RWMutex
	todos []result.Todo
}

func (m *memoryStorage) Delete(uuid string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	idx, err := m.getTodoIndexByUUID(uuid)
	if err != nil {
		return err
	}

	m.todos = append(m.todos[:idx], m.todos[idx+1:]...)
	return nil
}

func (m *memoryStorage) AddTodo(msg string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	newTodo := result.Todo{
		UUID: uuid.New().String(),
		Todo: msg,
	}
	m.todos = append(m.todos, newTodo)

	return nil
}

func (m *memoryStorage) ListTodos() ([]result.Todo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.todos, nil
}

func (m *memoryStorage) getTodoIndexByUUID(uuid string) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for idx, todo := range m.todos {
		if todo.UUID == uuid {
			return idx, nil
		}
	}

	return 0, fmt.Errorf("Todo does not exists")
}
