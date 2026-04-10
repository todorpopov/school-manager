package directors

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
)

type IDirectorService interface {
	CreateDirector(ctx context.Context, tx pgx.Tx, createDirector *CreateDirector) (*Director, *exceptions.AppError)
	GetDirectorById(ctx context.Context, tx pgx.Tx, directorId int32) (*Director, *exceptions.AppError)
	GetDirectorByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Director, *exceptions.AppError)
	GetDirectors(ctx context.Context, tx pgx.Tx) ([]Director, *exceptions.AppError)
	UpdateDirector(ctx context.Context, tx pgx.Tx, updateDirector *UpdateDirector) (*Director, *exceptions.AppError)
	DeleteDirector(ctx context.Context, tx pgx.Tx, directorId int32) *exceptions.AppError
}

type DirectorService struct {
	directorRepo IDirectorRepository
	userSvc      users.IUserService
	txFactory    persistence.ITransactionFactory
}

func NewDirectorService(directorRepo IDirectorRepository, userSvc users.IUserService, txFactory persistence.ITransactionFactory) *DirectorService {
	return &DirectorService{directorRepo, userSvc, txFactory}
}

func (ds *DirectorService) CreateDirector(ctx context.Context, tx pgx.Tx, createDirector *CreateDirector) (*Director, *exceptions.AppError) {
	validationErr := ValidateCreateDirector(createDirector)
	if validationErr != nil {
		return nil, validationErr
	}

	var txToUse pgx.Tx
	var txErr *exceptions.AppError
	committed := false

	if tx != nil {
		txToUse = tx
	} else {
		txToUse, txErr = ds.txFactory.BeginTransaction(ctx)
		if txErr != nil {
			return nil, txErr
		}
		defer func() {
			if !committed {
				_ = ds.txFactory.CommitOrRollback(ctx, txToUse, txErr)
			}
		}()
	}

	createUser := &users.CreateUser{
		FirstName: createDirector.FirstName,
		LastName:  createDirector.LastName,
		Email:     createDirector.Email,
		Password:  createDirector.Password,
		Roles:     []string{"DIRECTOR"},
	}

	user, userErr := ds.userSvc.CreateUser(ctx, txToUse, createUser)
	if userErr != nil {
		txErr = userErr
		return nil, userErr
	}

	directorRecord, directorErr := ds.directorRepo.CreateDirector(ctx, txToUse, user.UserId)
	if directorErr != nil {
		txErr = directorErr
		return nil, directorErr
	}

	if tx == nil {
		commitErr := ds.txFactory.CommitOrRollback(ctx, txToUse, nil)
		if commitErr != nil {
			txErr = commitErr
			return nil, commitErr
		}
		committed = true
	}

	director := &Director{
		DirectorId: directorRecord.DirectorId,
		UserId:     user.UserId,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Roles:      user.Roles,
	}

	return director, nil
}

func (ds *DirectorService) GetDirectorById(ctx context.Context, tx pgx.Tx, directorId int32) (*Director, *exceptions.AppError) {
	if msg := domain_model.ValidateId(directorId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid director ID", map[string]string{"director_id": msg})
	}

	director, err := ds.directorRepo.GetDirectorById(ctx, tx, directorId)
	if err != nil {
		return nil, err
	}

	rolesMap, rolesErr := ds.userSvc.GetUsersRoles(ctx, tx, []int32{director.UserId})
	if rolesErr != nil {
		return nil, rolesErr
	}
	director.Roles = rolesMap[director.UserId]

	return director, nil
}

func (ds *DirectorService) GetDirectorByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Director, *exceptions.AppError) {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}

	director, err := ds.directorRepo.GetDirectorByUserId(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	rolesMap, rolesErr := ds.userSvc.GetUsersRoles(ctx, tx, []int32{director.UserId})
	if rolesErr != nil {
		return nil, rolesErr
	}
	director.Roles = rolesMap[director.UserId]

	return director, nil
}

func (ds *DirectorService) GetDirectors(ctx context.Context, tx pgx.Tx) ([]Director, *exceptions.AppError) {
	directors, err := ds.directorRepo.GetDirectors(ctx, tx)
	if err != nil {
		return nil, err
	}

	if len(directors) > 0 {
		userIds := make([]int32, len(directors))
		for i := range directors {
			userIds[i] = directors[i].UserId
		}

		rolesMap, rolesErr := ds.userSvc.GetUsersRoles(ctx, tx, userIds)
		if rolesErr != nil {
			return nil, rolesErr
		}

		for i := range directors {
			directors[i].Roles = rolesMap[directors[i].UserId]
		}
	}

	return directors, nil
}

func (ds *DirectorService) UpdateDirector(ctx context.Context, tx pgx.Tx, updateDirector *UpdateDirector) (*Director, *exceptions.AppError) {
	if validationErr := ValidateUpdateDirector(updateDirector); validationErr != nil {
		return nil, validationErr
	}

	var txToUse pgx.Tx
	var txErr *exceptions.AppError
	committed := false

	if tx != nil {
		txToUse = tx
	} else {
		txToUse, txErr = ds.txFactory.BeginTransaction(ctx)
		if txErr != nil {
			return nil, txErr
		}
		defer func() {
			if !committed {
				_ = ds.txFactory.CommitOrRollback(ctx, txToUse, txErr)
			}
		}()
	}

	director, getErr := ds.directorRepo.GetDirectorById(ctx, txToUse, updateDirector.DirectorId)
	if getErr != nil {
		txErr = getErr
		return nil, getErr
	}

	rolesMap, rolesErr := ds.userSvc.GetUsersRoles(ctx, txToUse, []int32{director.UserId})
	if rolesErr != nil {
		txErr = rolesErr
		return nil, rolesErr
	}
	currentRoles := rolesMap[director.UserId]

	updateUser := &users.UpdateUser{
		UserId:    director.UserId,
		FirstName: updateDirector.FirstName,
		LastName:  updateDirector.LastName,
		Email:     updateDirector.Email,
		Roles:     currentRoles,
	}

	user, userErr := ds.userSvc.UpdateUser(ctx, txToUse, updateUser)
	if userErr != nil {
		txErr = userErr
		return nil, userErr
	}

	if tx == nil {
		commitErr := ds.txFactory.CommitOrRollback(ctx, txToUse, nil)
		if commitErr != nil {
			txErr = commitErr
			return nil, commitErr
		}
		committed = true
	}

	updatedDirector := &Director{
		DirectorId: director.DirectorId,
		UserId:     user.UserId,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Roles:      user.Roles,
	}

	return updatedDirector, nil
}

func (ds *DirectorService) DeleteDirector(ctx context.Context, tx pgx.Tx, directorId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(directorId); msg != "" {
		return exceptions.NewValidationError("Invalid director ID", map[string]string{"director_id": msg})
	}
	return ds.directorRepo.DeleteDirector(ctx, tx, directorId)
}
