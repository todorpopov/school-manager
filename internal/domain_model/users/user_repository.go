package users

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, tx pgx.Tx, createUser *CreateUser) (*User, *exceptions.AppError)
	GetUserById(ctx context.Context, tx pgx.Tx, userId int32) (*User, *exceptions.AppError)
	GetUsers(ctx context.Context, tx pgx.Tx) ([]User, *exceptions.AppError)
	GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*User, *exceptions.AppError)
	UpdateUser(ctx context.Context, tx pgx.Tx, updateUser *UpdateUser) (*User, *exceptions.AppError)
	UpdateUserPassword(ctx context.Context, tx pgx.Tx, updateUserPass *UpdateUserPassword) *exceptions.AppError
	DeleteUser(ctx context.Context, tx pgx.Tx, userId int32) *exceptions.AppError
}

type UserRepository struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewUserRepository(db *persistence.Database, logger *zap.Logger) *UserRepository {
	return &UserRepository{db, logger}
}

func (ur *UserRepository) CreateUser(ctx context.Context, tx pgx.Tx, createUser *CreateUser) (*User, *exceptions.AppError) {
	validationErr := ValidateCreateUser(createUser)
	if validationErr != nil {
		return nil, validationErr
	}

	sql := "INSERT INTO users (first_name, last_name, email, password)" +
		"VALUES ($1, $2, $3, $4) " +
		"RETURNING user_id, first_name, last_name, email, password;"

	var user User
	var err error

	if tx != nil {
		ur.logger.Debug("Creating user in transaction")
		err = tx.QueryRow(ctx, sql, createUser.FirstName, createUser.LastName, createUser.Email, createUser.Password).
			Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	} else {
		ur.logger.Debug("Creating user without transaction")
		err = ur.db.Pool.QueryRow(ctx, sql, createUser.FirstName, createUser.LastName, createUser.Email, createUser.Password).
			Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	}

	if err != nil {
		ur.logger.Error("Failed to create user", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &user, nil
}

func (ur *UserRepository) GetUserById(ctx context.Context, tx pgx.Tx, userId int32) (*User, *exceptions.AppError) {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}

	sql := "SELECT user_id, first_name, last_name, email, password FROM users WHERE user_id = $1;"

	var user User
	var err error

	if tx != nil {
		ur.logger.Debug("Getting user by id transaction", zap.Int32("user_id", userId))
		err = tx.QueryRow(ctx, sql, userId).Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	} else {
		ur.logger.Debug("Getting user by id without transaction", zap.Int32("user_id", userId))
		err = ur.db.Pool.QueryRow(ctx, sql, userId).Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	}

	if err != nil {
		ur.logger.Error("Failed to get user by id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*User, *exceptions.AppError) {
	if msg := domain_model.ValidateEmail(&email, true); msg != "" {
		return nil, exceptions.NewValidationError("Invalid email", map[string]string{"email": msg})
	}

	sql := "SELECT user_id, first_name, last_name, email, password FROM users WHERE email = $1;"

	var user User
	var err error
	if tx != nil {
		ur.logger.Debug("Getting user by email transaction", zap.String("email", email))
		err = tx.QueryRow(ctx, sql, email).Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	} else {
		ur.logger.Debug("Getting user by email without transaction", zap.String("email", email))
		err = ur.db.Pool.QueryRow(ctx, sql, email).Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	}

	if err != nil {
		ur.logger.Error("Failed to get user by email", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	return &user, nil
}

func (ur *UserRepository) GetUsers(ctx context.Context, tx pgx.Tx) ([]User, *exceptions.AppError) {
	sql := "SELECT user_id, first_name, last_name, email, password FROM users;"

	var users []User
	var err error
	var rows pgx.Rows

	if tx != nil {
		rows, err = tx.Query(ctx, sql)
	} else {
		rows, err = ur.db.Pool.Query(ctx, sql)
	}

	if err != nil {
		ur.logger.Error("Failed to get users", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password)
		if err != nil {
			ur.logger.Error("Failed to scan user row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, tx pgx.Tx, updateUser *UpdateUser) (*User, *exceptions.AppError) {
	if validationErr := ValidateUpdateUser(updateUser); validationErr != nil {
		return nil, validationErr
	}

	sql := "UPDATE users " +
		" SET first_name = $1, last_name = $2, email = $3 " +
		"WHERE user_id = $4 " +
		"RETURNING user_id, first_name, last_name, email, password;"

	var user User
	var err error

	if tx != nil {
		ur.logger.Debug("Updating user in transaction")
		err = tx.QueryRow(ctx, sql, updateUser.FirstName, updateUser.LastName, updateUser.Email, updateUser.UserId).
			Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	} else {
		ur.logger.Debug("Updating user without transaction")
		err = ur.db.Pool.QueryRow(ctx, sql, updateUser.FirstName, updateUser.LastName, updateUser.Email, updateUser.UserId).
			Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	}

	if err != nil {
		ur.logger.Error("Failed to update user", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &user, nil
}

func (ur *UserRepository) UpdateUserPassword(ctx context.Context, tx pgx.Tx, updateUserPass *UpdateUserPassword) *exceptions.AppError {
	if err := ValidateUpdateUserPassword(updateUserPass); err != nil {
		return err
	}

	sql := "UPDATE users SET password = $1 WHERE user_id = $2;"

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		ur.logger.Debug("Updating user password in transaction")
		cmdTag, err = tx.Exec(ctx, sql, updateUserPass.Password, updateUserPass.UserId)
	} else {
		ur.logger.Debug("Updating user password without transaction")
		cmdTag, err = ur.db.Pool.Exec(ctx, sql, updateUserPass.Password, updateUserPass.UserId)
	}
	if err != nil {
		ur.logger.Error("Failed to update user password", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		ur.logger.Error("Failed to update user password - user not found", zap.Int32("user_id", updateUserPass.UserId))
		return exceptions.NewNotFoundError("User not found")
	}

	return nil
}

func (ur *UserRepository) DeleteUser(ctx context.Context, tx pgx.Tx, userId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}

	sql := "DELETE FROM users WHERE user_id = $1;"

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		ur.logger.Debug("Deleting user in transaction")
		cmdTag, err = tx.Exec(ctx, sql, userId)
	} else {
		ur.logger.Debug("Deleting user without transaction")
		cmdTag, err = ur.db.Pool.Exec(ctx, sql, userId)
	}
	if err != nil {
		ur.logger.Error("Failed to delete user", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		ur.logger.Error("Failed to delete user - user not found", zap.Int32("user_id", userId))
		return exceptions.NewNotFoundError("User not found")
	}

	return nil
}
