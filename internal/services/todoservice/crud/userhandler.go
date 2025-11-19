package crud

import (
	"html/template"
	"net/http"
	"todo/internal/repositories/user/results"
)

const (
	sessionErrorKey = "error"
)

type PageParams struct {
	Error string
}

func (s *service) registerHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/register.html")
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		s.registerPostHandler(w, r)
		return
	}

	sessionData, err := s.session.GetAll(s.sessionID)
	if err != nil {
		http.Error(w, "cannot get session "+err.Error(), http.StatusInternalServerError)
		return
	}

	errorMsg := sessionData[sessionErrorKey]
	s.session.Delete(s.sessionID)

	err = tmpl.Execute(w, PageParams{
		Error: errorMsg,
	})

	if err != nil {
		http.Error(w, "cannot execute template", http.StatusInternalServerError)
	}
}

func (s *service) registerPostHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if ok, msg := s.registerValidator(r); !ok {
		s.session.Set(s.sessionID, sessionErrorKey, msg)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	err = s.user.AddUser(&results.User{
		Name:  r.FormValue("name"),
		Email: r.FormValue("email"),
	}, r.FormValue("password"))
	if err != nil {
		http.Error(w, "cannot save user "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.registrationService.Register(r.FormValue("email"))
	if err != nil {
		http.Error(w, "cannot register user"+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}

func (s *service) loginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		s.loginPostHandler(w, r)
		return
	}

	sessionData, err := s.session.GetAll(s.sessionID)
	if err != nil {
		http.Error(w, "cannot get session "+err.Error(), http.StatusInternalServerError)
		return
	}

	errorMsg := sessionData[sessionErrorKey]
	s.session.Delete(s.sessionID)

	err = tmpl.Execute(w, PageParams{
		Error: errorMsg,
	})

	if err != nil {
		http.Error(w, "cannot execute template", http.StatusInternalServerError)
	}
}

func (s *service) loginPostHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	ok, err := s.loginValidator(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	if !ok {
		s.session.Set(s.sessionID, sessionErrorKey, "email or password incorrect")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return

	}

	s.session.Set(s.sessionID, "user", r.FormValue("email"))
	http.Redirect(w, r, "/todo", http.StatusSeeOther)

}

func (s *service) logoutHandler(w http.ResponseWriter, r *http.Request) {
	s.session.Delete(s.sessionID)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *service) homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}

func (s *service) registerValidator(r *http.Request) (bool, string) {
	if r.FormValue("name") == "" || r.FormValue("email") == "" || r.FormValue("password") == "" {
		return false, "All fields required"
	}

	if r.FormValue("password") != r.FormValue("repeat_password") {
		return false, "Password does not match"
	}

	return true, ""
}

func (s *service) loginValidator(r *http.Request) (bool, error) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		s.session.Set(s.sessionID, sessionErrorKey, "All fields required")
		return false, nil
	}

	ok, err := s.user.Authenticate(email, password)
	if err != nil {
		return false, err

	}

	return ok, nil
}

func (s *service) confirmHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}

	err := s.registrationService.Confirm(token)
	if err != nil {
		http.Error(w, "could not register "+err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
