package api

import (
	"encoding/json"
	"io"
	"net/http"
	repositoryContracts "todo/internal/repositories/contracts"
	"todo/internal/services/contracts"
)

type todo struct {
	Todo string `json:"todo"`
}

func New(
	appConfig repositoryContracts.AppConfig,
	storage repositoryContracts.Storage,
) contracts.TodoService {
	return &service{
		appConfig: appConfig,
		storage:   storage,
	}
}

type service struct {
	appConfig repositoryContracts.AppConfig
	storage   repositoryContracts.Storage
}

// Serve implements contracts.TodoService.
func (s *service) Serve() {
	http.HandleFunc("/add", s.jsonMiddleware(s.addHandler))
	http.HandleFunc("/list", s.jsonMiddleware(s.listHandler))
	http.HandleFunc("/delete", s.jsonMiddleware(s.deleteHandler))

	http.ListenAndServe(s.appConfig.AppPort(), nil)
}

func (s *service) addHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusInternalServerError)
		return
	}

	var todoData todo
	err = json.Unmarshal(body, &todoData)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	err = s.storage.AddTodo(todoData.Todo)
	if err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}
}

func (s *service) listHandler(w http.ResponseWriter, r *http.Request) {
	todoList, err := s.storage.ListTodos()
	if err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(todoList)
	if err != nil {
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
		return
	}
}

func (s *service) deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "not a delete method", http.StatusBadRequest)
		return
	}

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		http.Error(w, "uuid parameter is required", http.StatusBadRequest)
		return
	}

	err := s.storage.Delete(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *service) jsonMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		next(w, r)
	}
}
