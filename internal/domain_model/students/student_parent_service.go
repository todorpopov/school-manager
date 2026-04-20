package students

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/parents"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type IStudentParentService interface {
	LinkParentToStudent(ctx context.Context, tx pgx.Tx, studentId int32, parentId int32) *exceptions.AppError
	UnlinkParentFromStudent(ctx context.Context, tx pgx.Tx, studentId int32, parentId int32) *exceptions.AppError
	GetParentsForStudent(ctx context.Context, tx pgx.Tx, studentId int32) ([]parents.Parent, *exceptions.AppError)
	GetStudentsForParent(ctx context.Context, tx pgx.Tx, parentId int32) ([]Student, *exceptions.AppError)
}

type StudentParentService struct {
	studentParentRepo IStudentParentRepository
}

func NewStudentParentService(studentParentRepo IStudentParentRepository) *StudentParentService {
	return &StudentParentService{studentParentRepo}
}

func (sps *StudentParentService) LinkParentToStudent(ctx context.Context, tx pgx.Tx, studentId int32, parentId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(studentId); msg != "" {
		return exceptions.NewValidationError("Invalid student ID", map[string]string{"student_id": msg})
	}
	if msg := domain_model.ValidateId(parentId); msg != "" {
		return exceptions.NewValidationError("Invalid parent ID", map[string]string{"parent_id": msg})
	}

	count, err := sps.studentParentRepo.CountParentsForStudent(ctx, tx, studentId)
	if err != nil {
		return err
	}

	if count >= MaxParentsPerStudent {
		return exceptions.NewValidationError(
			fmt.Sprintf("Student already has the maximum number of parents (%d)", MaxParentsPerStudent),
			map[string]string{"student_id": fmt.Sprintf("Student can have at most %d parents", MaxParentsPerStudent)},
		)
	}

	return sps.studentParentRepo.LinkParentToStudent(ctx, tx, studentId, parentId)
}

func (sps *StudentParentService) UnlinkParentFromStudent(ctx context.Context, tx pgx.Tx, studentId int32, parentId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(studentId); msg != "" {
		return exceptions.NewValidationError("Invalid student ID", map[string]string{"student_id": msg})
	}
	if msg := domain_model.ValidateId(parentId); msg != "" {
		return exceptions.NewValidationError("Invalid parent ID", map[string]string{"parent_id": msg})
	}

	return sps.studentParentRepo.UnlinkParentFromStudent(ctx, tx, studentId, parentId)
}

func (sps *StudentParentService) GetParentsForStudent(ctx context.Context, tx pgx.Tx, studentId int32) ([]parents.Parent, *exceptions.AppError) {
	if msg := domain_model.ValidateId(studentId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid student ID", map[string]string{"student_id": msg})
	}

	return sps.studentParentRepo.GetParentsForStudent(ctx, tx, studentId)
}

func (sps *StudentParentService) GetStudentsForParent(ctx context.Context, tx pgx.Tx, parentId int32) ([]Student, *exceptions.AppError) {
	if msg := domain_model.ValidateId(parentId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid parent ID", map[string]string{"parent_id": msg})
	}

	return sps.studentParentRepo.GetStudentsForParent(ctx, tx, parentId)
}
