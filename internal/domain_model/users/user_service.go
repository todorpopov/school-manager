package users

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
)

type IUserService interface {
	CreateUser(ctx context.Context, tx pgx.Tx, createUser *CreateUser) (*User, *exceptions.AppError)
	GetUserById(ctx context.Context, tx pgx.Tx, userId int32) (*User, *exceptions.AppError)
	GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*User, *exceptions.AppError)
	GetUserByEmailWithPass(ctx context.Context, tx pgx.Tx, email string) (*User, *exceptions.AppError)
	GetUsers(ctx context.Context, tx pgx.Tx) ([]User, *exceptions.AppError)
	UpdateUser(ctx context.Context, tx pgx.Tx, updateUser *UpdateUser) (*User, *exceptions.AppError)
	UpdateUserPassword(ctx context.Context, tx pgx.Tx, updateUserPass *UpdateUserPassword) *exceptions.AppError
	DeleteUser(ctx context.Context, tx pgx.Tx, userId int32) *exceptions.AppError

	GetUsersRoles(ctx context.Context, tx pgx.Tx, userIds []int32) (map[int32][]string, *exceptions.AppError)
	AreRolesValid(ctx context.Context, roles []string) (bool, *exceptions.AppError)
}

type UserService struct {
	bcryptSvc internal.IBCryptService
	usrRepo   IUserRepository
	txFactory persistence.ITransactionFactory
}

func nullPassword(usr *User) {
	if usr != nil {
		usr.Password = nil
	}
}

func NewUserService(bcryptSvc internal.IBCryptService, usrRepo IUserRepository, txFactory persistence.ITransactionFactory) *UserService {
	return &UserService{bcryptSvc, usrRepo, txFactory}
}

func (us *UserService) CreateUser(ctx context.Context, tx pgx.Tx, createUser *CreateUser) (*User, *exceptions.AppError) {
	validationErr := ValidateCreateUser(createUser)
	if validationErr != nil {
		return nil, validationErr
	}

	hash, success := us.bcryptSvc.HashPassword(createUser.Password)
	if success != true {
		return nil, exceptions.NewAppError("INTERNAL_ERROR", "Failed to hash password", nil)
	}

	createUserCpy := *createUser
	createUserCpy.Password = hash

	// Use provided transaction or create new one
	var txToUse pgx.Tx
	var txErr *exceptions.AppError
	committed := false

	if tx != nil {
		txToUse = tx
	} else {
		txToUse, txErr = us.txFactory.BeginTransaction(ctx)
		if txErr != nil {
			return nil, txErr
		}
		defer func() {
			if !committed {
				_ = us.txFactory.CommitOrRollback(ctx, txToUse, txErr)
			}
		}()
	}

	usr, err := us.usrRepo.CreateUser(ctx, txToUse, &createUserCpy)
	if err != nil {
		txErr = err
		return nil, err
	}

	// Commit if we created the transaction
	if tx == nil {
		commitErr := us.txFactory.CommitOrRollback(ctx, txToUse, nil)
		if commitErr != nil {
			txErr = commitErr
			return nil, commitErr
		}
		committed = true
	}

	nullPassword(usr)
	return usr, nil
}

func (us *UserService) GetUserById(ctx context.Context, tx pgx.Tx, userId int32) (*User, *exceptions.AppError) {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}
	usr, err := us.usrRepo.GetUserById(ctx, tx, userId)
	nullPassword(usr)
	return usr, err
}

func (us *UserService) GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*User, *exceptions.AppError) {
	if msg := domain_model.ValidateEmail(&email, true); msg != "" {
		return nil, exceptions.NewValidationError("Invalid email", map[string]string{"email": msg})
	}
	usr, err := us.usrRepo.GetUserByEmail(ctx, tx, email)
	nullPassword(usr)
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
	if validationErr := ValidateUpdateUser(updateUser); validationErr != nil {
		return nil, validationErr
	}

	// Use provided transaction or create new one
	var txToUse pgx.Tx
	var txErr *exceptions.AppError
	committed := false

	if tx != nil {
		txToUse = tx
	} else {
		txToUse, txErr = us.txFactory.BeginTransaction(ctx)
		if txErr != nil {
			return nil, txErr
		}
		defer func() {
			if !committed {
				_ = us.txFactory.CommitOrRollback(ctx, txToUse, txErr)
			}
		}()
	}

	usr, err := us.usrRepo.UpdateUser(ctx, txToUse, updateUser)
	if err != nil {
		txErr = err
		return nil, err
	}

	// Commit if we created the transaction
	if tx == nil {
		commitErr := us.txFactory.CommitOrRollback(ctx, txToUse, nil)
		if commitErr != nil {
			txErr = commitErr
			return nil, commitErr
		}
		committed = true
	}

	nullPassword(usr)
	return usr, nil
}

func (us *UserService) UpdateUserPassword(ctx context.Context, tx pgx.Tx, updateUserPass *UpdateUserPassword) *exceptions.AppError {
	if validationErr := ValidateUpdateUserPassword(updateUserPass); validationErr != nil {
		return validationErr
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

func (us *UserService) GetUsersRoles(ctx context.Context, tx pgx.Tx, userIds []int32) (map[int32][]string, *exceptions.AppError) {
	if messages := domain_model.ValidateIds(userIds, true); len(messages) != 0 {
		return nil, exceptions.NewValidationError("Invalid IDs", messages)
	}
	return us.usrRepo.GetUsersRoles(ctx, tx, userIds)
}

func (us *UserService) AreRolesValid(ctx context.Context, roles []string) (bool, *exceptions.AppError) {
	if messages := domain_model.ValidateRoleNames(roles, true); len(messages) != 0 {
		return false, exceptions.NewValidationError("Invalid roles", messages)
	}
	return us.usrRepo.AreRolesValid(ctx, roles)
}
