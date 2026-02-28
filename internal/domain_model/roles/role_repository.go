package roles

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type IRoleRepository interface {
	CreateRole(ctx context.Context, tx pgx.Tx, createRole *CreateRole) (*Role, *exceptions.AppError)
	GetRoleById(ctx context.Context, tx pgx.Tx, roleId int32) (*Role, *exceptions.AppError)
	GetRoles(ctx context.Context, tx pgx.Tx) ([]Role, *exceptions.AppError)
	DeleteRole(ctx context.Context, tx pgx.Tx, roleId int32) *exceptions.AppError
}

type RoleRepository struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewRoleRepository(db *persistence.Database, logger *zap.Logger) *RoleRepository {
	return &RoleRepository{db, logger}
}

func (rr *RoleRepository) CreateRole(ctx context.Context, tx pgx.Tx, createRole *CreateRole) (*Role, *exceptions.AppError) {
	validationErr := ValidateCreateRole(createRole)
	if validationErr != nil {
		return nil, validationErr
	}

	sql := "INSERT INTO roles (role_name) VALUES ($1) RETURNING role_id, role_name;"

	var role Role
	var err error

	if tx != nil {
		rr.logger.Debug("Creating role in transaction")
		err = tx.QueryRow(ctx, sql, createRole.RoleName).Scan(&role.RoleId, &role.RoleName)
	} else {
		rr.logger.Debug("Creating role without transaction")
		err = rr.db.Pool.QueryRow(ctx, sql, createRole.RoleName).Scan(&role.RoleId, &role.RoleName)
	}

	if err != nil {
		rr.logger.Error("Failed to create role", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &role, nil
}

func (rr *RoleRepository) GetRoleById(ctx context.Context, tx pgx.Tx, roleId int32) (*Role, *exceptions.AppError) {
	if msg := domain_model.ValidateId(roleId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"role_id": msg})
	}

	sql := "SELECT role_id, role_name FROM roles WHERE role_id = $1;"

	var role Role
	var err error

	if tx != nil {
		rr.logger.Debug("Getting role by ID in transaction")
		err = tx.QueryRow(ctx, sql, roleId).Scan(&role.RoleId, &role.RoleName)
	} else {
		rr.logger.Debug("Getting role by ID without transaction")
		err = rr.db.Pool.QueryRow(ctx, sql, roleId).Scan(&role.RoleId, &role.RoleName)
	}

	if err != nil {
		rr.logger.Error("Failed to get role by ID", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &role, nil
}

func (rr *RoleRepository) GetRoles(ctx context.Context, tx pgx.Tx) ([]Role, *exceptions.AppError) {
	sql := "SELECT role_id, role_name FROM roles;"

	var roles []Role
	var err error
	var rows pgx.Rows

	if tx != nil {
		rr.logger.Debug("Getting roles in transaction")
		rows, err = tx.Query(ctx, sql)
	} else {
		rr.logger.Debug("Getting roles without transaction")
		rows, err = rr.db.Pool.Query(ctx, sql)
	}

	if err != nil {
		rr.logger.Error("Failed to get roles", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var role Role
		err = rows.Scan(&role.RoleId, &role.RoleName)
		if err != nil {
			rr.logger.Error("Failed to scan role", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (rr *RoleRepository) DeleteRole(ctx context.Context, tx pgx.Tx, roleId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(roleId); msg != "" {
		return exceptions.NewValidationError("Invalid user ID", map[string]string{"role_id": msg})
	}

	sql := "DELETE FROM roles WHERE role_id = $1;"

	var err error
	var cmdTag pgconn.CommandTag

	if tx != nil {
		rr.logger.Debug("Deleting role in transaction")
		cmdTag, err = tx.Exec(ctx, sql, roleId)
	} else {
		rr.logger.Debug("Deleting role without transaction")
		cmdTag, err = rr.db.Pool.Exec(ctx, sql, roleId)
	}

	if err != nil {
		rr.logger.Error("Failed to delete role", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		rr.logger.Debug("Failed to delete role - role not found", zap.Int32("role_id", roleId))
		return exceptions.NewNotFoundError("Role not found")
	}

	return nil
}
