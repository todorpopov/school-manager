package teachers

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
)

type ITeacherService interface {
	CreateTeacher(ctx context.Context, tx pgx.Tx, createTeacher *CreateTeacher) (*Teacher, *exceptions.AppError)
	GetTeacherById(ctx context.Context, tx pgx.Tx, teacherId int32) (*Teacher, *exceptions.AppError)
	GetTeacherByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Teacher, *exceptions.AppError)
	GetTeachers(ctx context.Context, tx pgx.Tx) ([]Teacher, *exceptions.AppError)
	UpdateTeacher(ctx context.Context, tx pgx.Tx, updateTeacher *UpdateTeacher) (*Teacher, *exceptions.AppError)
	DeleteTeacher(ctx context.Context, tx pgx.Tx, teacherId int32) *exceptions.AppError
}

type TeacherService struct {
	teacherRepo ITeacherRepository
	userSvc     users.IUserService
	txFactory   persistence.ITransactionFactory
}

func NewTeacherService(teacherRepo ITeacherRepository, userSvc users.IUserService, txFactory persistence.ITransactionFactory) *TeacherService {
	return &TeacherService{teacherRepo, userSvc, txFactory}
}

func (ts *TeacherService) CreateTeacher(ctx context.Context, tx pgx.Tx, createTeacher *CreateTeacher) (*Teacher, *exceptions.AppError) {
	validationErr := ValidateCreateTeacher(createTeacher)
	if validationErr != nil {
		return nil, validationErr
	}

	var txToUse pgx.Tx
	var txErr *exceptions.AppError
	committed := false

	if tx != nil {
		txToUse = tx
	} else {
		txToUse, txErr = ts.txFactory.BeginTransaction(ctx)
		if txErr != nil {
			return nil, txErr
		}
		defer func() {
			if !committed {
				_ = ts.txFactory.CommitOrRollback(ctx, txToUse, txErr)
			}
		}()
	}

	createUser := &users.CreateUser{
		FirstName: createTeacher.FirstName,
		LastName:  createTeacher.LastName,
		Email:     createTeacher.Email,
		Password:  createTeacher.Password,
		Roles:     []string{"TEACHER"},
	}

	user, userErr := ts.userSvc.CreateUser(ctx, txToUse, createUser)
	if userErr != nil {
		txErr = userErr
		return nil, userErr
	}

	teacherRecord, teacherErr := ts.teacherRepo.CreateTeacher(ctx, txToUse, user.UserId)
	if teacherErr != nil {
		txErr = teacherErr
		return nil, teacherErr
	}

	if tx == nil {
		commitErr := ts.txFactory.CommitOrRollback(ctx, txToUse, nil)
		if commitErr != nil {
			txErr = commitErr
			return nil, commitErr
		}
		committed = true
	}

	teacher := &Teacher{
		TeacherId: teacherRecord.TeacherId,
		UserId:    user.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Roles:     user.Roles,
	}

	return teacher, nil
}

func (ts *TeacherService) GetTeacherById(ctx context.Context, tx pgx.Tx, teacherId int32) (*Teacher, *exceptions.AppError) {
	if msg := domain_model.ValidateId(teacherId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid teacher ID", map[string]string{"teacher_id": msg})
	}

	teacher, err := ts.teacherRepo.GetTeacherById(ctx, tx, teacherId)
	if err != nil {
		return nil, err
	}

	rolesMap, rolesErr := ts.userSvc.GetUsersRoles(ctx, tx, []int32{teacher.UserId})
	if rolesErr != nil {
		return nil, rolesErr
	}
	teacher.Roles = rolesMap[teacher.UserId]

	return teacher, nil
}

func (ts *TeacherService) GetTeacherByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Teacher, *exceptions.AppError) {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}

	teacher, err := ts.teacherRepo.GetTeacherByUserId(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	rolesMap, rolesErr := ts.userSvc.GetUsersRoles(ctx, tx, []int32{teacher.UserId})
	if rolesErr != nil {
		return nil, rolesErr
	}
	teacher.Roles = rolesMap[teacher.UserId]

	return teacher, nil
}

func (ts *TeacherService) GetTeachers(ctx context.Context, tx pgx.Tx) ([]Teacher, *exceptions.AppError) {
	teachers, err := ts.teacherRepo.GetTeachers(ctx, tx)
	if err != nil {
		return nil, err
	}

	if len(teachers) > 0 {
		userIds := make([]int32, len(teachers))
		for i := range teachers {
			userIds[i] = teachers[i].UserId
		}

		rolesMap, rolesErr := ts.userSvc.GetUsersRoles(ctx, tx, userIds)
		if rolesErr != nil {
			return nil, rolesErr
		}

		for i := range teachers {
			teachers[i].Roles = rolesMap[teachers[i].UserId]
		}
	}

	return teachers, nil
}

func (ts *TeacherService) UpdateTeacher(ctx context.Context, tx pgx.Tx, updateTeacher *UpdateTeacher) (*Teacher, *exceptions.AppError) {
	if validationErr := ValidateUpdateTeacher(updateTeacher); validationErr != nil {
		return nil, validationErr
	}

	var txToUse pgx.Tx
	var txErr *exceptions.AppError
	committed := false

	if tx != nil {
		txToUse = tx
	} else {
		txToUse, txErr = ts.txFactory.BeginTransaction(ctx)
		if txErr != nil {
			return nil, txErr
		}
		defer func() {
			if !committed {
				_ = ts.txFactory.CommitOrRollback(ctx, txToUse, txErr)
			}
		}()
	}

	teacher, getErr := ts.teacherRepo.GetTeacherById(ctx, txToUse, updateTeacher.TeacherId)
	if getErr != nil {
		txErr = getErr
		return nil, getErr
	}

	rolesMap, rolesErr := ts.userSvc.GetUsersRoles(ctx, txToUse, []int32{teacher.UserId})
	if rolesErr != nil {
		txErr = rolesErr
		return nil, rolesErr
	}
	currentRoles := rolesMap[teacher.UserId]

	updateUser := &users.UpdateUser{
		UserId:    teacher.UserId,
		FirstName: updateTeacher.FirstName,
		LastName:  updateTeacher.LastName,
		Email:     updateTeacher.Email,
		Roles:     currentRoles,
	}

	user, userErr := ts.userSvc.UpdateUser(ctx, txToUse, updateUser)
	if userErr != nil {
		txErr = userErr
		return nil, userErr
	}

	if tx == nil {
		commitErr := ts.txFactory.CommitOrRollback(ctx, txToUse, nil)
		if commitErr != nil {
			txErr = commitErr
			return nil, commitErr
		}
		committed = true
	}

	updatedTeacher := &Teacher{
		TeacherId: teacher.TeacherId,
		UserId:    user.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Roles:     user.Roles,
	}

	return updatedTeacher, nil
}

func (ts *TeacherService) DeleteTeacher(ctx context.Context, tx pgx.Tx, teacherId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(teacherId); msg != "" {
		return exceptions.NewValidationError("Invalid teacher ID", map[string]string{"teacher_id": msg})
	}
	return ts.teacherRepo.DeleteTeacher(ctx, tx, teacherId)
}

