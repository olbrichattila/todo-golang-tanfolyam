package app

import (
	"fmt"

	"todo/internal/repositories/appconfing/env"
	"todo/internal/repositories/contracts"
	fileSession "todo/internal/repositories/session/file"
	"todo/internal/repositories/storage/database"
	"todo/internal/repositories/storage/file"
	"todo/internal/repositories/storage/memory"
	userDBRepository "todo/internal/repositories/user/database"
	serviceContracts "todo/internal/services/contracts"
	registrationService "todo/internal/services/registration"
	"todo/internal/services/todoservice/api"
	"todo/internal/services/todoservice/cmd"
	"todo/internal/services/todoservice/crud"

	databaseConfirmationRepository "todo/internal/repositories/confirmation/database"
	emailNotificationService "todo/internal/services/notification/email"
)

const environmentVariableName = "STORAGE_TYPE"

type AppFactory interface {
	Storage() contracts.Storage
	Serve()
}

type app struct {
	storage        contracts.Storage
	appConfig      contracts.AppConfig
	session        contracts.Session
	userRepository contracts.User
	todoService    serviceContracts.TodoService

	confirmationRepository contracts.Confirmation
	notificationService    serviceContracts.NotificationService
	registrationService    serviceContracts.RegistrationService
}

func New() AppFactory {
	app := &app{}
	app.initSession()
	app.initAppConfig()
	app.initStorage()
	app.initUserRepository()

	app.initConfirmationRepository()
	app.initNotificationService()
	app.initRegistrationService()

	app.initTodoService()

	return app
}

func (a *app) initRegistrationService() {
	a.registrationService = registrationService.New(
		a.userRepository,
		a.confirmationRepository,
		a.notificationService,
	)
}

func (a *app) initConfirmationRepository() {
	a.confirmationRepository = databaseConfirmationRepository.New(a.appConfig)
}

func (a *app) initNotificationService() {
	a.notificationService = emailNotificationService.New(a.appConfig)
}

func (a *app) initUserRepository() {
	a.userRepository = userDBRepository.New(a.appConfig)
}

func (a *app) initSession() {
	a.session = fileSession.New()
}

func (a *app) initAppConfig() {
	a.appConfig = env.New()
}

func (a *app) initTodoService() {
	switch a.appConfig.FrontendType() {
	case "cmd":
		a.todoService = cmd.New(a.storage)
	case "api":
		a.todoService = api.New(a.appConfig, a.storage)
	case "crud":
		a.todoService = crud.New(
			a.appConfig,
			a.storage,
			a.session,
			a.userRepository,
			a.registrationService,
		)
	default:
		panic(fmt.Sprintf("missing environment variable for frontend, %s, can be cmd  or api", "FRONTEND"))
	}
}

func (a *app) initStorage() {
	switch a.appConfig.StorageType() {
	case "file":
		a.storage = file.New()
	case "memory":
		a.storage = memory.New()
	case "db":
		a.storage = database.New(a.appConfig)
	default:
		panic(fmt.Sprintf("missing environment variable %s, can be file or memory", environmentVariableName))
	}
}

func (a *app) Storage() contracts.Storage {
	return a.storage
}

func (a *app) Serve() {
	a.todoService.Serve()
}
