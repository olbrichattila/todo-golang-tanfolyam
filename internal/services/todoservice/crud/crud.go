package crud

import (
	"net/http"
	repositoryContracts "todo/internal/repositories/contracts"
	"todo/internal/repositories/storage/result"
	"todo/internal/services/contracts"
)

func New(
	appConfig repositoryContracts.AppConfig,
	storage repositoryContracts.Storage,
	session repositoryContracts.Session,
	user repositoryContracts.User,
	registrationService contracts.RegistrationService,
) contracts.TodoService {
	// TODO error check
	return &service{
		appConfig:           appConfig,
		storage:             storage,
		session:             session,
		user:                user,
		registrationService: registrationService,
	}
}

type service struct {
	appConfig           repositoryContracts.AppConfig
	storage             repositoryContracts.Storage
	session             repositoryContracts.Session
	user                repositoryContracts.User
	registrationService contracts.RegistrationService
	sessionID           string
}

type PageData struct {
	Todos []result.Todo
}

// Serve implements contracts.TodoService.
func (s *service) Serve() {
	fs := http.FileServer(http.Dir("static"))
	http.HandleFunc("/", s.homeHandler)
	http.HandleFunc("/todo", s.sessionMiddleware(s.authMiddleware(s.todoHandler)))
	http.HandleFunc("/delete", s.sessionMiddleware(s.authMiddleware(s.deleteHandler)))

	http.HandleFunc("/register", s.sessionMiddleware(s.registerHandler))
	http.HandleFunc("/confirm", s.sessionMiddleware(s.confirmHandler))
	http.HandleFunc("/login", s.sessionMiddleware(s.loginHandler))
	http.HandleFunc("/logout", s.sessionMiddleware(s.logoutHandler))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(s.appConfig.AppPort(), nil)
}
