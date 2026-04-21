package users

import (
	"context"
	"slices"

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

	GetUsersRoles(ctx context.Context, tx pgx.Tx, userIds []int32) (map[int32][]string, *exceptions.AppError)
	AreRolesValid(ctx context.Context, roles []string) (bool, *exceptions.AppError)
	deleteUserRoles(ctx context.Context, tx pgx.Tx, userId int32) *exceptions.AppError
}

type UserRepository struct {
	db        *persistence.Database
	logger    *zap.Logger
	txFactory persistence.ITransactionFactory
}

func NewUserRepository(db *persistence.Database, logger *zap.Logger) *UserRepository {
	txFactory := persistence.NewTransactionFactory(db)
	return &UserRepository{db, logger, txFactory}
}

func (ur *UserRepository) CreateUser(ctx context.Context, tx pgx.Tx, createUser *CreateUser) (*User, *exceptions.AppError) {
	validationErr := ValidateCreateUser(createUser)
	if validationErr != nil {
		return nil, validationErr
	}

	if _, rolesErr2 := ur.AreRolesValid(ctx, createUser.Roles); rolesErr2 != nil {
		ur.logger.Error("Failed to validate roles", zap.Error(rolesErr2))
		return nil, rolesErr2
	}

	var txToUse pgx.Tx
	var txErr *exceptions.AppError
	committed := false

	if tx != nil {
		txToUse = tx
	} else {
		txToUse, txErr = ur.txFactory.BeginTransaction(ctx)
		if txErr != nil {
			return nil, txErr
		}
		defer func() {
			if !committed {
				_ = ur.txFactory.CommitOrRollback(ctx, txToUse, txErr)
			}
		}()
	}

	sql := `
		WITH new_user AS (
			INSERT INTO users (first_name, last_name, email, password)
			VALUES ($1, $2, $3, $4)
			RETURNING user_id, first_name, last_name, email, password
		)
		INSERT INTO user_roles (user_id, role_id)
		SELECT new_user.user_id, r.role_id
		FROM new_user
		CROSS JOIN unnest($5::text[]) AS rn(role_name)
		INNER JOIN roles r ON r.role_name = rn.role_name
		RETURNING (SELECT user_id FROM new_user), 
				  (SELECT first_name FROM new_user), 
				  (SELECT last_name FROM new_user), 
				  (SELECT email FROM new_user), 
				  (SELECT password FROM new_user);
		`

	var user User
	var err error

	ur.logger.Debug("Creating user", zap.Bool("has_transaction", tx != nil))
	err = txToUse.QueryRow(ctx, sql, createUser.FirstName, createUser.LastName, createUser.Email, createUser.Password, createUser.Roles).
		Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password)

	if err != nil {
		ur.logger.Error("Failed to create user", zap.Error(err))
		txErr = exceptions.PgErrorToAppError(err)
		return nil, txErr
	}

	if tx == nil {
		commitErr := ur.txFactory.CommitOrRollback(ctx, txToUse, nil)
		if commitErr != nil {
			txErr = commitErr
			return nil, commitErr
		}
		committed = true
	}

	user.Roles = createUser.Roles
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

	rolesMap, rolesErr := ur.GetUsersRoles(ctx, tx, []int32{userId})
	if rolesErr != nil {
		ur.logger.Error("Failed to get user roles", zap.Error(rolesErr))
		return nil, rolesErr
	}
	user.Roles = rolesMap[user.UserId]

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

	rolesMap, rolesErr := ur.GetUsersRoles(ctx, tx, []int32{user.UserId})
	if rolesErr != nil {
		ur.logger.Error("Failed to get user roles", zap.Error(rolesErr))
		return nil, rolesErr
	}
	user.Roles = rolesMap[user.UserId]

	return &user, nil
}

