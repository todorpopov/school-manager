package classes

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type IClassService interface {
	CreateClass(ctx context.Context, tx pgx.Tx, createClass *CreateClass) (*Class, *exceptions.AppError)
	GetClassById(ctx context.Context, tx pgx.Tx, classId int32) (*Class, *exceptions.AppError)
	GetClasses(ctx context.Context, tx pgx.Tx) ([]Class, *exceptions.AppError)
	GetClassesBySchoolId(ctx context.Context, tx pgx.Tx, schoolId int32) ([]Class, *exceptions.AppError)
	DeleteClass(ctx context.Context, tx pgx.Tx, classId int32) *exceptions.AppError
}

type ClassService struct {
	classRepo IClassRepository
}

func NewClassService(classRepo IClassRepository) *ClassService {
	return &ClassService{classRepo}
}

func (cs *ClassService) CreateClass(ctx context.Context, tx pgx.Tx, createClass *CreateClass) (*Class, *exceptions.AppError) {
	validationErr := ValidateCreateClass(createClass)
	if validationErr != nil {
		return nil, validationErr
	}

	return cs.classRepo.CreateClass(ctx, tx, createClass)
}

func (cs *ClassService) GetClassById(ctx context.Context, tx pgx.Tx, classId int32) (*Class, *exceptions.AppError) {
	if msg := domain_model.ValidateId(classId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid class ID", map[string]string{"class_id": msg})
	}

	return cs.classRepo.GetClassById(ctx, tx, classId)
}

func (cs *ClassService) GetClasses(ctx context.Context, tx pgx.Tx) ([]Class, *exceptions.AppError) {
	return cs.classRepo.GetClasses(ctx, tx)
}

func (cs *ClassService) GetClassesBySchoolId(ctx context.Context, tx pgx.Tx, schoolId int32) ([]Class, *exceptions.AppError) {
	if msg := domain_model.ValidateId(schoolId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid school ID", map[string]string{"school_id": msg})
	}

	return cs.classRepo.GetClassesBySchoolId(ctx, tx, schoolId)
}

func (cs *ClassService) DeleteClass(ctx context.Context, tx pgx.Tx, classId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(classId); msg != "" {
		return exceptions.NewValidationError("Invalid class ID", map[string]string{"class_id": msg})
	}

	return cs.classRepo.DeleteClass(ctx, tx, classId)
}
