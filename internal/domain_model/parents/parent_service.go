package parents

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
)

type IParentService interface {
	CreateParent(ctx context.Context, tx pgx.Tx, createParent *CreateParent) (*Parent, *exceptions.AppError)
	GetParentById(ctx context.Context, tx pgx.Tx, parentId int32) (*Parent, *exceptions.AppError)
	GetParentByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Parent, *exceptions.AppError)
	GetParents(ctx context.Context, tx pgx.Tx) ([]Parent, *exceptions.AppError)
	UpdateParent(ctx context.Context, tx pgx.Tx, updateParent *UpdateParent) (*Parent, *exceptions.AppError)
	DeleteParent(ctx context.Context, tx pgx.Tx, parentId int32) *exceptions.AppError
}

type ParentService struct {
	parentRepo IParentRepository
	userSvc    users.IUserService
	txFactory  persistence.ITransactionFactory
}

func NewParentService(parentRepo IParentRepository, userSvc users.IUserService, txFactory persistence.ITransactionFactory) *ParentService {
	return &ParentService{parentRepo, userSvc, txFactory}
}

func (ps *ParentService) CreateParent(ctx context.Context, tx pgx.Tx, createParent *CreateParent) (*Parent, *exceptions.AppError) {
	validationErr := ValidateCreateParent(createParent)
	if validationErr != nil {
		return nil, validationErr
	}

	var txToUse pgx.Tx
	var txErr *exceptions.AppError
	committed := false

	if tx != nil {
		txToUse = tx
	} else {
		txToUse, txErr = ps.txFactory.BeginTransaction(ctx)
		if txErr != nil {
			return nil, txErr
		}
		defer func() {
			if !committed {
				_ = ps.txFactory.CommitOrRollback(ctx, txToUse, txErr)
			}
		}()
	}

	createUser := &users.CreateUser{
		FirstName: createParent.FirstName,
		LastName:  createParent.LastName,
		Email:     createParent.Email,
		Password:  createParent.Password,
		Roles:     []string{"PARENT"},
	}

	user, userErr := ps.userSvc.CreateUser(ctx, txToUse, createUser)
	if userErr != nil {
		txErr = userErr
		return nil, userErr
	}

	parentRecord, parentErr := ps.parentRepo.CreateParent(ctx, txToUse, user.UserId)
	if parentErr != nil {
		txErr = parentErr
		return nil, parentErr
	}

	if tx == nil {
		commitErr := ps.txFactory.CommitOrRollback(ctx, txToUse, nil)
		if commitErr != nil {
			txErr = commitErr
			return nil, commitErr
		}
		committed = true
	}

	parent := &Parent{
		ParentId:  parentRecord.ParentId,
		UserId:    user.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Roles:     user.Roles,
	}

	return parent, nil
}

func (ps *ParentService) GetParentById(ctx context.Context, tx pgx.Tx, parentId int32) (*Parent, *exceptions.AppError) {
	if msg := domain_model.ValidateId(parentId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid parent ID", map[string]string{"parent_id": msg})
	}

	parent, err := ps.parentRepo.GetParentById(ctx, tx, parentId)
	if err != nil {
		return nil, err
	}

	rolesMap, rolesErr := ps.userSvc.GetUsersRoles(ctx, tx, []int32{parent.UserId})
	if rolesErr != nil {
		return nil, rolesErr
	}
	parent.Roles = rolesMap[parent.UserId]

	return parent, nil
}

func (ps *ParentService) GetParentByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Parent, *exceptions.AppError) {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}

	parent, err := ps.parentRepo.GetParentByUserId(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	rolesMap, rolesErr := ps.userSvc.GetUsersRoles(ctx, tx, []int32{parent.UserId})
	if rolesErr != nil {
		return nil, rolesErr
	}
	parent.Roles = rolesMap[parent.UserId]

	return parent, nil
}

