package students

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
)

type IStudentService interface {
	CreateStudent(ctx context.Context, tx pgx.Tx, createStudent *CreateStudent) (*Student, *exceptions.AppError)
	GetStudentById(ctx context.Context, tx pgx.Tx, studentId int32) (*Student, *exceptions.AppError)
	GetStudentByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Student, *exceptions.AppError)
	GetStudents(ctx context.Context, tx pgx.Tx) ([]Student, *exceptions.AppError)
	UpdateStudent(ctx context.Context, tx pgx.Tx, updateStudent *UpdateStudent) (*Student, *exceptions.AppError)
	DeleteStudent(ctx context.Context, tx pgx.Tx, studentId int32) *exceptions.AppError
}

type StudentService struct {
	studentRepo IStudentRepository
	userSvc     users.IUserService
	txFactory   persistence.ITransactionFactory
}

func NewStudentService(studentRepo IStudentRepository, userSvc users.IUserService, txFactory persistence.ITransactionFactory) *StudentService {
	return &StudentService{studentRepo, userSvc, txFactory}
}

func (ss *StudentService) CreateStudent(ctx context.Context, tx pgx.Tx, createStudent *CreateStudent) (*Student, *exceptions.AppError) {
	validationErr := ValidateCreateStudent(createStudent)
	if validationErr != nil {
		return nil, validationErr
	}

	var txToUse pgx.Tx
	var txErr *exceptions.AppError
	committed := false

	if tx != nil {
		txToUse = tx
	} else {
		txToUse, txErr = ss.txFactory.BeginTransaction(ctx)
		if txErr != nil {
			return nil, txErr
		}
		defer func() {
			if !committed {
				_ = ss.txFactory.CommitOrRollback(ctx, txToUse, txErr)
			}
		}()
	}

	createUser := &users.CreateUser{
		FirstName: createStudent.FirstName,
		LastName:  createStudent.LastName,
		Email:     createStudent.Email,
		Password:  createStudent.Password,
		Roles:     []string{"STUDENT"},
	}

	user, userErr := ss.userSvc.CreateUser(ctx, txToUse, createUser)
	if userErr != nil {
		txErr = userErr
		return nil, userErr
	}

	studentRecord, studentErr := ss.studentRepo.CreateStudent(ctx, txToUse, createStudent.SchoolId, user.UserId, createStudent.ClassId)
	if studentErr != nil {
		txErr = studentErr
		return nil, studentErr
	}

	if tx == nil {
		commitErr := ss.txFactory.CommitOrRollback(ctx, txToUse, nil)
		if commitErr != nil {
			txErr = commitErr
			return nil, commitErr
		}
		committed = true
	}

	student, getErr := ss.studentRepo.GetStudentById(ctx, nil, studentRecord.StudentId)
	if getErr != nil {
		return nil, getErr
	}

	rolesMap, rolesErr := ss.userSvc.GetUsersRoles(ctx, nil, []int32{student.UserId})
	if rolesErr != nil {
		return nil, rolesErr
	}
	student.Roles = rolesMap[student.UserId]

	return student, nil
}

func (ss *StudentService) GetStudentById(ctx context.Context, tx pgx.Tx, studentId int32) (*Student, *exceptions.AppError) {
	if msg := domain_model.ValidateId(studentId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid student ID", map[string]string{"student_id": msg})
	}

	student, err := ss.studentRepo.GetStudentById(ctx, tx, studentId)
	if err != nil {
		return nil, err
	}

	rolesMap, rolesErr := ss.userSvc.GetUsersRoles(ctx, tx, []int32{student.UserId})
	if rolesErr != nil {
		return nil, rolesErr
	}
	student.Roles = rolesMap[student.UserId]

	return student, nil
}

func (ss *StudentService) GetStudentByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Student, *exceptions.AppError) {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}

	student, err := ss.studentRepo.GetStudentByUserId(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	rolesMap, rolesErr := ss.userSvc.GetUsersRoles(ctx, tx, []int32{student.UserId})
	if rolesErr != nil {
		return nil, rolesErr
	}
	student.Roles = rolesMap[student.UserId]

	return student, nil
}

func (ss *StudentService) GetStudents(ctx context.Context, tx pgx.Tx) ([]Student, *exceptions.AppError) {
	students, err := ss.studentRepo.GetStudents(ctx, tx)
	if err != nil {
		return nil, err
	}

	if len(students) > 0 {
		userIds := make([]int32, len(students))
		for i := range students {
			userIds[i] = students[i].UserId
		}

		rolesMap, rolesErr := ss.userSvc.GetUsersRoles(ctx, tx, userIds)
		if rolesErr != nil {
			return nil, rolesErr
		}

		for i := range students {
			students[i].Roles = rolesMap[students[i].UserId]
		}
	}

	return students, nil
}

