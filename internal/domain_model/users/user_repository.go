package users

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"github.com/ydb-platform/ydb-go-sdk/v3/config"
	"go.uber.org/zap"
)

type UserRepo interface {
	CreateUser(ctx context.Context, tx pgx.Tx, createUser *CreateUser) (*User, *exceptions.AppError)
	GetUserById(ctx context.Context, tx pgx.Tx, userId int32) (*User, *exceptions.AppError)
	GetUsers(ctx context.Context, tx pgx.Tx) ([]User, *exceptions.AppError)
	UpdateUser(ctx context.Context, tx pgx.Tx, user *User) (*User, *exceptions.AppError)
	UpdateUserPassword(ctx context.Context, tx pgx.Tx, userId int32, password string) *exceptions.AppError
	DeleteUser(ctx context.Context, tx pgx.Tx, userId int32) *exceptions.AppError
}

type UserRepository struct {
	db     *persistence.Database
	env    *config.Config
	logger *zap.Logger
}

func NewUserRepository(db *persistence.Database, env *config.Config, logger *zap.Logger) *UserRepository {
	return &UserRepository{db, env, logger}
}

func (ur *UserRepository) CreateUser(ctx context.Context, tx pgx.Tx, createUser *CreateUser) (*User, *exceptions.AppError) {
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
		return nil, domain_model.PgErrorToAppError(err)
	}

	return &user, nil
}

func (ur *UserRepository) GetUserById(ctx context.Context, tx pgx.Tx, userId int32) (*User, *exceptions.AppError) {
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
		return nil, domain_model.PgErrorToAppError(err)
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
		return nil, domain_model.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password)
		if err != nil {
			ur.logger.Error("Failed to scan user row", zap.Error(err))
			return nil, domain_model.PgErrorToAppError(err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, tx pgx.Tx, updateUser *UpdateUser) (*User, *exceptions.AppError) {
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
		return nil, domain_model.PgErrorToAppError(err)
	}

	return &user, nil
}

func (ur *UserRepository) UpdateUserPassword(ctx context.Context, tx pgx.Tx, userId int32, password string) *exceptions.AppError {
	sql := "UPDATE users SET password = $1 WHERE user_id = $2;"

	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, sql, password, userId)
	} else {
		_, err = ur.db.Pool.Exec(ctx, sql, password, userId)
	}
	return domain_model.PgErrorToAppError(err)
}

func (ur *UserRepository) DeleteUser(ctx context.Context, tx pgx.Tx, userId int32) *exceptions.AppError {
	sql := "DELETE FROM users WHERE user_id = $1;"

	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, sql, userId)
	} else {
		_, err = ur.db.Pool.Exec(ctx, sql, userId)
	}
	return domain_model.PgErrorToAppError(err)
}