func (ps *ParentService) GetParents(ctx context.Context, tx pgx.Tx) ([]Parent, *exceptions.AppError) {
	parents, err := ps.parentRepo.GetParents(ctx, tx)
	if err != nil {
		return nil, err
	}

	if len(parents) > 0 {
		userIds := make([]int32, len(parents))
		for i := range parents {
			userIds[i] = parents[i].UserId
		}

		rolesMap, rolesErr := ps.userSvc.GetUsersRoles(ctx, tx, userIds)
		if rolesErr != nil {
			return nil, rolesErr
		}

		for i := range parents {
			parents[i].Roles = rolesMap[parents[i].UserId]
		}
	}

	return parents, nil
}

func (ps *ParentService) UpdateParent(ctx context.Context, tx pgx.Tx, updateParent *UpdateParent) (*Parent, *exceptions.AppError) {
	if validationErr := ValidateUpdateParent(updateParent); validationErr != nil {
		return nil, validationErr
	}

	var txToUse pgx.Tx
	var txErr *exceptions.AppError
	committed := false

	if tx != nil {
		txToUse = tx
	} else {
		txToUse, txErr = ps.txFactory.BeginTransaction(ctx)
		if txErr != nil {
			return nil, txErr
		}
		defer func() {
			if !committed {
				_ = ps.txFactory.CommitOrRollback(ctx, txToUse, txErr)
			}
		}()
	}

	parent, getErr := ps.parentRepo.GetParentById(ctx, txToUse, updateParent.ParentId)
	if getErr != nil {
		txErr = getErr
		return nil, getErr
	}

	rolesMap, rolesErr := ps.userSvc.GetUsersRoles(ctx, txToUse, []int32{parent.UserId})
	if rolesErr != nil {
		txErr = rolesErr
		return nil, rolesErr
	}
	currentRoles := rolesMap[parent.UserId]

	updateUser := &users.UpdateUser{
		UserId:    parent.UserId,
		FirstName: updateParent.FirstName,
		LastName:  updateParent.LastName,
		Email:     updateParent.Email,
		Roles:     currentRoles,
	}

	user, userErr := ps.userSvc.UpdateUser(ctx, txToUse, updateUser)
	if userErr != nil {
		txErr = userErr
		return nil, userErr
	}

	if tx == nil {
		commitErr := ps.txFactory.CommitOrRollback(ctx, txToUse, nil)
		if commitErr != nil {
			txErr = commitErr
			return nil, commitErr
		}
		committed = true
	}

	updatedParent := &Parent{
		ParentId:  parent.ParentId,
		UserId:    user.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Roles:     user.Roles,
	}

	return updatedParent, nil
}

func (ps *ParentService) DeleteParent(ctx context.Context, tx pgx.Tx, parentId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(parentId); msg != "" {
		return exceptions.NewValidationError("Invalid parent ID", map[string]string{"parent_id": msg})
	}

	var txToUse pgx.Tx
	var txErr *exceptions.AppError
	committed := false

	if tx != nil {
		txToUse = tx
	} else {
		txToUse, txErr = ps.txFactory.BeginTransaction(ctx)
		if txErr != nil {
			return txErr
		}
		defer func() {
			if !committed {
				_ = ps.txFactory.CommitOrRollback(ctx, txToUse, txErr)
			}
		}()
	}

	parent, getErr := ps.parentRepo.GetParentById(ctx, txToUse, parentId)
	if getErr != nil {
		txErr = getErr
		return getErr
	}

	deleteErr := ps.parentRepo.DeleteParent(ctx, txToUse, parentId)
	if deleteErr != nil {
		txErr = deleteErr
		return deleteErr
	}

	userDeleteErr := ps.userSvc.DeleteUser(ctx, txToUse, parent.UserId)
	if userDeleteErr != nil {
		txErr = userDeleteErr
		return userDeleteErr
	}

	if tx == nil {
		commitErr := ps.txFactory.CommitOrRollback(ctx, txToUse, nil)
		if commitErr != nil {
			txErr = commitErr
			return commitErr
		}
		committed = true
	}

	return nil
}