func (ss *StudentService) UpdateStudent(ctx context.Context, tx pgx.Tx, updateStudent *UpdateStudent) (*Student, *exceptions.AppError) {
	if validationErr := ValidateUpdateStudent(updateStudent); validationErr != nil {
		return nil, validationErr
	}

	var txToUse pgx.Tx
	var txErr *exceptions.AppError
	committed := false

	if tx != nil {
		txToUse = tx
	} else {
		txToUse, txErr = ss.txFactory.BeginTransaction(ctx)
		if txErr != nil {
			return nil, txErr
		}
		defer func() {
			if !committed {
				_ = ss.txFactory.CommitOrRollback(ctx, txToUse, txErr)
			}
		}()
	}

	student, getErr := ss.studentRepo.GetStudentById(ctx, txToUse, updateStudent.StudentId)
	if getErr != nil {
		txErr = getErr
		return nil, getErr
	}

	rolesMap, rolesErr := ss.userSvc.GetUsersRoles(ctx, txToUse, []int32{student.UserId})
	if rolesErr != nil {
		txErr = rolesErr
		return nil, rolesErr
	}
	currentRoles := rolesMap[student.UserId]

	updateUser := &users.UpdateUser{
		UserId:    student.UserId,
		FirstName: updateStudent.FirstName,
		LastName:  updateStudent.LastName,
		Email:     updateStudent.Email,
		Roles:     currentRoles,
	}

	user, userErr := ss.userSvc.UpdateUser(ctx, txToUse, updateUser)
	if userErr != nil {
		txErr = userErr
		return nil, userErr
	}
	if updateStudent.SchoolId != nil && student.School.SchoolId != *updateStudent.SchoolId {
		updateSchoolErr := ss.studentRepo.UpdateStudentSchool(ctx, txToUse, updateStudent.StudentId, *updateStudent.SchoolId)
		if updateSchoolErr != nil {
			txErr = updateSchoolErr
			return nil, updateSchoolErr
		}
	}

	if student.Class != nil && updateStudent.ClassId != nil {		if student.Class.ClassId != *updateStudent.ClassId {
			updateClassErr := ss.studentRepo.UpdateStudentClass(ctx, txToUse, updateStudent.StudentId, updateStudent.ClassId)
			if updateClassErr != nil {
				txErr = updateClassErr
				return nil, updateClassErr
			}
		}
	} else if student.Class == nil && updateStudent.ClassId != nil {
		updateClassErr := ss.studentRepo.UpdateStudentClass(ctx, txToUse, updateStudent.StudentId, updateStudent.ClassId)
		if updateClassErr != nil {
			txErr = updateClassErr
			return nil, updateClassErr
		}
	} else if student.Class != nil && updateStudent.ClassId == nil {
		updateClassErr := ss.studentRepo.UpdateStudentClass(ctx, txToUse, updateStudent.StudentId, nil)
		if updateClassErr != nil {
			txErr = updateClassErr
			return nil, updateClassErr
		}
	}

	if tx == nil {
		commitErr := ss.txFactory.CommitOrRollback(ctx, txToUse, nil)
		if commitErr != nil {
			txErr = commitErr
			return nil, commitErr
		}
		committed = true
	}

	updatedStudent, getUpdatedErr := ss.studentRepo.GetStudentById(ctx, nil, updateStudent.StudentId)
	if getUpdatedErr != nil {
		return nil, getUpdatedErr
	}

	updatedStudent.Roles = user.Roles

	return updatedStudent, nil
}

func (ss *StudentService) DeleteStudent(ctx context.Context, tx pgx.Tx, studentId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(studentId); msg != "" {
		return exceptions.NewValidationError("Invalid student ID", map[string]string{"student_id": msg})
	}

	var txToUse pgx.Tx
	var txErr *exceptions.AppError
	committed := false

	if tx != nil {
		txToUse = tx
	} else {
		txToUse, txErr = ss.txFactory.BeginTransaction(ctx)
		if txErr != nil {
			return txErr
		}
		defer func() {
			if !committed {
				_ = ss.txFactory.CommitOrRollback(ctx, txToUse, txErr)
			}
		}()
	}

	student, getErr := ss.studentRepo.GetStudentById(ctx, txToUse, studentId)
	if getErr != nil {
		txErr = getErr
		return getErr
	}

	deleteErr := ss.studentRepo.DeleteStudent(ctx, txToUse, studentId)
	if deleteErr != nil {
		txErr = deleteErr
		return deleteErr
	}

	userDeleteErr := ss.userSvc.DeleteUser(ctx, txToUse, student.UserId)
	if userDeleteErr != nil {
		txErr = userDeleteErr
		return userDeleteErr
	}

	if tx == nil {
		commitErr := ss.txFactory.CommitOrRollback(ctx, txToUse, nil)
		if commitErr != nil {
			txErr = commitErr
			return commitErr
		}
		committed = true
	}

	return nil
}
