package users

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type IUserService interface {
	CreateUser(ctx context.Context, tx pgx.Tx, createUser *CreateUser) (*User, *exceptions.AppError)
	GetUserById(ctx context.Context, tx pgx.Tx, userId int32) (*User, *exceptions.AppError)
	GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*User, *exceptions.AppError)
	GetUserByEmailWithPass(ctx context.Context, tx pgx.Tx, email string) (*User, *exceptions.AppError)
	GetUsers(ctx context.Context, tx pgx.Tx) (*[]User, *exceptions.AppError)
	UpdateUser(ctx context.Context, tx pgx.Tx, updateUser *UpdateUser) (*User, *exceptions.AppError)
	UpdateUserPassword(ctx context.Context, tx pgx.Tx, updateUserPass *UpdateUserPassword) *exceptions.AppError
	DeleteUser(ctx context.Context, tx pgx.Tx, userId int32) *exceptions.AppError
}

type UserService struct {
	bcryptSvc internal.IBCryptService
	usrRepo   IUserRepository
}

func NewUserService(bcryptSvc internal.IBCryptService, usrRepo IUserRepository) *UserService {
	return &UserService{bcryptSvc, usrRepo}
}

func (us *UserService) CreateUser(ctx context.Context, tx pgx.Tx, createUser *CreateUser) (*User, *exceptions.AppError) {
	err := ValidateCreateUser(createUser)
	if err != nil {
		return nil, err
	}

	hash, success := us.bcryptSvc.HashPassword(createUser.Password)
	if success != true {
		return nil, exceptions.NewAppError("INTERNAL_ERROR", "Failed to hash password", nil)
	}

	createUser.Password = hash
	usr, err := us.usrRepo.CreateUser(ctx, tx, createUser)
	if usr != nil {
		usr.Password = nil
	}
	return usr, err
}

func (us *UserService) GetUserById(ctx context.Context, tx pgx.Tx, userId int32) (*User, *exceptions.AppError) {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}
	usr, err := us.usrRepo.GetUserById(ctx, tx, userId)
	if usr != nil {
		usr.Password = nil
	}
	return usr, err
}

func (us *UserService) GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*User, *exceptions.AppError) {
	if msg := domain_model.ValidateEmail(&email, true); msg != "" {
		return nil, exceptions.NewValidationError("Invalid email", map[string]string{"email": msg})
	}
	usr, err := us.usrRepo.GetUserByEmail(ctx, tx, email)
	if usr != nil {
		usr.Password = nil
	}
	return usr, err
}

func (us *UserService) GetUserByEmailWithPass(ctx context.Context, tx pgx.Tx, email string) (*User, *exceptions.AppError) {
	if msg := domain_model.ValidateEmail(&email, true); msg != "" {
		return nil, exceptions.NewValidationError("Invalid email", map[string]string{"email": msg})
	}
	return us.usrRepo.GetUserByEmail(ctx, tx, email)
}

func (us *UserService) GetUsers(ctx context.Context, tx pgx.Tx) ([]User, *exceptions.AppError) {
	users, err := us.usrRepo.GetUsers(ctx, tx)
	for i := range users {
		users[i].Password = nil
	}
	return users, err
}

func (us *UserService) UpdateUser(ctx context.Context, tx pgx.Tx, updateUser *UpdateUser) (*User, *exceptions.AppError) {
	if err := ValidateUpdateUser(updateUser); err != nil {
		return nil, err
	}
	usr, err := us.usrRepo.UpdateUser(ctx, tx, updateUser)
	if usr != nil {
		usr.Password = nil
	}
	return usr, err
}

func (us *UserService) UpdateUserPassword(ctx context.Context, tx pgx.Tx, updateUserPass *UpdateUserPassword) *exceptions.AppError {
	if err := ValidateUpdateUserPassword(updateUserPass); err != nil {
		return err
	}
	hash, success := us.bcryptSvc.HashPassword(updateUserPass.Password)
	if success != true {
		return exceptions.NewAppError("INTERNAL_ERROR", "Failed to hash password", nil)
	}

	updateUserPass.Password = hash
	return us.usrRepo.UpdateUserPassword(ctx, tx, updateUserPass)
}

func (us *UserService) DeleteUser(ctx context.Context, tx pgx.Tx, userId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}
	return us.usrRepo.DeleteUser(ctx, tx, userId)
}
