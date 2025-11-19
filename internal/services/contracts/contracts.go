package contracts

type TodoService interface {
	Serve()
}

type RegistrationService interface {
	Register(email string) error
	Confirm(token string) error
}

type NotificationService interface {
	SendRegistrationNotification(email, token string) error
}
