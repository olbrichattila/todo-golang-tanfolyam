package helpers

import "fmt"

type SelectedMenu int

const (
	AddTodo SelectedMenu = iota
	ListTodos
	DeleteTodo
	NotSelected
)

func DisplayMenu() (SelectedMenu, error) {
	for {
		fmt.Println("1. Add todo")
		fmt.Println("2. List todos")
		fmt.Println("3. Delete todo")
		fmt.Println("-------------")

		input, err := GetLine()
		if err != nil {
			return NotSelected, fmt.Errorf("menu not selected %w", err)
		}

		if input == "1" {
			return AddTodo, nil
		}

		if input == "2" {
			return ListTodos, nil
		}

		if input == "3" {
			return DeleteTodo, nil
		}

		fmt.Println("Incorrect answer type 1 or 2")
	}
}
