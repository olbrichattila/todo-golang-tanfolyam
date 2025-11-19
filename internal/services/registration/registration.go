package registration

import (
	"fmt"
	repositoryContracts "todo/internal/repositories/contracts"
	"todo/internal/services/contracts"
)

func New(
	userRepository repositoryContracts.User,
	confirmationRepository repositoryContracts.Confirmation,
	notificationService contracts.NotificationService,
) contracts.RegistrationService {
	// TODO error handling of nils
	return &regService{
		userRepository:         userRepository,
		confirmationRepository: confirmationRepository,
		notificationService:    notificationService,
	}
}

type regService struct {
	userRepository         repositoryContracts.User
	confirmationRepository repositoryContracts.Confirmation
	notificationService    contracts.NotificationService
}

func (r *regService) Confirm(token string) error {
	email, err := r.confirmationRepository.FindByToken(token)
	if err != nil {
		return err
	}

	if email == "" {
		return fmt.Errorf("Cannot find confirmation information")
	}

	err = r.userRepository.Activate(email)
	if err != nil {
		return err
	}

	return r.confirmationRepository.Delete(email)
}

func (r *regService) Register(email string) error {
	token, err := r.confirmationRepository.RegisterConfirmation(email)
	if err != nil {
		return err
	}

	return r.notificationService.SendRegistrationNotification(email, token)
}
