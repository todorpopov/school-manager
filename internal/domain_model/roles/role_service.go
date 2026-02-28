package roles

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type IRoleService interface {
	CreateRole(ctx context.Context, tx pgx.Tx, createRole *CreateRole) (*Role, *exceptions.AppError)
	GetRoleById(ctx context.Context, tx pgx.Tx, roleId int32) (*Role, *exceptions.AppError)
	GetRoles(ctx context.Context, tx pgx.Tx) ([]Role, *exceptions.AppError)
	DeleteRole(ctx context.Context, tx pgx.Tx, roleId int32) *exceptions.AppError
}

type RoleService struct {
	roleRepo IRoleRepository
}

func NewRoleService(roleRepo IRoleRepository) *RoleService {
	return &RoleService{roleRepo}
}

func (rs *RoleService) CreateRole(ctx context.Context, tx pgx.Tx, createRole *CreateRole) (*Role, *exceptions.AppError) {
	validationErr := ValidateCreateRole(createRole)
	if validationErr != nil {
		return nil, validationErr
	}
	return rs.roleRepo.CreateRole(ctx, tx, createRole)
}

func (rs *RoleService) GetRoleById(ctx context.Context, tx pgx.Tx, roleId int32) (*Role, *exceptions.AppError) {
	if msg := domain_model.ValidateId(roleId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid role ID", map[string]string{"role_id": msg})
	}
	return rs.roleRepo.GetRoleById(ctx, tx, roleId)
}

func (rs *RoleService) GetRoles(ctx context.Context, tx pgx.Tx) ([]Role, *exceptions.AppError) {
	return rs.roleRepo.GetRoles(ctx, tx)
}

func (rs *RoleService) DeleteRole(ctx context.Context, tx pgx.Tx, roleId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(roleId); msg != "" {
		return exceptions.NewValidationError("Invalid role ID", map[string]string{"role_id": msg})
	}
	return rs.roleRepo.DeleteRole(ctx, tx, roleId)
}
