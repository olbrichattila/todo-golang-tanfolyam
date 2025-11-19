package crud

import (
	"html/template"
	"net/http"
)

func (s *service) todoHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		todoToAdd := r.FormValue("todo")
		s.storage.AddTodo(todoToAdd)

		http.Redirect(w, r, "/todo", http.StatusSeeOther)

		return
	}

	todoList, err := s.storage.ListTodos()
	if err != nil {
		http.Error(w, "cannot load todos", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, PageData{
		Todos: todoList,
	})
	if err != nil {
		http.Error(w, "template execution error", http.StatusInternalServerError)
	}
}

func (s *service) deleteHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		http.Error(w, "uuid parameter is required", http.StatusBadRequest)
		return
	}

	err := s.storage.Delete(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}
