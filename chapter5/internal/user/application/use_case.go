package application

import (
	"context"
	"errors"

	"github.com/PacktPublishing/Domain-Driven-Design-with-Go/chapter5/internal/user/domain"
)

var (
	// RegistrationUseCaseInternalError defines unexpected error during registration
	RegistrationUseCaseInternalError = errors.New("registration.internalError")
	// RegistrationUseCaseUserAlreadyCreated defines an error when such user is already created
	RegistrationUseCaseUserAlreadyCreated = errors.New("registration.userAlreadyCreated")
)

// RegistrationUseCase defines a contract for registration use case
type RegistrationUseCase interface {
	Execute(ctx context.Context, user domain.User, password string) (*domain.User, error)
}

// defaultRegistrationUseCase is an actual implementation of RegistrationUseCase
type defaultRegistrationUseCase struct {
	repository domain.UserRepository
}

// NewRegistrationUseCase creates new defaultRegistrationUseCase
func NewRegistrationUseCase(repository domain.UserRepository) RegistrationUseCase {
	return &defaultRegistrationUseCase{
		repository: repository,
	}
}

// Execute performs complete registration and return the result
func (c *defaultRegistrationUseCase) Execute(ctx context.Context, user domain.User, password string) (*domain.User, error) {
	users, err := c.repository.SearchByUsername(ctx, user.Username)
	if err != nil {
		return nil, RegistrationUseCaseInternalError
	} else if len(users) > 0 {
		return nil, RegistrationUseCaseUserAlreadyCreated
	}

	newUser, err := c.repository.Create(ctx, user, password)
	if err != nil {
		return nil, RegistrationUseCaseInternalError
	}

	return newUser, nil
}