func (ur *UserRepository) GetUsers(ctx context.Context, tx pgx.Tx) ([]User, *exceptions.AppError) {
	sql := "SELECT user_id, first_name, last_name, email, password FROM users;"

	var users []User
	var err error
	var rows pgx.Rows

	if tx != nil {
		ur.logger.Debug("Getting users in transaction")
		rows, err = tx.Query(ctx, sql)
	} else {
		ur.logger.Debug("Getting users without transaction")
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

	if err = rows.Err(); err != nil {
		ur.logger.Error("Error iterating users rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	if len(users) > 0 {
		userIds := make([]int32, len(users))
		for i := range users {
			userIds[i] = users[i].UserId
		}

		rolesMap, rolesErr := ur.GetUsersRoles(ctx, tx, userIds)
		if rolesErr != nil {
			ur.logger.Error("Failed to get user roles", zap.Error(rolesErr))
			return nil, rolesErr
		}

		for i := range users {
			users[i].Roles = rolesMap[users[i].UserId]
		}
	}

	return users, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, tx pgx.Tx, updateUser *UpdateUser) (*User, *exceptions.AppError) {
	if validationErr := ValidateUpdateUser(updateUser); validationErr != nil {
		return nil, validationErr
	}

	var txToUse pgx.Tx
	var txErr *exceptions.AppError
	committed := false

	if tx != nil {
		txToUse = tx
	} else {
		txToUse, txErr = ur.txFactory.BeginTransaction(ctx)
		if txErr != nil {
			return nil, txErr
		}
		defer func() {
			if !committed {
				_ = ur.txFactory.CommitOrRollback(ctx, txToUse, txErr)
			}
		}()
	}

	rolesMap, rolesErr := ur.GetUsersRoles(ctx, txToUse, []int32{updateUser.UserId})
	if rolesErr != nil {
		ur.logger.Error("Failed to get user roles", zap.Error(rolesErr))
		txErr = rolesErr
		return nil, rolesErr
	}

	existingUserRoles := rolesMap[updateUser.UserId]
	slices.Sort(existingUserRoles)
	slices.Sort(updateUser.Roles)

	var err error
	var user User

	if slices.Compare(existingUserRoles, updateUser.Roles) != 0 {
		ur.logger.Debug("Updating roles", zap.Strings("existing_roles", existingUserRoles), zap.Strings("updated_roles", updateUser.Roles))

		if _, rolesErr2 := ur.AreRolesValid(ctx, updateUser.Roles); rolesErr2 != nil {
			ur.logger.Error("Failed to validate roles", zap.Error(rolesErr2))
			txErr = rolesErr2
			return nil, rolesErr2
		}

		if delRolesErr := ur.deleteUserRoles(ctx, txToUse, updateUser.UserId); delRolesErr != nil {
			ur.logger.Error("Failed to delete user roles", zap.Error(delRolesErr))
			txErr = delRolesErr
			return nil, delRolesErr
		}

		sql := `
			WITH updated_user AS (
				UPDATE users
				SET first_name = $1, last_name = $2, email = $3
				WHERE user_id = $4
				RETURNING user_id, first_name, last_name, email, password
			)
			INSERT INTO user_roles (user_id, role_id)
			SELECT u.user_id, r.role_id
			FROM updated_user u
			CROSS JOIN unnest($5::text[]) AS rn(role_name)
			INNER JOIN roles r ON r.role_name = rn.role_name
			RETURNING (SELECT user_id FROM updated_user), 
					  (SELECT first_name FROM updated_user), 
					  (SELECT last_name FROM updated_user), 
					  (SELECT email FROM updated_user), 
					  (SELECT password FROM updated_user);
		`

		ur.logger.Debug("Updating user with roles", zap.Bool("has_transaction", tx != nil))
		err = txToUse.QueryRow(ctx, sql, updateUser.FirstName, updateUser.LastName, updateUser.Email, updateUser.UserId, updateUser.Roles).
			Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	} else {
		sql := "UPDATE users " +
			" SET first_name = $1, last_name = $2, email = $3 " +
			"WHERE user_id = $4 " +
			"RETURNING user_id, first_name, last_name, email, password;"

		ur.logger.Debug("Updating user", zap.Bool("has_transaction", tx != nil))
		err = txToUse.QueryRow(ctx, sql, updateUser.FirstName, updateUser.LastName, updateUser.Email, updateUser.UserId).
			Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	}

	if err != nil {
		ur.logger.Error("Failed to update user", zap.Error(err))
		txErr = exceptions.PgErrorToAppError(err)
		return nil, txErr
	}

	if tx == nil {
		commitErr := ur.txFactory.CommitOrRollback(ctx, txToUse, nil)
		if commitErr != nil {
			txErr = commitErr
			return nil, commitErr
		}
		committed = true
	}

	user.Roles = updateUser.Roles
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

	var txToUse pgx.Tx
	var txErr *exceptions.AppError
	committed := false

	if tx != nil {
		txToUse = tx
	} else {
		txToUse, txErr = ur.txFactory.BeginTransaction(ctx)
		if txErr != nil {
			return txErr
		}
		defer func() {
			if !committed {
				_ = ur.txFactory.CommitOrRollback(ctx, txToUse, txErr)
			}
		}()
	}

	sql := "DELETE FROM users WHERE user_id = $1;"

	var err error
	var cmdTag pgconn.CommandTag

	ur.logger.Debug("Deleting user", zap.Bool("has_transaction", tx != nil))
	cmdTag, err = txToUse.Exec(ctx, sql, userId)

	if err != nil {
		ur.logger.Error("Failed to delete user", zap.Error(err))
		txErr = exceptions.PgErrorToAppError(err)
		return txErr
	}

	if cmdTag.RowsAffected() == 0 {
		ur.logger.Error("Failed to delete user - user not found", zap.Int32("user_id", userId))
		txErr = exceptions.NewNotFoundError("User not found")
		return txErr
	}

	delRolesErr := ur.deleteUserRoles(ctx, txToUse, userId)
	if delRolesErr != nil {
		ur.logger.Error("Failed to delete user roles", zap.Error(delRolesErr))
		txErr = delRolesErr
		return delRolesErr
	}

	if tx == nil {
		commitErr := ur.txFactory.CommitOrRollback(ctx, txToUse, nil)
		if commitErr != nil {
			txErr = commitErr
			return commitErr
		}
		committed = true
	}

	return nil
}

func (ur *UserRepository) GetUsersRoles(ctx context.Context, tx pgx.Tx, userIds []int32) (map[int32][]string, *exceptions.AppError) {
	if messages := domain_model.ValidateIds(userIds, true); len(messages) != 0 {
		return nil, exceptions.NewValidationError("Invalid IDs", messages)
	}

	sql := `
		SELECT ur.user_id, r.role_name 
		FROM user_roles ur
		INNER JOIN roles r ON ur.role_id = r.role_id
		WHERE ur.user_id = ANY($1)
		ORDER BY ur.user_id, r.role_name;
	`

	var rows pgx.Rows
	var err error

	if tx != nil {
		ur.logger.Debug("Getting users roles in transaction")
		rows, err = tx.Query(ctx, sql, userIds)
	} else {
		ur.logger.Debug("Getting users roles without transaction")
		rows, err = ur.db.Pool.Query(ctx, sql, userIds)
	}

	if err != nil {
		ur.logger.Error("Failed to get users roles", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	rolesMap := make(map[int32][]string)

	for _, userId := range userIds {
		rolesMap[userId] = []string{}
	}

	for rows.Next() {
		var userId int32
		var roleName string
		err = rows.Scan(&userId, &roleName)
		if err != nil {
			ur.logger.Error("Failed to scan user role row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}
		rolesMap[userId] = append(rolesMap[userId], roleName)
	}

	if err = rows.Err(); err != nil {
		ur.logger.Error("Error iterating user roles rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return rolesMap, nil
}

func (ur *UserRepository) AreRolesValid(ctx context.Context, roles []string) (bool, *exceptions.AppError) {
	if messages := domain_model.ValidateRoleNames(roles, true); len(messages) != 0 {
		return false, exceptions.NewValidationError("Invalid roles", messages)
	}

	validRolesSql := `
		SELECT COUNT(*) 
		FROM unnest($1::text[]) AS rn(role_name)
		LEFT JOIN roles r ON r.role_name = rn.role_name
		WHERE r.role_id IS NULL;
	`

	var invalidCount int
	var validRolesErr error

	validRolesErr = ur.db.Pool.QueryRow(ctx, validRolesSql, roles).Scan(&invalidCount)
	if validRolesErr != nil {
		ur.logger.Error("Failed to validate roles", zap.Error(validRolesErr))
		return false, exceptions.PgErrorToAppError(validRolesErr)
	}

	if invalidCount > 0 {
		ur.logger.Warn("Invalid roles provided", zap.Strings("roles", roles), zap.Int("invalid_count", invalidCount))
		return false, exceptions.NewAppError("ROLES_DO_NOT_EXIST", "The provided roles do not exist", nil)
	}

	return true, nil
}

func (ur *UserRepository) deleteUserRoles(ctx context.Context, tx pgx.Tx, userId int32) *exceptions.AppError {
	deleteRolesSql := "DELETE FROM user_roles WHERE user_id = $1;"
	var deleteErr error
	if tx != nil {
		_, deleteErr = tx.Exec(ctx, deleteRolesSql, userId)
	} else {
		_, deleteErr = ur.db.Pool.Exec(ctx, deleteRolesSql, userId)
	}

	if deleteErr != nil {
		ur.logger.Error("Failed to delete existing roles", zap.Error(deleteErr))
		return exceptions.PgErrorToAppError(deleteErr)
	}
	return nil
}
