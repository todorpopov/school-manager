package users

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"go.uber.org/zap"
)

type IUserService interface {
	validateCreateUser(createUser *CreateUser) *exceptions.AppError

	CreateUser(ctx context.Context, tx pgx.Tx, createUser *CreateUser) (*User, *exceptions.AppError)
}

type UserService struct {
	usrRepo IUserRepository
	logger  *zap.Logger
}

func NewUserService(usrRepo IUserRepository, logger *zap.Logger) *UserService {
	return &UserService{usrRepo, logger}
}

func (us *UserService) validateCreateUser(createUser *CreateUser) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateString(&createUser.FirstName, 1, 255, true)
	if msg != "" {
		messages["first_name"] = msg
	}

	msg = domain_model.ValidateString(&createUser.LastName, 1, 255, true)
	if msg != "" {
		messages["last_name"] = msg
	}

	msg = domain_model.ValidateEmail(&createUser.Email, true)
	if msg != "" {
		messages["email"] = msg
	}

	msg = domain_model.ValidatePassword(&createUser.Password, true)
	if msg != "" {
		messages["password"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during user creation", messages)
	}
	return nil
}

func (us *UserService) CreateUser(ctx context.Context, tx pgx.Tx, createUser *CreateUser) (*User, *exceptions.AppError) {
	err := us.validateCreateUser(createUser)
	if err != nil {
		return nil, err
	}

	usr, err := us.usrRepo.CreateUser(ctx, tx, createUser)
	if err == nil {
		usr.Password = nil
	}

	return usr, err
}
