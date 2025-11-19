package cmd

import (
	"fmt"
	"todo/internal/helpers"
	repositoryContracts "todo/internal/repositories/contracts"
	"todo/internal/repositories/storage/result"
	"todo/internal/services/contracts"
)

func New(storage repositoryContracts.Storage) contracts.TodoService {
	return &service{
		storage: storage,
	}
}

type service struct {
	storage repositoryContracts.Storage
}

// Serve implements contracts.TodoService.
func (s *service) Serve() {
	for {
		selectedMenu, err := helpers.DisplayMenu()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if selectedMenu == helpers.AddTodo {
			fmt.Println("Please type: ")
			input, err := helpers.GetLine()
			if err != nil {
				fmt.Println(err.Error())
				break
			}

			err = s.storage.AddTodo(input)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}

		if selectedMenu == helpers.ListTodos {
			todos, err := s.storage.ListTodos()
			if err != nil {
				fmt.Println(err.Error())
				break
			}

			s.displayTodos(todos)
		}

		if selectedMenu == helpers.DeleteTodo {
			fmt.Print("UUID of the todo: ")
			delUUID, err := helpers.GetLine()
			if err != nil {
				fmt.Println(err.Error())
				break
			}

			err = s.storage.Delete(delUUID)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func (s *service) displayTodos(todos []result.Todo) {
	fmt.Println("Your todos:")
	for i, todo := range todos {
		fmt.Printf("%d: %s | %s\n", i+1, todo.Todo, todo.UUID)
	}
	fmt.Println("---------")
}
